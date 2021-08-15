package bucket

// Package bucket defines bucket types and data structures.
import (
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

	// Parses page data and returns a list of keys found and the next pagination key if applicable.
	ParsePage([]byte) ([]string, string, error)
}

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
