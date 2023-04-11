package s3filemanager

import (
	"os"
	"path/filepath"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/sirupsen/logrus"
)

type S3BucketDownloader struct {
	bucket string
	folder string
	logger *logrus.Logger
}

func NewS3BucketDownloader(bucket string, folder string, logger *logrus.Logger) *S3BucketDownloader {
	return &S3BucketDownloader{bucket: bucket, folder: folder, logger: logger}
}

func (s3bd *S3BucketDownloader) Run(sess *session.Session) error {
	s3bd.logger.Infof("Started downloading bucket %s to %s", s3bd.bucket, s3bd.folder)
	err := s3bd.downloadToFolder(sess)
	if err != nil {
		s3bd.logger.Infof("Download failed: %s", err)
	}
	return err
}

type SaveFileAfterWrite struct {
	s3manager.BatchDownloadObject
	logger *logrus.Logger
}

func NewSaveFileAfterWrite(object *s3.Object, bucket string, folder string, logger *logrus.Logger) SaveFileAfterWrite {
	batchDownloadObject := SaveFileAfterWrite{}
	batchDownloadObject.Object = &s3.GetObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(*object.Key),
	}
	batchDownloadObject.Writer = NewWriterForObject(object, folder, logger)
	batchDownloadObject.After = batchDownloadObject.saveAfterDownload
	batchDownloadObject.logger = logger
	return batchDownloadObject
}

func (o *SaveFileAfterWrite) saveAfterDownload() error {
	n, err := o.BatchDownloadObject.Writer.(*bufferizedFileWriter).save()
	if n < 1 {
		o.logger.Warnf("No data saved for %s", o.BatchDownloadObject.Writer.(*bufferizedFileWriter).path)
	}
	return err
}

type bufferizedFileWriter struct {
	path string

	buf    *aws.WriteAtBuffer
	logger *logrus.Logger
}

func NewWriterForObject(object *s3.Object, basedir string, logger *logrus.Logger) *bufferizedFileWriter {
	writer := &bufferizedFileWriter{logger: logger}
	writer.path = filepath.Join(basedir, *object.Key)

	return writer
}

func (w *bufferizedFileWriter) WriteAt(p []byte, off int64) (n int, err error) {
	w.logger.Debugf("Called WriteAt for %s", w.path)
	if w.buf == nil {
		w.buf = &aws.WriteAtBuffer{}
	}
	return w.buf.WriteAt(p, off)
}

func (w *bufferizedFileWriter) save() (n int, err error) {
	w.logger.Debugf("Saving %s", w.path)
	if w.buf == nil {
		w.logger.Warnf("Buffer empty")
		return 0, nil
	}
	var f *os.File
	err = os.MkdirAll(filepath.Dir(w.path), 0775)
	if err != nil {
		w.logger.Errorf("Cannot create base folder %s: %s", filepath.Dir(w.path), err)
		return 0, err
	}
	w.logger.Debugf("Created base folder %s", filepath.Dir(w.path))

	f, err = os.Create(w.path)
	if err != nil {
		w.logger.Errorf("Cannot create file %s: %s", w.path, err)
		return 0, err
	}

	defer f.Close()
	w.logger.Debugf("Open files: %d", countOpenFiles())
	return f.Write(w.buf.Bytes())
}

func (s3bd *S3BucketDownloader) downloadToFolder(sess *session.Session) error {
	client := s3.New(sess)

	params := &s3.ListObjectsInput{
		Bucket: aws.String(s3bd.bucket),
	}

	objects := []s3manager.BatchDownloadObject{}
	resp, _ := client.ListObjects(params)
	for _, object := range resp.Contents {
		objects = append(objects, NewSaveFileAfterWrite(object, s3bd.bucket, s3bd.folder, s3bd.logger).BatchDownloadObject)
	}

	iter := &s3manager.DownloadObjectsIterator{Objects: objects}
	downloader := s3manager.NewDownloader(sess)
	err := downloader.DownloadWithIterator(aws.BackgroundContext(), iter)

	s3bd.logger.Infof("Successfully downloaded %d files from folder %s to bucket %s", len(resp.Contents), s3bd.folder, s3bd.bucket)
	return err
}
