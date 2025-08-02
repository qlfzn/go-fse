package crawler

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"path"
	"time"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
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
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	Path      string    `json:"path"`
	Size      int64     `json:"size"`
	Type      string    `json:"type"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type Crawler struct {
	Config Config
	Client *minio.Client
}

func NewCrawler(cfg Config) (*Crawler, error) {
	client, err := minio.New(cfg.MinioEndpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(cfg.AccessKey, cfg.SecretKey, ""),
		Secure: cfg.UseSSL,
	})
	if err != nil {
		return nil, err
	}
	return &Crawler{Config: cfg, Client: client}, nil
}

func (c *Crawler) ListAndSendFiles() error {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	objectCh := c.Client.ListObjects(ctx, c.Config.BucketName, minio.ListObjectsOptions{})
	for object := range objectCh {
		if object.Err != nil {
			log.Printf("error fetching object: %v", object.Err)
			continue
		}

		fileObj := File{
			ID:        object.ETag,
			Name:      path.Base(object.Key),
			Path:      object.Key,
			Size:      object.Size,
			Type:      path.Ext(object.Key),
			CreatedAt: object.LastModified,
			UpdatedAt: object.LastModified,
		}

		data, _ := json.Marshal(fileObj)
		res, err := http.Post(c.Config.IndexerURL+"/index", "application/json", bytes.NewBuffer(data))
		if err != nil {
			log.Printf("Failed to send file to indexer: %v", err)
			continue
		}
		res.Body.Close()
	}
	return nil
}
