package minio

import (
	"context"
	"fmt"
	"io"
	"mime/multipart"

	"github.com/minio/minio-go/v7"
)

type Client struct {
	mcl      *minio.Client
	bucket   string
	endpoint string
}

func New(mcl *minio.Client, bucket string, endpoint string) *Client {
	return &Client{
		mcl:      mcl,
		bucket:   bucket,
		endpoint: endpoint,
	}
}

func (c *Client) Endpoint() string {
	return c.endpoint
}

func (c *Client) UploadPhoto(ctx context.Context, file *multipart.FileHeader) (string, error) {
	src, err := file.Open()
	if err != nil {
		return "", fmt.Errorf("file.Open: %w", err)
	}
	defer src.Close()

	return c.UploadFromReader(ctx, src, file.Filename, file.Size, file.Header.Get("Content-Type"))
}

func (c *Client) GetObject(ctx context.Context, objectName string) (*minio.Object, error) {
	return c.mcl.GetObject(ctx, c.bucket, objectName, minio.GetObjectOptions{})
}

func (c *Client) DeleteObject(ctx context.Context, objectName string) error {
	return c.mcl.RemoveObject(ctx, c.bucket, objectName, minio.RemoveObjectOptions{})
}

func (c *Client) UploadFromReader(ctx context.Context, reader io.Reader, filename string, size int64, contentType string) (string, error) {
	info, err := c.mcl.PutObject(
		ctx,
		c.bucket,
		filename,
		reader,
		size,
		minio.PutObjectOptions{
			ContentType: contentType,
		},
	)
	if err != nil {
		return "", fmt.Errorf("c.mcl.PutObject: %w", err)
	}

	return fmt.Sprintf("http://%s/%s/%s", c.endpoint, c.bucket, info.Key), nil
}
