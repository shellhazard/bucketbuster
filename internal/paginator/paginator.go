package paginator

import (
	"io"
	"net/http"

	"github.com/unva1idated/bucketbuster/internal/bucket"
)

// Paginates the target bucket, fetching a page and returning a list
// of keys as well as the next pagination key if applicable.
func Paginate(b bucket.Bucket, paginationKey string) ([]string, string, error) {
	var keys []string
	var newPaginationKey string
	var targetURL string

	// Use the provided pagination key if not empty.
	targetURL = b.PageURL(paginationKey)

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
