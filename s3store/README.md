- [s3store](#s3store)
- [Supported tools](#supported-tools)
  - [Upload an entire folder to an S3 bucket](#upload-an-entire-folder-to-an-s3-bucket)
  - [Download an entire bucket to a local folder](#download-an-entire-bucket-to-a-local-folder)
- [Using Noobaa](#using-noobaa)

# s3store
Golang library to transfer files from and to S3 compatible storage.
* Targetted to run on K8s/OpenShift environments, can run as CLI application
* Easy to integrate as imported package in any client application

# Supported tools
The following tools leverage the [s3 Golang package](https://pkg.go.dev/github.com/aws/aws-sdk-go/service/s3)
to implement utility functions that simplify the client applications.

The tools are connected to the S3 service using the following environment variables:
* `AWS_ACCESS_KEY_ID`: the AWS access key
* `AWS_SECRET_ACCESS_KEY`: the AWS secret key
* `S3_ENDPOINT`: the S3 endpoint
* `S3_REGION`: the S3 region

This function in the `s3filemanager` package creates as [AWS Session](https://pkg.go.dev/github.com/aws/aws-sdk-go/aws/session)
using the values from the above environment variables:
```golang
func ConnectWithEnvVariables() (*session.Session, error) {
  ...
}
```

Once the `Session` has been established, it is used to execute the tools.

## Upload an entire folder to an S3 bucket
[S3FolderUploader](./s3filemanager/s3_folder_uploader.go) is a Golang struct that implements the function to upload a given local folder to a remote S3 bucket.

The bucket is created, if it doesn't exist.

All subfolders will be included, and the target bucket preserves the folder structure.

The following snippet creates and execute the `S3FolderUploader` tool:
```golang
  ...
  session, err := s3filemanager.ConnectWithEnvVariables()
  if err != nil {
    log.Fatal(err)
  }

  uploader := s3filemanager.NewS3FolderUploader(bucket, folder, logger)
  err = uploader.Run(session)
  ...
```

`WithMatcher` function allows to select only the files whose path contains the given `matcher` string.

`WithPrexif` function allows to add a fixed prefix to all the file names.

`S3FolderUploader` uses [s3manager.Uploader](https://pkg.go.dev/github.com/stripe/aws-go/service/s3/s3manager#Uploader) 
to upload the entire folder recursively, using the `UploadWithIterator` scanner function to iterate through the files to be uploaded.

> This implementation minimizes the memory footprint and the number of open files: files are read and closed just after the read is complete.
> It defines a `lazyFileReader` implementation of the `io.Reader` interface to postpone the actual file reading until it's really needed.<br/>

`S3FolderUploader` can be tested using the provided `main` function with the following options:
```bash
go run main.go -b BUCKET_NAME -debug -f SOURCE_FOLDER -m upload
```

## Download an entire bucket to a local folder
[S3BucketDownloader](./s3filemanager/s3_bucket_downloader.go) is a Golang struct that implements the function to fully download a remote S3 bucket to a local folder.

The destination folder is created, if it doesn't exist.

All subfolders will be included, and the target folder preserves the bucket structure.

The following snippet creates and execute the `S3BucketDownloader` tool:
```golang
  ...
  session, err := s3filemanager.ConnectWithEnvVariables()
  if err != nil {
    log.Fatal(err)
  }

  downloader := s3filemanager.NewS3BucketDownloader(bucket, folder, logger)
  err = downloader.Run(session)
  ...
```

`S3BucketDownloader` uses [s3manager.Downloader](https://pkg.go.dev/github.com/stripe/aws-go/service/s3/s3manager#Downloader) 
to download the entire bucket, using the `DownloadWithIterator` function to iterate through the files to be downloaded.
 
> This implementation minimizes the memory footprint and the number of open files: the downloaded content is bufferized
> by the `bufferizedFileWriter` implementation of the `io.WriterAt` interface, and the file is created and closed just after the download
> of the S3 object completes, by the `SaveFileAfterWrite` struct that extends `s3manager.BatchDownloadObject`.<br/>

`S3BucketDownloader` can be tested using the provided `main` function with the following options:
```bash
go run main.go -b BUCKET_NAME -debug -f DESTINATION_FOLDER -m download
```
# Using Noobaa
These tools work with all S3-compatible endpoints, and [Noobaa](https://www.noobaa.io/) is a data gateway for objects
that provides the same S3 API and management tools. If you can access a K8s/OpenShift cluster, you don't need an AWS account
to create the S3-compatible storage.

Once you install the `noobaa` [CLI](https://www.noobaa.io/noobaa-operator-cli.html), you can install the NooBaa operator and services on a
given namespace as:
```bash
# Initialize the TARGET_NS
TARGET_NS= ...
noobaa install -n $TARGET_NS
```

You can get quick-start information by running:
```bash
oc descrive noobaa noobaa -n $TARGET_NS
```
and look for the section starting with:
```
  ...
  Readme:               

  Welcome to NooBaa!
  -----------------
  ...
```

To test the `s3-tools` functionality, initialize the S3 connection to the `NooBaa` instance first:
```bash
. ./noobaa-init.sh $TARGET_NS
```
The above script initializes all the environment variables for you.

**Note**: unless you manually specified a value in the field `spec.region`, the default S3 region for NooBaa instances is `us-east-1`

You can use the following commands to start playing with the S3 CLI (install [aws CLI](https://docs.aws.amazon.com/cli/latest/userguide/getting-started-install.html) first):
```bash
# Create s3 alias
alias s3='aws --endpoint "https://$S3_ENDPOINT" --no-verify-ssl s3'

# list objects buckets
s3 ls
# list objects of bucket, recursively
s3 ls --recursive s3://mybucket
# count objects of bucket, recursively
s3 ls --summarize --recursive s3://mybucket

# remove objects from bucket, recursively
s3 rm --recursive s3://mybucket
# remove empty bucket
s3 rb s3://mybucket
# remove bucket and all content
s3 rb s3://mybucket --force
```