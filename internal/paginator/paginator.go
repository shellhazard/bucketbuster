package paginator

import (
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/unva1idated/bucketbuster/internal/bucket"
)

const (
	esc = 27
)

// Paginates the target bucket, fetching a page and returning a list
// of keys as well as the next pagination key if applicable.
func Paginate(b bucket.Bucket, paginationKey string) ([]string, string, error) {
	var keys []string
	var newPaginationKey string
	var targetURL string

	// Use the provided pagination key if not empty.
	if paginationKey == "" {
		targetURL = b.URL()
	} else {
		targetURL = b.PageURL(paginationKey)
	}

	// Fetch the target page.
	resp, err := http.Get(targetURL)
	if err != nil {
		return nil, "", err
	}
	defer resp.Body.Close()

	// Read entire body into memory.
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, "", err
	}

	// Parse the page.
	keys, newPaginationKey, err = b.ParsePage(body)
	if err != nil {
		return nil, "", err
	}
	return keys, newPaginationKey, nil
}

func WritePaginatorStatus(keynum int, elapsed time.Duration) {
	// Move cursor up and clear line
	fmt.Printf("%c[%dA", esc)
	fmt.Printf("%c[2K\r", esc)
	fmt.Printf("\r\rTime elapsed: %s, Total keys: %d", elapsed.Round(1*time.Second), keynum)
}
