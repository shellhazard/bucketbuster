package bucket

import (
	"encoding/xml"
	"fmt"
	"net/url"
	"strings"
)

// Type S3Bucket represents a generic S3 compatible storage bucket.
type S3Bucket struct {
	// The base URL of the bucket.
	baseURL string

	// The generated name of the bucket.
	name string
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

func NewS3Bucket(baseURL string, name string) S3Bucket {
	return S3Bucket{
		baseURL: baseURL,
		name:    name,
	}
}

// Returns the name of the bucket, which is the URL for a generic S3 bucket.
func (bucket S3Bucket) Name() string {
	return bucket.name
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
	return fmt.Sprintf("%s?list-type=2&start-after=%s", bucket.URL(), url.QueryEscape(paginationKey))
}

// Returns the URL used to fetch the resource with the specified key.
// TODO: Support extracting download key from metadata.
func (bucket S3Bucket) ResourceURL(key string) string {
	burl := bucket.URL()
	if !strings.HasSuffix(burl, "/") {
		burl = fmt.Sprintf("%s/", burl)
	}
	return fmt.Sprintf("%s%s", burl, url.QueryEscape(key))
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
