package minio

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"io/ioutil"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

type Client struct {
	client   *minio.Client
	bucket   string
	endpoint string
}

type Config struct {
	AccessKey       string `json:"accessKey"`
	SecretAccessKey string `json:"secretAccessKey"`
	Region          string `json:"region"`
	Endpoint        string `json:"host"`
	Bucket          string `json:"bucket"`
	UseSSL          bool   `json:"disableSSL"`
}

var MinioClient *Client

// NewClient return new instance of the s3 client with session.
func NewClient(accessKey, secretAccessKey, Endpoint, bucket, region string,
	useSSL bool) (*Client, error) {

	// Initialize minio client object.
	minioClient, err := minio.New(Endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(accessKey, secretAccessKey, ""),
		Secure: useSSL,
		Region: region,
	})
	if err != nil {
		return nil, err
	}
	return &Client{client: minioClient,
		bucket:   bucket,
		endpoint: Endpoint}, nil
}

func NewClientFromConfig(config *Config) (*Client, error) {
	return NewClient(config.AccessKey,
		config.SecretAccessKey,
		config.Endpoint,
		config.Bucket,
		config.Region,
		config.UseSSL,
	)
}

// Put - puts bytes by prefix and path with filename in minio client's bucket.
func (c *Client) Put(ctx context.Context, b []byte, path, fileName string) (*minio.UploadInfo, error) {

	reader := bytes.NewReader(b)
	uploadInfo, err := c.client.PutObject(ctx, c.bucket, fileName, reader, -1, minio.PutObjectOptions{ContentType: "application/octet-stream"})
	if err != nil {
		return nil, err
	}

	return &uploadInfo, nil
}

// PutByReader - puts data from io.ReadSeeker by prefix and path with filename in minio client's bucket.
func (c *Client) PutByReader(ctx context.Context, reader io.ReadSeeker, fileName string) (*minio.UploadInfo, error) {

	uploadInfo, err := c.client.PutObject(ctx, c.bucket, fileName, reader, -1, minio.PutObjectOptions{ContentType: "application/octet-stream"})
	if err != nil {
		return nil, err
	}

	return &uploadInfo, nil
}

// Get - returns bytes from minio host by prefix with path and filename.
func (c *Client) Get(ctx context.Context, path, fileName string) ([]byte, error) {
	out, err := c.GetReadCloser(ctx, path, fileName)
	if err != nil {
		return nil, err
	}

	defer out.Close()

	return ioutil.ReadAll(out)
}

// GetByReadCloser - returns io.ReadCloser from minio host by prefix with path and filename.
func (c *Client) GetReadCloser(ctx context.Context, path, fileName string) (io.ReadCloser, error) {
	out, err := c.client.GetObject(ctx, c.bucket, fileName, minio.GetObjectOptions{})
	if err != nil {
		err = fmt.Errorf("getting object %s from bucket %s failed. Error: %v",
			fileName, c.bucket, err)
		return nil, err
	}

	return out, nil
}

// Delete - deletes object from s3 host by prefix with path and filename.
func (c *Client) Delete(ctx context.Context, path, fileName string) error {
	opts := minio.RemoveObjectOptions{
		GovernanceBypass: true,
		VersionID:        "1",
	}

	if err := c.client.RemoveObject(ctx, c.bucket, fileName, opts); err != nil {
		return err
	}

	return nil

}
