package bucket

import (
	"encoding/xml"
	"fmt"
	"net/url"
	"strings"
)

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
