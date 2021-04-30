# Cloud providers and known service URLs for storage providers

Last updated 23/04/21.

Generally, buckets are accessed as such, where `baseurl` may or may not include the region:

```
example.baseurl.tld
example.<region>.baseurl.tld
<region>.baseurl.tld/example
```

Regexes below are in PCRE format.

## Amazon (S3)

Cloudfront is not always for S3 buckets and supports subdomain only. 

Source: https://docs.aws.amazon.com/general/latest/gr/rande.html

```
<name>.s3.amazonaws.com
(?i)[A-Z\d-\.]{3,63}\.s3\.amazonaws\.com

s3.amazonaws.com/<name>
(?i)s3\.amazonaws\.com\/[A-Z\d-\.]{3,63}

<name>.s3.<region>.amazonaws.com
(?i)[A-Z\d-\.]{3,63}\.s3\.[A-Z\d-]+\.amazonaws\.com

s3.<region>.amazonaws.com/<name>
(?i)s3\.[A-Z\d-]\.amazonaws\.com\/[A-Z\d-\.]{3,63}

s3-<region>.amazonaws.com/<name>
(?i)s3-[A-Z\d-\.]+\.amazonaws\.com\/[A-Z\d-\.]{3,63}

<name>.cloudfront.net
(?i)[A-Z\d-\.]{3,63}\.cloudfront\.net/g
```

## Google Storage Buckets (modified S3, should be compatible)

Source: https://cloud.google.com/storage/docs/request-endpoints

```
<name>.storage.googleapis.com
(?i)[A-Z\d-\.]{3,63}\.storage\.googleapis\.com

storage.googleapis.com/<name>
(?i)storage\.googleapis\.com\/[A-Z\d-\.]{3,63}

www.googleapis.com/storage/v1/b/<name>/o/ (not S3 format)
(?i)www.googleapis\.com\/storage\/v\d\/b\/[A-Z\d-\.]{3,63}
```

## DigitalOcean Spaces (S3)

Supports aliasing.

Source: https://docs.digitalocean.com/products/spaces/

```
<name>.<region>.digitaloceanspaces.com
(?i)[A-Z\d-\.]{3,63}\.[A-Za-z\d-\.]+\.digitaloceanspaces\.com

<region>.digitaloceanspaces.com/<name>
(?i)[A-Z\d-]+\.digitaloceanspaces\.com\/[A-Z\d-\.]{3,63}
```

## Azure Blob Storage (S3)

Buckets can be accessed by subdomain only.

```
<name>.blob.core.windows.net
(?i)[A-Z\d-]{3,63}\.blob\.core\.windows\.net
```

## Linode (S3)

Source: https://status.linode.com

```
<name>.<region>.linodeobjects.com
(?i)[a-z0-09][a-z0-9-]*[a-z0-9]?\.[A-Za-z\d-\.]+\.linodeobjects\.com

<region>.linodeobjects.com/<name>
(?i)[a-z\d-\.]+\.linodeobjects\.com\/[a-z0-09][a-z0-9-]*[a-z0-9]?
```

## Vultr (S3)

Source: https://www.vultr.com/docs/vultr-object-storage

```
<name>.<region>.vultrobjects.com
(?i)[A-Z\d-\.]{1,63}\.[A-Z\d-\.]+\.vultrobjects\.com

<region>.vultrobjects.com/<name>
(?i)[A-Z\d-\.]+\.vultrobjects\.com\/[A-Z\d-\.]{1,255}
```

## Backblaze (S3 after May 4th, 2020)

Not much info on this one yet.

Source: https://docs.fastly.com/en/guides/backblaze-b2-cloud-storage

```
<name>.s3.<region>.backblazeb2.com
(?i)[A-Z\d-]{6,50}\.s3\.[A-Z\d-]+\.backblazeb2\.com

s3.<region>.backblazeb2.com/<name>
(?i)s3\.[A-Z\d-]+\.backblazeb2\.com\/[A-Z\d-]{6,50}
```

## Wasabi (S3)

Source: https://wasabi-support.zendesk.com/hc/en-us/articles/360015106031-What-are-the-service-URLs-for-Wasabi-s-different-storage-regions-

```
<name>.s3.wasabisys.com
(?i)[A-Z\d-\.]{3,63}\.s3\.wasabisys\.com

s3.wasabisys.com/<name>
(?i)s3\.wasabisys\.com\/[A-Z\d-\.]{3,63}

<name>.s3.<region>.wasabisys.com
(?i)[A-Z\d-\.]{3,63}\.s3\.[A-Z\d-]+\.wasabisys\.com

s3.<region>.wasabisys.com/<name>
(?i)s3\.[A-Z\d-]\.wasabisys\.com\/[A-Z\d-\.]{3,63}
```

## DreamHost (S3)

Supports aliasing.

Source: https://help.dreamhost.com/hc/en-us/articles/215253408-How-to-create-a-DNS-alias-for-DreamObjects-buckets

```
<name>.objects-<region>.dream.io
(?i)[A-Z\d-]{3,63}\.objects-[A-Z\d-]+\.dream\.io

objects-<region>.dream.io
(?i)objects-[A-Z\d-]+\.dream\.io\/[A-Z\d-]{3,63}
```

## IBM Cloud Object Storage (S3)

Source: https://cloud.ibm.com/objectstorage/

```
<name>.s3.<region>.cloud-object-storage.appdomain.cloud
(?i)[A-Z\d-\.]{3,63}\.s3\.[A-Z\d-]+\.cloud-object-storage\.appdomain\.cloud

s3.<region>.cloud-object-storage.appdomain.cloud/<name>
(?i)s3\.[A-Z\d-]+\.cloud-object-storage\.appdomain\.cloud\/[A-Z\d-\.]{3,63}
```

## Firebase Storage (custom protocol)

Source: https://firebase.google.com/docs/storage/web/download-files

```
firebasestorage.googleapis.com/v0/b/<name>.appspot.com/o/
(?i)firebasestorage\.googleapis\.com\/v\d\/b\/[A-Z\d-\.]+
```

## Generic S3 catcher

Very rough. Regexes should be used in listed order to prevent mismatching.

```
<name>.s3.<region>.<domain>.<tld>
(?i)[A-Za-z\d-\.]+\.s3\.[A-Za-z\d-]+\.[A-Za-z\d-]+\.[A-Za-z\d-]{2,63}

s3.<region>.<domain>.<tld>/<name>
(?i)s3\.[A-Za-z\d-]+\.[A-Za-z\d-]+\.[A-Za-z\d-]{2,63}\/[A-Za-z\d-\.]+

<name>.s3.<domain>.<tld>
(?i)[A-Za-z\d-\.]+\.s3\.[A-Za-z\d-]+\.[A-Za-z\d-]{2,63}

s3.<domain>.<tld>/<name>
(?i)s3\.[A-Za-z\d-]+\.[A-Za-z\d-]{2,63}\/[A-Za-z\d-\.]+
```


# Undocumented

## Rackspace (maybe S3?)

Haven't found a public listing for one of these yet.

Source: https://docs.rackspace.com/docs/cloud-files/v1/general-api-info/service-access/

```
storage101.<region>.clouddrive.com
snet-storage101.<region>.clouddrive.com
```

## Alibaba Cloud (OSS)

Haven't found a public listing for one of these yet.

Source: https://blog.anynines.com/concourse-configure-s3-resource-to-work-with-alibaba-cloud-oss/

```
oss-<region>.aliyuncs.com
```