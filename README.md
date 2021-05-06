# bucketbuster

A standalone tool to analyse and index public storage buckets without depending on monolithic SDKs. Supports any S3 compliant storage service as well as Google Cloud Storage, Firestore and Azure.

## Purpose

I wanted a tool to index and analyse generic public storage buckets without depending on large SDKs from their providers and I couldn't find one that did everything I wanted it to do. I also wanted to be able to stop the execution of the program without losing the keys that had been indexed so far. I hacked this together in a night so there are probably bugs. Improvements are planned.

## Install

```
go install github.com/unva1idated/bucketbuster@latest
```

## Usage

The program will attempt to fingerprint the URL you provide with a specific cloud provider (Firebase). If it can't find a match, it will assume your provided URL is the root of a generic S3 compatible storage bucket.

```
# Enumerate keys in the target S3 bucket and write their URLs to links.txt, then pass to wget
bucketbuster --url https://example.s3.eu-central-1.amazonaws.com -o links.txt
wget -i ./links.txt

# Enumerate keys in the target Firebase storage bucket and write 
# the keys and URLs to files.csv, then pass to massivedl
bucketbuster --url https://firebasestorage.googleapis.com/v0/b/example.appspot.com/o/ --format csv -o files.csv
massivedl -p 30 -i data.csv

# Take a list of bucket URLs as input, indexing up to 30 at a time.
# Output is written to #-bucketurl.txt.
bucketbuster -i input-buckets.txt -c 30 -f csv

# Start enumeration from a specific key and append key names to output.txt (without overwriting it)
bucketbuster -u https://example.s3.amazonaws.com -s examplekey -f key --append
```

## Notes

- Any S3-compatible provider is supported, but I'm interested in supporting other public storage providers with their own APIs. Make an issue or PR if you want one added (preferably with an example URL).
- Firebase storage buckets provide a specific key for pagination, but S3 buckets let you start from the key of any resource if you have it.

## Todo

- [Done] Indexing: Paginate S3 and Firestore bucket to determine number of keys, write keys to file.
- [Done] Different output formats: 
	- List of bucket keys.
	- List of URLs to pipe into wget.
	- Key/URL in csv format.
- [Done] Index multiple buckets simultaneously: Accept list of bucket URLs as input.
- [Done, kind of] Improve logging.
- Support more storage bucket providers/URLs
- Notify on finding keys with dangerous file extensions: .sql etc.
- Test bucket validity: Handle errors (NoSuchKey, AccessDenied, NoSuchBucket, PermanentRedirect)
- Write tests

## Acknowledgements 

* The crew at [Exploitee.rs](https://exploitee.rs/) for being awesome.
* dimkouv, for [massivedl](https://github.com/dimkouv/massivedl)
