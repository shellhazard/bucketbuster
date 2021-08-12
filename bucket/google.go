package bucket

import (
	"encoding/xml"
	"fmt"
	"net/url"
	"strings"
)

// Type GoogleStorageBucket represents a Google Cloud Storage bucket.
type GoogleStorageBucket struct {
	// The name of the bucket.
	name string
}

// Type GoogleStorageBucketPage is a helper type for storing XML data.
type GoogleStorageBucketPage struct {
	XMLName     xml.Name `xml:"ListBucketResult"`
	Text        string   `xml:",chardata"`
	Xmlns       string   `xml:"xmlns,attr"`
	Name        string   `xml:"Name"`
	Prefix      string   `xml:"Prefix"`
	Marker      string   `xml:"Marker"`
	NextMarker  string   `xml:"NextMarker"`
	IsTruncated bool     `xml:"IsTruncated"`
	Contents    []struct {
		Text           string `xml:",chardata"`
		Key            string `xml:"Key"`
		Generation     string `xml:"Generation"`
		MetaGeneration string `xml:"MetaGeneration"`
		LastModified   string `xml:"LastModified"`
		ETag           string `xml:"ETag"`
		Size           string `xml:"Size"`
	} `xml:"Contents"`
}

func NewGoogleStorageBucket(name string) GoogleStorageBucket {
	return GoogleStorageBucket{
		name: name,
	}
}

// Returns the name of the bucket, which is the URL for a generic S3 bucket.
func (bucket GoogleStorageBucket) Name() string {
	return bucket.name
}

// Returns the URL of the bucket.
func (bucket GoogleStorageBucket) URL() string {
	return fmt.Sprintf("https://%s.storage.googleapis.com/", bucket.name)
}

// Returns the URL pointing to the position in the bucket indicated by the pagination key.
func (bucket GoogleStorageBucket) PageURL(paginationKey string) string {
	if paginationKey == "" {
		return bucket.URL()
	}
	return fmt.Sprintf("%s?marker=%s", bucket.URL(), url.QueryEscape(paginationKey))
}

// Returns the URL used to fetch the resource with the specified key.
// TODO: Support extracting download key from metadata.
func (bucket GoogleStorageBucket) ResourceURL(key string) string {
	burl := bucket.URL()
	if !strings.HasSuffix(burl, "/") {
		burl = fmt.Sprintf("%s/", burl)
	}
	return fmt.Sprintf("%s%s", burl, url.QueryEscape(key))
}

// Parses a response to a page request and returns a slice of the keys
// and the next pagination key if applicable.
func (bucket GoogleStorageBucket) ParsePage(data []byte) ([]string, string, error) {
	var keys []string
	var page GoogleStorageBucketPage
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
