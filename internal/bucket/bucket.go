package bucket

// Package bucket defines bucket types and data structures.

import (
	"encoding/json"
	"encoding/xml"
	"errors"
	"fmt"
	"net/url"
	"regexp"
	"strings"

	"github.com/shellhazard/bucketbuster/internal/utils"
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

// Type AzureStorageBucket represents an Azure Storage bucket.
type AzureStorageBucket struct {
	// The name of the account
	accountname string

	// The name of the container
	container string
}

// Type AzureStorageBucketPage is a helper type for storing XML data.
type AzureStorageBucketPage struct {
	XMLName         xml.Name `xml:"EnumerationResults"`
	Text            string   `xml:",chardata"`
	ServiceEndpoint string   `xml:"ServiceEndpoint,attr"`
	ContainerName   string   `xml:"ContainerName,attr"`
	Blobs           struct {
		Text string `xml:",chardata"`
		Blob []struct {
			Text       string `xml:",chardata"`
			Name       string `xml:"Name"`
			Url        string `xml:"Url"`
			Properties struct {
				Text               string `xml:",chardata"`
				LastModified       string `xml:"Last-Modified"`
				Etag               string `xml:"Etag"`
				ContentLength      string `xml:"Content-Length"`
				ContentType        string `xml:"Content-Type"`
				ContentEncoding    string `xml:"Content-Encoding"`
				ContentLanguage    string `xml:"Content-Language"`
				ContentMD5         string `xml:"Content-MD5"`
				CacheControl       string `xml:"Cache-Control"`
				ContentDisposition string `xml:"Content-Disposition"`
				BlobType           string `xml:"BlobType"`
				LeaseStatus        string `xml:"LeaseStatus"`
				LeaseState         string `xml:"LeaseState"`
			} `xml:"Properties"`
		} `xml:"Blob"`
	} `xml:"Blobs"`
	NextMarker string `xml:"NextMarker"`
}

func NewAzureStorageBucket(accountname string, container string) AzureStorageBucket {
	return AzureStorageBucket{
		accountname: accountname,
		container:   container,
	}
}

// Returns the name of the bucket, which is the URL for a generic S3 bucket.
func (bucket AzureStorageBucket) Name() string {
	return fmt.Sprintf("%s-%s", bucket.accountname, bucket.container)
}

// Returns the URL of the bucket.
func (bucket AzureStorageBucket) URL() string {
	return fmt.Sprintf("https://%s.blob.core.windows.net/%s", bucket.accountname, bucket.container)
}

// Returns the URL pointing to the position in the bucket indicated by the pagination key.
func (bucket AzureStorageBucket) PageURL(paginationKey string) string {
	if paginationKey == "" {
		return fmt.Sprintf("%s?restype=container&comp=list", bucket.URL())
	}
	return fmt.Sprintf("%s?restype=container&comp=list&marker=%s", bucket.URL(), url.QueryEscape(paginationKey))
}

// Returns the URL used to fetch the resource with the specified key.
// TODO: Support extracting download key from metadata.
func (bucket AzureStorageBucket) ResourceURL(key string) string {
	burl := bucket.URL()
	if !strings.HasSuffix(burl, "/") {
		burl = fmt.Sprintf("%s/", burl)
	}
	return fmt.Sprintf("%s%s", burl, url.QueryEscape(key))
}

// Parses a response to a page request and returns a slice of the keys
// and the next pagination key if applicable.
func (bucket AzureStorageBucket) ParsePage(data []byte) ([]string, string, error) {
	var keys []string
	var page AzureStorageBucketPage
	var token string
	err := xml.Unmarshal(data, &page)
	if err != nil {
		return nil, "", err
	}
	for _, k := range page.Blobs.Blob {
		keys = append(keys, k.Name)
	}
	if page.NextMarker != "" {
		token = page.NextMarker
	}
	return keys, token, nil
}

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

// -----

// Attempts to fingerprint the kind of bucket based on the URL.
// Returns a Bucket object.
func ParseURL(input string) (Bucket, error) {
	urlData, err := url.Parse(input)
	if err != nil {
		return nil, err
	}

	if urlData.Host == "" {
		return nil, errors.New("Invalid URL (missing scheme?)")
	}

	// Parse fragments
	pathFragments := utils.CleanStringSlice(strings.Split(urlData.Path, "/"))
	hostFragments := utils.CleanStringSlice(strings.Split(urlData.Host, "."))

	// Fingerprint Firestore buckets
	matched, err := regexp.Match(`(?i)firebasestorage\.googleapis\.com\/v\d\/b\/[A-Za-z\d-\.]+`, []byte(input))
	if err != nil {
		return nil, err
	}
	if matched {
		if len(pathFragments) >= 3 {
			if pathFragments[0] == "v0" && pathFragments[1] == "b" {
				// Extract name
				return NewFirestoreBucket(pathFragments[2]), nil
			}
		}
	}

	// Fingerprint Azure buckets
	matched, err = regexp.Match(`(?i)[A-Z\d-]{3,63}\.blob\.core\.windows\.net`, []byte(input))
	if err != nil {
		return nil, err
	}
	if matched {
		if len(pathFragments) >= 1 && len(hostFragments) >= 5 {
			// Extract name
			reversedHostFragments := hostFragments
			utils.ReverseAny(reversedHostFragments)
			return NewAzureStorageBucket(reversedHostFragments[4], pathFragments[0]), nil
		}
	}

	// Fingerprint Google Storage buckets
	// Name in subdomain
	matched, err = regexp.Match(`(?i)[A-Z\d-\.]{3,63}\.storage\.googleapis\.com`, []byte(input))
	if err != nil {
		return nil, err
	}
	if matched {
		if len(hostFragments) >= 4 {
			// Extract name
			reversedHostFragments := hostFragments
			utils.ReverseAny(reversedHostFragments)
			return NewGoogleStorageBucket(reversedHostFragments[3]), nil
		}
	}

	// Name in path
	matched, err = regexp.Match(`(?i)storage\.googleapis\.com\/[A-Z\d-\.]{3,63}`, []byte(input))
	if err != nil {
		return nil, err
	}
	if matched {
		if len(pathFragments) >= 1 {
			return NewGoogleStorageBucket(pathFragments[0]), nil
		}
	}

	// Otherwise, assume it's a generic generic S3 bucket
	urlData.RawQuery = ""

	// Build bucket name
	name := fmt.Sprintf("%s", urlData.Host)
	pathstring := strings.Join(pathFragments, "-")
	if pathstring != "" {
		name = fmt.Sprintf("%s-%s", urlData.Host, pathstring)
	}
	return NewS3Bucket(urlData.String(), name), nil
}
