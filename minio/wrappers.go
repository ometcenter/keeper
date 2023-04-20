package minio

import (
	"os"
	"strings"
)

func MustGetClient(c *Config) *Client {
	minioClient, err := NewClientFromConfig(c)
	if err != nil {
		panic(err)
	}

	return minioClient
}

func PublicConfig() *Config {
	return &Config{
		AccessKey:       os.Getenv("MINIO_ACCESS_KEY_ID"),
		SecretAccessKey: os.Getenv("MINIO_SECRET_ACCESS_KEY"),
		Bucket:          os.Getenv("MINIO_BUCKET"),
		UseSSL:          strings.EqualFold(os.Getenv("MINIO_USE_SSL"), "true"),
		Region:          os.Getenv("MINIO_REGION"), // location
		Endpoint:        os.Getenv("MINIO_ENDPOINT"),
	}
}
