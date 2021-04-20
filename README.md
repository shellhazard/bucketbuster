# bucketbuster

A standalone tool to analyse and index public Amazon S3 and Google Cloud Storage buckets without depending on monolithic SDKs.

## Purpose

I wanted a tool to index and analyse generic public storage buckets without depending on large SDKs from their providers and I couldn't find one that did everything I wanted it to do. I also wanted to be able to stop the execution of the program without losing the keys that had been indexed so far. I hacked this together in a night so there are probably bugs and only Amazon S3 and Firestore buckets are supported right now. Improvements are planned.

## Install

```
go install github.com/unva1idated/bucketbuster@latest
```

## Usage

The program will attempt to fingerprint the URL you provide with a specific cloud provider. If it can't find a match, it will assume your provided URL is the root of a generic S3 compatible storage bucket.

```
# Enumerate keys in the target S3 bucket and write their URLs to links.txt, then pass to wget
bucketbuster --url https://example.s3.eu-central-1.amazonaws.com -o links.txt
wget -i ./links.txt

# Enumerate keys in the target Firebase storage bucket and write 
# the keys and URLs to files.csv, then pass to massivedl
bucketbuster --url https://firebasestorage.googleapis.com/v0/b/example.appspot.com/o/ --format csv -o files.csv
massivedl -p 30 -i data.csv

# Start enumeration from a specific key and append key names to output.txt (without overwriting it)
bucketbuster -u https://example.s3.amazonaws.com -s examplekey -f key --append

```

## Notes

- Firebase storage buckets provide a specific key for pagination, but Amazon S3 buckets let you start from the key of any resource if you have it.
- Not all S3 providers use the same XML structure.

## Todo

- [Done] Indexing: Paginate S3 and Firestore bucket to determine number of keys, write keys to file.
- [Done] Different output formats: 
	- List of bucket keys.
	- List of URLs to pipe into wget.
	- Key/URL in csv format.
- Notify on finding keys with dangerous file extensions: .sql etc.
- Test bucket validity: Handle errors (NoSuchKey, AccessDenied, NoSuchBucket, PermanentRedirect)
- Write tests
- Index multiple buckets simultaneously: Accept list of bucket URLs as input.
- Support more storage bucket providers/URLs (some boilerplate is there already)
- Improve logging.
- Bucket analysis:
	- Test public index.
	- Test first three file access?
	- Test random three file access?
	- Test write (with a flag).

## Acknowledgements 

* The crew at [Exploitee.rs](https://exploitee.rs/) for being awesome.
* dimkouv, for [massivedl](https://github.com/dimkouv/massivedl)