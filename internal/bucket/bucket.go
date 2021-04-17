package bucket

// Package bucket defines bucket types and data structures.

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"net/url"
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

// Type AmazonS3Bucket represents a storage bucket hosted on Amazon S3.
type AmazonS3Bucket struct {
	// The bucket name.
	name string

	// The AWS server region of the bucket, if applicable.
	region string
}

// Type AmazonS3BucketPage is a helper type for storing XML data.
type AmazonS3BucketPage struct {
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

func NewAmazonS3Bucket(name, region string) AmazonS3Bucket {
	return AmazonS3Bucket{
		name:   name,
		region: region,
	}
}

func (bucket AmazonS3Bucket) Name() string {
	return bucket.name
}

// Returns the URL of the bucket.
func (bucket AmazonS3Bucket) URL() string {
	if bucket.region != "" {
		return fmt.Sprintf("https://%s.s3.%s.amazonaws.com", bucket.name, bucket.region)
	}
	return fmt.Sprintf("https://%s.s3.amazonaws.com", bucket.name)
}

// Returns the URL pointing to the position in the bucket indicated by the pagination key.
func (bucket AmazonS3Bucket) PageURL(paginationKey string) string {
	if paginationKey == "" {
		return bucket.URL()
	}
	return fmt.Sprintf("%s?list-type=2&start-after=%s", bucket.URL(), url.QueryEscape(paginationKey))
}

// Returns the URL used to fetch the resource with the specified key.
// TODO: Support extracting download key from metadata.
func (bucket AmazonS3Bucket) ResourceURL(key string) string {
	return fmt.Sprintf("%s/%s", bucket.URL(), url.QueryEscape(key))
}

// Parses a response to a page request and returns a slice of the keys
// and the next pagination key if applicable.
func (bucket AmazonS3Bucket) ParsePage(data []byte) ([]string, string, error) {
	var keys []string
	var page AmazonS3BucketPage
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

// Type DigitalOceanS3Bucket represents a storage bucket hosted on
// DigitalOceanSpaces.
type DigitalOceanS3Bucket struct {
	// The bucket name.
	name string

	// The DO server region of the bucket, if applicable.
	region string
}

func NewDigitalOceanS3Bucket(name, region string) DigitalOceanS3Bucket {
	return DigitalOceanS3Bucket{
		name:   name,
		region: region,
	}
}

// Returns the URL of the bucket.
func (bucket DigitalOceanS3Bucket) URL() string {
	if bucket.region != "" {
		return fmt.Sprintf("https://%s.%s.digitaloceanspaces.com", bucket.name, bucket.region)
	}
	return fmt.Sprintf("https://%s.nyc3.digitaloceanspaces.com", bucket.name)
}

// Type GoogleStorageS3Bucket represents a storage bucket hosted on Google's Storage API.
type GoogleStorageS3Bucket struct {
	// The bucket name.
	name string
}

func NewGoogleStorageS3Bucket(name string) GoogleStorageS3Bucket {
	return GoogleStorageS3Bucket{
		name: name,
	}
}

func (bucket GoogleStorageS3Bucket) URL() string {
	return fmt.Sprintf("https://%s.storage.googleapis.com/", bucket.name)
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

// Type GenericS3Bucket represents a generic S3 compatible storage bucket.
type GenericS3Bucket struct {
	// The base URL of the bucket.
	baseURL string
}

func NewGenericS3Bucket(baseURL string) GenericS3Bucket {
	return GenericS3Bucket{
		baseURL: baseURL,
	}
}

// Returns the name of the bucket, which is the URL for a generic S3 bucket.
func (bucket GenericS3Bucket) Name() string {
	return bucket.baseURL
}

// Returns the URL of the bucket.
func (bucket GenericS3Bucket) URL() string {
	return bucket.baseURL
}

// Returns the URL pointing to the position in the bucket indicated by the pagination key.
func (bucket GenericS3Bucket) PageURL(paginationKey string) string {
	if paginationKey == "" {
		return bucket.URL()
	}
	return fmt.Sprintf("%s?list-type=2&start-after=%s", bucket.URL(), paginationKey)
}

// Returns the URL used to fetch the resource with the specified key.
// TODO: Support extracting download key from metadata.
func (bucket GenericS3Bucket) ResourceURL(key string) string {
	return fmt.Sprintf("%s/%s", bucket.URL(), url.QueryEscape(key))
}

// Parses a response to a page request and returns a slice of the keys
// and the next pagination key if applicable.
func (bucket GenericS3Bucket) ParsePage(data []byte) ([]string, string, error) {
	// Not implemented
	return nil, "", nil
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
	workingPathFragments := strings.Split(urlData.Path, "/")
	pathFragments := make([]string, 0)
	for _, f := range workingPathFragments {
		trimmed := strings.TrimSpace(f)
		if len(trimmed) > 0 {
			pathFragments = append(pathFragments, trimmed)
		}
	}
	workingHostFragments := strings.Split(urlData.Host, ".")
	hostFragments := make([]string, 0)
	for _, f := range workingHostFragments {
		trimmed := strings.TrimSpace(f)
		if len(trimmed) > 0 {
			hostFragments = append(hostFragments, trimmed)
		}

	}

	// Fingerprint Firestore buckets
	if urlData.Host == "firebasestorage.googleapis.com" {
		if len(pathFragments) >= 3 {
			if pathFragments[0] == "v0" && pathFragments[1] == "b" {
				// Probably Firestore
				name := pathFragments[2]
				return NewFirestoreBucket(name)
			}
		}
	}

	// Fingerprint Amazon S3 buckets
	// https://s3.amazonaws.com/example
	// https://s3.eu-central-1.amazonaws.com/example
	// https://example.s3.amazonaws.com/
	// https://example.s3.eu-central-1.amazonaws.com/
	if strings.Contains(urlData.Host, "amazonaws.com") {
		hostFragsReversed := hostFragments
		utils.ReverseAny(hostFragsReversed)
		// First case
		if len(hostFragsReversed) == 3 && len(pathFragments) >= 1 {
			if hostFragsReversed[0] == "com" &&
				hostFragsReversed[1] == "amazonaws" &&
				hostFragsReversed[2] == "s3" {
				// Probably AWS
				return NewAmazonS3Bucket(pathFragments[0], "")
			}
		}
		// Second case
		if len(hostFragsReversed) == 4 && len(pathFragments) >= 1 {
			if hostFragsReversed[0] == "com" &&
				hostFragsReversed[1] == "amazonaws" &&
				hostFragsReversed[3] == "s3" {
				// Probably AWS
				return NewAmazonS3Bucket(pathFragments[0], hostFragsReversed[2])
			}
		}
		// Third case
		if len(hostFragsReversed) == 4 {
			if hostFragsReversed[0] == "com" &&
				hostFragsReversed[1] == "amazonaws" &&
				hostFragsReversed[2] == "s3" {
				// Probably AWS
				return NewAmazonS3Bucket(hostFragsReversed[3], "")
			}
		}
		// Fourth case
		if len(hostFragsReversed) == 5 {
			if hostFragsReversed[0] == "com" &&
				hostFragsReversed[1] == "amazonaws" &&
				hostFragsReversed[3] == "s3" {
				// Probably AWS
				return NewAmazonS3Bucket(hostFragsReversed[4], hostFragsReversed[2])
			}
		}
	}

	// Otherwise, assume it's a generic generic S3 bucket
	urlData.RawQuery = ""
	return NewGenericS3Bucket(urlData.String())
}
