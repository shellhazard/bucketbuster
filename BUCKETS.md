# Cloud providers and known service URLs for storage providers

Last updated 18/04/21.

Generally, buckets are accessed as such, where `baseurl` may or may not include the region:

```
example.baseurl.tld
baseurl.tld/example
```

Regions are of course subject to change, so it's best to assume the region is a wildcard in your regexes. Known regions are listed here.

## Amazon (S3)

Source: https://docs.aws.amazon.com/general/latest/gr/rande.html

```
s3.amazonaws.com
s3.<region>.amazonaws.com
cloudfront.net (sometimes, subdomain only)
```

## Google Storage Buckets (modified S3, should be compatible)

Source: https://cloud.google.com/storage/docs/request-endpoints

```
storage.googleapis.com
```

## DigitalOcean Spaces (S3)

Supports aliasing.

Source: https://docs.digitalocean.com/products/spaces/

```
<region>.digitaloceanspaces.com
```

## Azure Blob Storage (S3)

Buckets can be accessed by subdomain only.

```
blob.core.windows.net
```

## Linode (S3)

Source: https://status.linode.com

```
<region>.linodeobjects.com
```

## Vultr (S3)

Source: Just trust me bro

```
<region>.vultrobjects.com
```

## Backblaze (S3)

Source: https://docs.fastly.com/en/guides/backblaze-b2-cloud-storage

```
s3.<region>.backblazeb2.com
backblazeb2.com
```

## Wasabi (S3)

Source: https://wasabi-support.zendesk.com/hc/en-us/articles/360015106031-What-are-the-service-URLs-for-Wasabi-s-different-storage-regions-

```
s3.wasabisys.com
s3.<region>.wasabisys.com
```

## DreamHost (S3)

Supports aliasing.

Source: https://help.dreamhost.com/hc/en-us/articles/215253408-How-to-create-a-DNS-alias-for-DreamObjects-buckets

```
objects-<region>.dream.io
```

## Rackspace (maybe S3?)

Haven't found a public listing for one of these yet.

Source: https://docs.rackspace.com/docs/cloud-files/v1/general-api-info/service-access/

```
storage101.<region>.clouddrive.com
snet-storage101.<region>.clouddrive.com
```

## IBM Cloud Object Storage (S3)

Source: https://cloud.ibm.com/objectstorage/

```
s3.<region>.cloud-object-storage.appdomain.cloud
```

## Firebase Storage (custom protocol)

Source: https://firebase.google.com/docs/storage/web/download-files

```
firebasestorage.googleapis.com/v0/b/<name>.appspot.com/o/
```