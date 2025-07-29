package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"

	"github.com/joho/godotenv"
)

type Config struct {
	MinioEndpoint string
	AccessKey     string
	SecretKey     string
	BucketName    string
	IndexerURL    string
	UseSSL        bool
}

type File struct {
	ID        string
	Name      string
	Path      string
	Size      string
	Type      string
	CreatedAt time.Time
	UpdateAt  time.Time
}

type Crawler struct {
	Config Config
	Client *minio.Client
}

func (c *Crawler) IsBucketExist(bucketName string) error {
	ctx := context.Background()
	found, err := c.Client.BucketExists(ctx, bucketName)
	if err != nil {
		return err
	}

	if found {
		fmt.Printf("Bucket '%s' exists.\n", bucketName)
	} else {
		fmt.Printf("Bucket '%s' does not exist.\n", bucketName)
	}

	return nil
}

func loadConfig() Config {
	err := godotenv.Load()
	if err != nil {
		log.Printf("failed to load env: %s", err)
	}

	return Config{
		MinioEndpoint: "localhost:9000",
		AccessKey:     os.Getenv("MINIO_ROOT_USER"),
		SecretKey:     os.Getenv("MINIO_ROOT_PASSWORD"),
		BucketName:    "ql-bucket",
		UseSSL:        false,
	}
}

func main() {
	cfg := loadConfig()

	// initialise minio client
	minioClient, err := minio.New(cfg.MinioEndpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(cfg.AccessKey, cfg.SecretKey, ""),
		Secure: cfg.UseSSL,
	})
	if err != nil {
		log.Fatalf("Failed to create MinIO client: %v", err)
	}

	crawler := &Crawler{
		Config: cfg,
		Client: minioClient,
	}

	if err := crawler.IsBucketExist("ql-bucket"); err != nil {
		log.Fatalf("Bucket check failed: %v", err)
	}
}
