# Cloud providers and known service URLs for storage providers

Last updated 20/04/21.

Generally, buckets are accessed as such, where `baseurl` may or may not include the region:

```
example.baseurl.tld
example.<region>.baseurl.tld
<region>.baseurl.tld/example
```

## Amazon (S3)

Cloudfront is not always for S3 buckets and supports subdomain only. 

Source: https://docs.aws.amazon.com/general/latest/gr/rande.html

```
<name>.s3.amazonaws.com
/[A-Za-z\d-\.]{3,63}\.s3\.amazonaws\.com/g

s3.amazonaws.com/<name>
/s3\.amazonaws\.com\/[A-Za-z\d-\.]{3,63}/g

<name>.s3.<region>.amazonaws.com
/[A-Za-z\d-\.]{3,63}\.s3\.[A-Za-z\d-]+\.amazonaws\.com/g

s3.<region>.amazonaws.com/<name>
/s3\.[A-Za-z\d-]\.amazonaws\.com\/[A-Za-z\d-\.]{3,63}/g

<name>.cloudfront.net
/[A-Za-z\d-\.]{3,63}\.cloudfront\.net/g
```

## Google Storage Buckets (modified S3, should be compatible)

Source: https://cloud.google.com/storage/docs/request-endpoints

```
<name>.storage.googleapis.com
/[A-Za-z\d-\.]{3,63}\.storage\.googleapis\.com/g

storage.googleapis.com/<name>
/storage\.googleapis\.com\/[A-Za-z\d-\.]{3,63}/g

www.googleapis.com/storage/v1/b/<name>/o/ (not S3 format)
/www.googleapis\.com\/storage\/v\d\/b\/[A-Za-z\d-\.]{3,63}/g
```

## DigitalOcean Spaces (S3)

Supports aliasing.

Source: https://docs.digitalocean.com/products/spaces/

```
<name>.<region>.digitaloceanspaces.com
/[A-Za-z\d-\.]{3,63}\.[A-Za-z\d-\.]+\.digitaloceanspaces\.com/g

<region>.digitaloceanspaces.com/<name>
/[A-Za-z\d-]+\.digitaloceanspaces\.com\/[A-Za-z\d-\.]{3,63}/g
```

## Azure Blob Storage (S3)

Buckets can be accessed by subdomain only.

```
<name>.blob.core.windows.net
/[A-Za-z\d-]{3,63}\.blob\.core\.windows\.net/g
```

## Linode (S3)

Source: https://status.linode.com

```
<name>.<region>.linodeobjects.com
/[a-z0-09][a-z0-9-]*[a-z0-9]?\.[A-Za-z\d-\.]+\.linodeobjects\.com/g

<region>.linodeobjects.com/<name>
/[A-Za-z\d-\.]+\.linodeobjects\.com\/[a-z0-09][a-z0-9-]*[a-z0-9]?/g
```

## Vultr (S3)

Source: https://www.vultr.com/docs/vultr-object-storage

```
<name>.<region>.vultrobjects.com
/[A-Za-z\d-\.]{1,63}\.[A-Za-z\d-\.]+\.vultrobjects\.com/g

<region>.vultrobjects.com/<name>
/[A-Za-z\d-\.]+\.vultrobjects\.com\/[A-Za-z\d-\.]{1,255}/g
```

## Backblaze (S3 after May 4th, 2020)

Not much info on this one yet.

Source: https://docs.fastly.com/en/guides/backblaze-b2-cloud-storage

```
<name>.s3.<region>.backblazeb2.com
/[A-Za-z\d-]{6,50}\.s3\.[A-Za-z\d-]+\.backblazeb2\.com/g

s3.<region>.backblazeb2.com/<name>
/s3\.[A-Za-z\d-]+\.backblazeb2\.com\/[A-Za-z\d-]{6,50}/g
```

## Wasabi (S3)

Source: https://wasabi-support.zendesk.com/hc/en-us/articles/360015106031-What-are-the-service-URLs-for-Wasabi-s-different-storage-regions-

```
<name>.s3.wasabisys.com
/[A-Za-z\d-\.]{3,63}\.s3\.wasabisys\.com/g

s3.wasabisys.com/<name>
/s3\.wasabisys\.com\/[A-Za-z\d-\.]{3,63}/g

<name>.s3.<region>.wasabisys.com
/[A-Za-z\d-\.]{3,63}\.s3\.[A-Za-z\d-]+\.wasabisys\.com/g

s3.<region>.wasabisys.com/<name>
/s3\.[A-Za-z\d-]\.wasabisys\.com\/[A-Za-z\d-\.]{3,63}/g
```

## DreamHost (S3)

Supports aliasing.

Source: https://help.dreamhost.com/hc/en-us/articles/215253408-How-to-create-a-DNS-alias-for-DreamObjects-buckets

```
<name>.objects-<region>.dream.io
/[A-Za-z\d-]{3,63}\.objects-[A-Za-z\d-]+\.dream\.io/g

objects-<region>.dream.io
/objects-[A-Za-z\d-]+\.dream\.io\/[A-Za-z\d-]{3,63}/g
```

## IBM Cloud Object Storage (S3)

Source: https://cloud.ibm.com/objectstorage/

```
<name>.s3.<region>.cloud-object-storage.appdomain.cloud
/[A-Za-z\d-\.]{3,63}\.s3\.[A-Za-z\d-]+\.cloud-object-storage\.appdomain\.cloud/g

s3.<region>.cloud-object-storage.appdomain.cloud/<name>
/s3\.[A-Za-z\d-]+\.cloud-object-storage\.appdomain\.cloud\/[A-Za-z\d-\.]{3,63}/g
```

## Firebase Storage (custom protocol)

Source: https://firebase.google.com/docs/storage/web/download-files

```
firebasestorage.googleapis.com/v0/b/<name>.appspot.com/o/
/firebasestorage\.googleapis\.com\/v\d\/b\/[A-Za-z\d-\.]+/g
```

## Generic S3 catcher

Very rough. Regexes should be used in listed order to prevent mismatching.

```
<name>.s3.<region>.<domain>.<tld>
/[A-Za-z\d-\.]+\.s3\.[A-Za-z\d-]+\.[A-Za-z\d-]+\.[A-Za-z\d-]{2,63}/g

s3.<region>.<domain>.<tld>/<name>
/s3\.[A-Za-z\d-]+\.[A-Za-z\d-]+\.[A-Za-z\d-]{2,63}\/[A-Za-z\d-\.]+/g

<name>.s3.<domain>.<tld>
/[A-Za-z\d-\.]+\.s3\.[A-Za-z\d-]+\.[A-Za-z\d-]{2,63}/g

s3.<domain>.<tld>/<name>
/s3\.[A-Za-z\d-]+\.[A-Za-z\d-]{2,63}\/[A-Za-z\d-\.]+/g
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