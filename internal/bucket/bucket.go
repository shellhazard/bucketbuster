package bucket

// Package bucket defines bucket types and data structures.

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"net/url"
	"regexp"
	"strings"

	"github.com/unva1idated/bucketbuster/internal/utils"
)

// Interface Bucket defines a deconstructed storage bucket.
type Bucket interface {
	// Returns the name of the bucket.
	Name() string

	// Returns the URL to access the root of the bucket.
	URL() string

	// Returns a URL to access a page of the bucket.
	PageURL(string) string

	// Returns a URL to download a specific resource in the bucket.
	ResourceURL(string) string

	// Parses page data and returns a list of keys found, and
	ParsePage([]byte) ([]string, string, error)
}

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

// Type S3Bucket represents a generic S3 compatible storage bucket.
type S3Bucket struct {
	// The base URL of the bucket.
	baseURL string
}

// Type S3BucketPage is a helper type for storing XML data.
type S3BucketPage struct {
	XMLName     xml.Name `xml:"ListBucketResult"`
	Text        string   `xml:",chardata"`
	Xmlns       string   `xml:"xmlns,attr"`
	Name        string   `xml:"Name"`
	Prefix      string   `xml:"Prefix"`
	Marker      string   `xml:"Marker"`
	MaxKeys     string   `xml:"MaxKeys"`
	IsTruncated bool     `xml:"IsTruncated"`
	Contents    []struct {
		Text         string `xml:",chardata"`
		Key          string `xml:"Key"`
		LastModified string `xml:"LastModified"`
		ETag         string `xml:"ETag"`
		Size         string `xml:"Size"`
		StorageClass string `xml:"StorageClass"`
	} `xml:"Contents"`
}

func NewS3Bucket(baseURL string) S3Bucket {
	return S3Bucket{
		baseURL: baseURL,
	}
}

// Returns the name of the bucket, which is the URL for a generic S3 bucket.
func (bucket S3Bucket) Name() string {
	return bucket.baseURL
}

// Returns the URL of the bucket.
func (bucket S3Bucket) URL() string {
	return bucket.baseURL
}

// Returns the URL pointing to the position in the bucket indicated by the pagination key.
func (bucket S3Bucket) PageURL(paginationKey string) string {
	if paginationKey == "" {
		return bucket.URL()
	}
	return fmt.Sprintf("%s?list-type=2&start-after=%s", bucket.URL(), paginationKey)
}

// Returns the URL used to fetch the resource with the specified key.
// TODO: Support extracting download key from metadata.
func (bucket S3Bucket) ResourceURL(key string) string {
	return fmt.Sprintf("%s/%s", bucket.URL(), url.QueryEscape(key))
}

// Parses a response to a page request and returns a slice of the keys
// and the next pagination key if applicable.
func (bucket S3Bucket) ParsePage(data []byte) ([]string, string, error) {
	var keys []string
	var page S3BucketPage
	var token string
	err := xml.Unmarshal(data, &page)
	if err != nil {
		return nil, "", err
	}
	for _, k := range page.Contents {
		keys = append(keys, k.Key)
	}
	if page.IsTruncated {
		token = keys[len(keys)-1]
	}
	return keys, token, nil
}

// -----

// Attempts to fingerprint the kind of bucket based on the URL.
// Returns a Bucket object.
func ParseURL(input string) Bucket {
	urlData, err := url.Parse(input)
	if err != nil {
		fmt.Printf("Couldn't parse input URL: %s", err)
		return nil
	}

	// Parse fragments
	pathFragments := utils.CleanStringSlice(strings.Split(urlData.Path, "/"))
	// hostFragments := utils.CleanStringSlice(strings.Split(urlData.Host, "."))

	// Fingerprint Firestore buckets
	matched, err := regexp.Match(`firebasestorage\.googleapis\.com\/v\d\/b\/[A-Za-z\d-\.]+`, []byte(input))
	if matched {
		if len(pathFragments) >= 3 {
			if pathFragments[0] == "v0" && pathFragments[1] == "b" {
				// Extract name
				return NewFirestoreBucket(pathFragments[2])
			}
		}
	}

	// Otherwise, assume it's a generic generic S3 bucket
	urlData.RawQuery = ""
	return NewS3Bucket(urlData.String())
}
