package aws

import (
	"context"
	"fmt"
	"io"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/luke92/LocalStackCloudTestGolang/internal/pkg/cloud"
)

var _ cloud.BucketClient = S3{}

type S3 struct {
	timeout    time.Duration
	client     *s3.S3
	uploader   *s3manager.Uploader
	downloader *s3manager.Downloader
}

func NewS3(session *session.Session, timeout time.Duration) S3 {
	s3manager.NewUploader(session)
	return S3{
		timeout:    timeout,
		client:     s3.New(session),
		uploader:   s3manager.NewUploader(session),
		downloader: s3manager.NewDownloader(session),
	}
}

func (s S3) Create(ctx context.Context, bucket string) error {
	ctx, cancel := context.WithTimeout(ctx, s.timeout)
	defer cancel()

	if _, err := s.client.CreateBucketWithContext(ctx, &s3.CreateBucketInput{
		Bucket: aws.String(bucket),
	}); err != nil {
		return fmt.Errorf("create: %w", err)
	}

	if err := s.client.WaitUntilBucketExists(&s3.HeadBucketInput{
		Bucket: aws.String(bucket),
	}); err != nil {
		return fmt.Errorf("wait: %w", err)
	}

	return nil
}

func (s S3) Exists(ctx context.Context, bucket string) (bool, error) {
	ctx, cancel := context.WithTimeout(ctx, s.timeout)
	defer cancel()

	_, err := s.client.HeadBucketWithContext(ctx, &s3.HeadBucketInput{
		Bucket: aws.String(bucket),
	})

	return getBooleanAndError(err)
}

func (s S3) ExistsObject(ctx context.Context, bucket, fileName string) (bool, error) {
	ctx, cancel := context.WithTimeout(ctx, s.timeout)
	defer cancel()

	_, err := s.client.HeadObjectWithContext(ctx, &s3.HeadObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(fileName),
	})
	return getBooleanAndError(err)
}

func getBooleanAndError(err error) (bool, error) {
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			case "NotFound": // s3.ErrCodeNoSuchKey does not work, aws is missing this error code so we hardwire a string
				return false, nil
			default:
				return false, err
			}
		}
		return false, err
	}
	return true, nil
}

func (s S3) UploadObject(ctx context.Context, bucket, fileName string, body io.Reader) (string, error) {
	ctx, cancel := context.WithTimeout(ctx, s.timeout)
	defer cancel()

	res, err := s.uploader.UploadWithContext(ctx, &s3manager.UploadInput{
		Body:   body,
		Bucket: aws.String(bucket),
		Key:    aws.String(fileName),
	})
	if err != nil {
		return "", fmt.Errorf("upload: %w", err)
	}

	return res.Location, nil
}

func (s S3) DownloadObject(ctx context.Context, bucket, fileName string, body io.WriterAt) error {
	if _, err := s.downloader.DownloadWithContext(ctx, body, &s3.GetObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(fileName),
	}); err != nil {
		return fmt.Errorf("download: %w", err)
	}

	return nil
}

func (s S3) DeleteObject(ctx context.Context, bucket, fileName string) error {
	if _, err := s.client.DeleteObjectWithContext(ctx, &s3.DeleteObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(fileName),
	}); err != nil {
		return fmt.Errorf("delete: %w", err)
	}

	if err := s.client.WaitUntilObjectNotExists(&s3.HeadObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(fileName),
	}); err != nil {
		return fmt.Errorf("wait: %w", err)
	}

	return nil
}

func (s S3) ListObjects(ctx context.Context, bucket string) ([]*cloud.Object, error) {
	ctx, cancel := context.WithTimeout(ctx, s.timeout)
	defer cancel()

	res, err := s.client.ListObjectsV2WithContext(ctx, &s3.ListObjectsV2Input{
		Bucket: aws.String(bucket),
	})
	if err != nil {
		return nil, fmt.Errorf("list: %w", err)
	}

	objects := make([]*cloud.Object, len(res.Contents))

	for i, object := range res.Contents {
		objects[i] = &cloud.Object{
			Key:        *object.Key,
			Size:       *object.Size,
			ModifiedAt: *object.LastModified,
		}
	}

	return objects, nil
}

func (s S3) FetchObject(ctx context.Context, bucket, fileName string) (io.ReadCloser, error) {
	ctx, cancel := context.WithTimeout(ctx, s.timeout)
	defer cancel()

	res, err := s.client.GetObjectWithContext(ctx, &s3.GetObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(fileName),
	})
	if err != nil {
		return nil, err
	}

	return res.Body, nil
}
