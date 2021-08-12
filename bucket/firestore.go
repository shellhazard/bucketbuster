package bucket

import (
	"encoding/json"
	"fmt"
	"net/url"
)

// Type FirestoreBucket represents a Firebase storage bucket hosted on Google Cloud Platform.
type FirestoreBucket struct {
	// The bucket name.
	name string
}

// Type FirestoreBucketPage is a helper type for storing JSON data.
type FirestoreBucketPage struct {
	Prefixes []interface{} `json:"prefixes"`
	Items    []struct {
		Name   string `json:"name"`
		Bucket string `json:"bucket"`
	} `json:"items"`
	Nextpagetoken string `json:"nextPageToken"`
}

func NewFirestoreBucket(name string) FirestoreBucket {
	return FirestoreBucket{
		name: name,
	}
}

// Returns the name of the bucket.
func (bucket FirestoreBucket) Name() string {
	return bucket.name
}

// Returns the URL of the bucket.
func (bucket FirestoreBucket) URL() string {
	return fmt.Sprintf("https://firebasestorage.googleapis.com/v0/b/%s/o", bucket.name)
}

// Returns the URL pointing to the position in the bucket indicated by the pagination key.
func (bucket FirestoreBucket) PageURL(paginationKey string) string {
	if paginationKey == "" {
		return bucket.URL()
	}
	return fmt.Sprintf("%s?pageToken=%s", bucket.URL(), url.QueryEscape(paginationKey))
}

// Returns the URL used to fetch the resource with the specified key.
// TODO: Support extracting download key from metadata.
func (bucket FirestoreBucket) ResourceURL(key string) string {
	return fmt.Sprintf("%s/%s?alt=media", bucket.URL(), url.QueryEscape(key))
}

// Parses a response to a page request and returns a slice of the keys
// and the next pagination key if applicable.
func (bucket FirestoreBucket) ParsePage(data []byte) ([]string, string, error) {
	var keys []string
	var page FirestoreBucketPage
	var token string
	err := json.Unmarshal(data, &page)
	if err != nil {
		return nil, "", err
	}
	for _, k := range page.Items {
		keys = append(keys, k.Name)
	}
	if page.Nextpagetoken != "" {
		token = page.Nextpagetoken
	}
	return keys, token, nil
}
