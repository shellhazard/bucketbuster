package cmd

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"os"
	"os/signal"
	"sync/atomic"
	"syscall"
	"time"

	logger "log"

	"github.com/spf13/cobra"
	"github.com/unva1idated/bucketbuster/internal/bucket"
	"github.com/unva1idated/bucketbuster/internal/paginator"
)

const (
	version = "v0.2"

	// Escape character
	esc = 27
)

var (
	// Flags
	url         string // The URL to parse.
	startkey    string // The key to start paginating from.
	outfile     string // The path of the output file.
	format      string // The format to write the output in.
	appendFile  bool   // Enable or disable appending to the target file instead of overwriting it.
	verbose     bool   // Enable extended output from buckets.
	input       string // The list of bucket URLs to index.
	concurrency int    // The maximum number of buckets to index simultaneously.

	// Loggers
	log     = logger.New(os.Stderr, "[bucketbuster] ", 0)
	worklog = logger.New(os.Stderr, "[bucketbuster] ", 0)

	rootCmd = &cobra.Command{
		Use:   "bucketbuster",
		Short: "Tool for indexing public storage buckets.",
		Long: `A standalone tool to analyse and index public Amazon S3 
and Google Cloud Storage buckets. See github.com/unva1idated/bucketbuster`,
		Run: func(cmd *cobra.Command, args []string) {
			log.Printf("Starting Bucketbuster %s.\n", version)

			// Prepare timer and counters
			start := time.Now()

			// Atomic counters. Maybe overkill but works nicely.
			var totalKeys *int64 = new(int64)
			var startedBuckets *int64 = new(int64)
			var completedBuckets *int64 = new(int64)

			ticker := time.NewTicker(500 * time.Millisecond)

			// Prepare status ticker
			go func() {
				for {
					_, ok := <-ticker.C
					WritePaginatorStatus(start, totalKeys, startedBuckets, completedBuckets)
					if !ok {
						return
					}

				}
			}()
			defer ticker.Stop()

			// Parse input file
			if input != "" {
				// Disable extended output
				if !verbose {
					worklog.SetOutput(ioutil.Discard)
				}

				// Attempt to load file
				file, err := os.Open(input)
				if err != nil {
					log.Fatalf("Failed to open input file: %s", err)
				}
				scanner := bufio.NewScanner(file)
				scanner.Split(bufio.ScanLines)

				log.Printf("Loading input URLs from file %s.", input)
				// Catch interrupt globally
				c := make(chan os.Signal, 2)
				signal.Notify(c, os.Interrupt, syscall.SIGTERM)
				go func() {
					for {
						_, ok := <-c
						if ok {
							fmt.Println("")
							log.Fatalf("Interrupted.")
						}
						return
					}
				}()
				defer close(c)

				// Read URLs from the file line by line
				sem := make(chan bool, concurrency)
				for scanner.Scan() {
					b, err := bucket.ParseURL(scanner.Text())
					if err != nil {
						log.Printf("Error parsing URL: %s")
						continue
					}
					sem <- true
					go func() {
						defer func() {
							<-sem
							atomic.AddInt64(completedBuckets, 1)
						}()
						atomic.AddInt64(startedBuckets, 1)
						IndexBucket(b, totalKeys, startedBuckets, completedBuckets, false)
					}()
				}
				for i := 0; i < cap(sem); i++ {
					sem <- true
				}
				// Parse URL parameter
			} else if url != "" {
				b, err := bucket.ParseURL(url)
				if err != nil {
					log.Fatalf("Failed to parse input URL: %s", err)
				}
				atomic.AddInt64(startedBuckets, 1)
				IndexBucket(b, totalKeys, startedBuckets, completedBuckets, true)
				atomic.AddInt64(completedBuckets, 1)
			}

			// Completed work
			WritePaginatorStatus(start, totalKeys, startedBuckets, completedBuckets)
			fmt.Println("")
			log.Println("Done.")
		},
	}
)

// Prepare flags
func init() {
	// Set command line flags
	rootCmd.PersistentFlags().StringVarP(&url, "url", "u", "", "The URL of a bucket to analyse. Required.")
	rootCmd.PersistentFlags().StringVarP(&startkey, "startkey", "s", "", "Specify the key to start paginating from if required. Ignored if using --input flag.")
	rootCmd.PersistentFlags().StringVarP(&input, "input", "i", "", "A list of bucket URLs to index.")
	rootCmd.PersistentFlags().StringVarP(&outfile, "outfile", "o", "", "The file to output keys/URLs to. Default {number}-{bucket-url}.txt. Ignored if using --input flag.")
	rootCmd.PersistentFlags().StringVarP(&format, "format", "f", "url", "Specify the output format. \"url\" is the default and outputs resource URLs, \"key\" outputs the list of keys. \"csv\" outputs as key,url for use with massivedl.")
	rootCmd.PersistentFlags().BoolVarP(&appendFile, "append", "a", false, "Appends to the target file instead of overwriting it.")
	rootCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "Detailed logging output.")
	rootCmd.PersistentFlags().IntVarP(&concurrency, "concurrency", "c", 10, "The maximum number of buckets to index simultaneously. Default 10.")

	// rootCmd.MarkPersistentFlagRequired("url")
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func IndexBucket(b bucket.Bucket, keyCounter *int64, startedBuckets *int64, completedBuckets *int64, single bool) {
	// Prepare state
	var keys []string
	paginationKey := ""
	if single {
		paginationKey = startkey
	}
	first := true
	filename := ""

	// Determine filename
	if single {
		filename = fmt.Sprintf("%s.txt", b.Name())
	} else {
		filename = fmt.Sprintf("%v-%s.txt", atomic.LoadInt64(startedBuckets), b.Name())
	}

	// Open file and prepare writer
	var file *os.File
	_, err := os.Stat(filename)
	if err == nil && appendFile {
		worklog.Println("Appending to existing file %s", filename)
		f, err := os.OpenFile(filename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			fmt.Fprintln(os.Stderr, "Error opening outfile: %s", err)
			os.Exit(1)
		}
		file = f
	} else {
		worklog.Printf("Creating new file %s.", filename)
		f, err := os.Create(filename)
		if err != nil {
			fmt.Fprintln(os.Stderr, "Error creating outfile: %s", err)
			os.Exit(1)
		}
		file = f
	}
	writer := bufio.NewWriterSize(file, 128)

	// Set up signal channel to catch program termination
	c := make(chan os.Signal, 2)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		for {
			_, ok := <-c
			writer.Flush()
			file.Close()
			if ok && single {
				fmt.Println("")
				log.Printf("Interrupted. Last pagination key: %s\n", paginationKey)
			}
			return
		}
	}()
	defer close(c)

	worklog.Printf("Counting keys in %s.", b.Name())
	if paginationKey != "" {
		worklog.Printf("Starting from key %s.", paginationKey)
	}

	// Pull each page of the bucket
	for paginationKey != "" || first == true {
		first = false
		newKeys, newPaginationKey, err := paginator.Paginate(b, paginationKey)
		if err != nil {
			fmt.Println("")
			log.Printf("Error during pagination: %s", err)
			return
		}
		// Write each key to the file as we receive it
		// This way even if our program is cancelled, we can resume
		// from the most recent key.
		for _, k := range newKeys {
			keys = append(keys, k)
			atomic.AddInt64(keyCounter, 1)

			var writestr string
			switch format {
			case "keys":
				writestr = fmt.Sprintf("%s\n", k)
			case "csv":
				writestr = fmt.Sprintf("%s,%s\n", k, b.ResourceURL(k))
			default:
				writestr = fmt.Sprintf("%s\n", b.ResourceURL(k))
			}

			_, writeErr := writer.WriteString(writestr)
			if writeErr != nil {
				log.Printf("Error during write: %s", writeErr)
				return
			}
		}
		paginationKey = newPaginationKey
	}
}

func WritePaginatorStatus(start time.Time, keys *int64, startedBuckets *int64, completedBuckets *int64) {
	// Calc time
	elapsed := time.Since(start)

	// Move cursor up and clear line
	fmt.Printf("%c[%dA", esc)
	fmt.Printf("%c[2K\r", esc)
	fmt.Printf("\r\r[bucketbuster] Elapsed: %s, Total keys: %v, Buckets started: %v, Buckets completed: %v", elapsed.Round(1*time.Second), atomic.LoadInt64(keys), atomic.LoadInt64(startedBuckets), atomic.LoadInt64(completedBuckets))
}
