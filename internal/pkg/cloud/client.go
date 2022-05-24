package cloud

import (
	"context"
	"io"
)

type BucketClient interface {
	// Creates a new bucket.
	Create(ctx context.Context, bucket string) error
	// Upload a new object to a bucket and returns its URL to view/download.
	UploadObject(ctx context.Context, bucket, fileName string, body io.Reader) (string, error)
	// Downloads an existing object from a bucket.
	DownloadObject(ctx context.Context, bucket, fileName string, body io.WriterAt) error
	// Deletes an existing object from a bucket.
	DeleteObject(ctx context.Context, bucket, fileName string) error
	// Lists all objects in a bucket.
	ListObjects(ctx context.Context, bucket string) ([]*Object, error)
	// Returns an object from bucket for reading.
	FetchObject(ctx context.Context, bucket, fileName string) (io.ReadCloser, error)
}
