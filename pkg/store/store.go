package store

import (
	"fmt"

	"github.com/chartmuseum/storage"
)

type (
	// Uploader --
	Uploader struct {
		Backend storage.Backend
	}
)

// NewUploader --
func NewUploader(cloud, bucket, prefix, region, endpoint, sse string) (*Uploader, error) {
	var backend storage.Backend
	switch cloud {
	case "local":
		backend = storage.NewLocalFilesystemBackend(prefix)
	case "aws":
		backend = storage.NewAmazonS3Backend(bucket, prefix, region, endpoint, sse)
	case "azure":
		backend = storage.NewMicrosoftBlobBackend(bucket, prefix)
	case "google":
		backend = storage.NewGoogleCSBackend(bucket, prefix)
	default:
		return nil, fmt.Errorf("cloud provider " + cloud + " not supported")
	}
	uploader := Uploader{Backend: backend}
	return &uploader, nil
}

// Put --
func (uploader *Uploader) Put(filename string, contents []byte) error {
	if err := uploader.Backend.PutObject(filename, contents); err != nil {
		return err
	}
	return nil
}

// Get --
func (uploader *Uploader) Get(filename string) ([]byte, error) {
	obj, err := uploader.Backend.GetObject(filename)
	if err != nil {
		return nil, err
	}

	return obj.Content, nil
}
