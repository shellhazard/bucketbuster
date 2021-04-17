package cmd

import (
	"bufio"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/spf13/cobra"
	"github.com/unva1idated/bucketbuster/internal/bucket"
	"github.com/unva1idated/bucketbuster/internal/paginator"
)

const (
	version = "v0.1"
)

var (
	// Flags
	url        string // The URL to parse.
	startkey   string // The key to start paginating from.
	outfile    string // The path of the output file.
	format     string // The format to write the output in.
	appendFile bool   // Enable or disable appending to the target file instead of overwriting it.
	quiet      bool   // Enable or disable progress messages.

	rootCmd = &cobra.Command{
		Use:   "bucketbuster",
		Short: "Tool for indexing public storage buckets.",
		Long: `A standalone tool to analyse and index public Amazon S3 
and Google Cloud Storage buckets. See github.com/unva1idated/bucketbuster`,
		Run: func(cmd *cobra.Command, args []string) {
			if quiet == false {
				fmt.Printf("Starting Bucketbuster %s.\n", version)
			}
			b := bucket.ParseURL(url)

			// Prepare state
			var keys []string
			paginationKey := startkey
			first := true
			start := time.Now()

			// Open file and prepare writer
			var file *os.File
			_, err := os.Stat(outfile)
			if err == nil && appendFile {
				fmt.Println("Appending to existing file.")
				f, err := os.OpenFile(outfile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
				if err != nil {
					fmt.Fprintln(os.Stderr, "Error opening outfile: %s", err)
					os.Exit(1)
				}
				file = f
			} else {
				fmt.Println("Creating new file.")
				f, err := os.Create(outfile)
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
					if ok {
						fmt.Printf("\nInterrupted. Last pagination key: %s\n", paginationKey)

					}
					os.Exit(1)
					return
				}
			}()
			defer close(c)

			if quiet == false {
				fmt.Printf("Counting keys in %s\n.", b.Name())
				if paginationKey != "" {
					fmt.Printf("Starting from key %s.", paginationKey)
				}
				paginator.WritePaginatorStatus(0, time.Duration(0))
			}

			// Pull each page of the bucket
			for paginationKey != "" || first == true {
				first = false
				newKeys, newPaginationKey, err := paginator.Paginate(b, paginationKey)
				if err != nil {
					fmt.Fprintln(os.Stderr, "\nError during pagination:", err)
					os.Exit(1)
				}
				// Write each key to the file as we receive it
				// This way even if our program is cancelled, we can resume
				// from the most recent key.
				for _, k := range newKeys {
					keys = append(keys, k)

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
						fmt.Fprintln(os.Stderr, "\nError during write:", writeErr)
						os.Exit(1)
					}
				}
				paginationKey = newPaginationKey
				elapsed := time.Since(start)
				if quiet == false {
					paginator.WritePaginatorStatus(len(keys), elapsed)
				}
			}
			if quiet == false {
				fmt.Printf("\nDone.\n")
			}
		},
	}
)

// Prepare flags
func init() {
	rootCmd.PersistentFlags().StringVarP(&url, "url", "u", "", "The URL of a bucket to analyse. Required.")
	rootCmd.PersistentFlags().StringVarP(&startkey, "startkey", "s", "", "Specify the key to start paginating from if required.")
	rootCmd.PersistentFlags().StringVarP(&outfile, "outfile", "o", "bucket.txt", "The file to output keys/URLs to. Default bucket.txt.")
	rootCmd.PersistentFlags().StringVarP(&format, "format", "f", "url", "Specify the output format. \"url\" is the default and outputs resource URLs, \"key\" outputs the list of keys. \"csv\" outputs as key,url for use with massivedl.")
	rootCmd.PersistentFlags().BoolVarP(&appendFile, "append", "a", false, "Appends to the target file instead of overwriting it.")
	rootCmd.PersistentFlags().BoolVarP(&quiet, "quiet", "q", false, "Disables writing to stdout.")

	rootCmd.MarkPersistentFlagRequired("url")
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
