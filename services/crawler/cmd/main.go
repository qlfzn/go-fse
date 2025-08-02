package main

import (
	"log"
	"os"

	"github.com/joho/godotenv"
	crawler "github.com/qlfzn/go-fse/services/crawler/internal"
)

func main() {
	godotenv.Load()

	cfg := crawler.Config{
		MinioEndpoint: "localhost:9000",
		AccessKey:     os.Getenv("MINIO_ROOT_USER"),
		SecretKey:     os.Getenv("MINIO_ROOT_PASSWORD"),
		BucketName:    os.Getenv("MINIO_BUCKET_NAME"),
		IndexerURL:    os.Getenv("INDEXER_URL"),
		UseSSL:        false,
	}

	cr, err := crawler.NewCrawler(cfg)
	if err != nil {
		log.Fatalf("failed to init crawler: %v", err)
	}

	if err := cr.ListAndSendFiles(); err != nil {
		log.Fatalf("crawler failed: %v", err)
	}
}
