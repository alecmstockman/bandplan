package storage

import (
	"context"
	"errors"
	"fmt"
	"io"
	"log"
	"net/url"
	"os"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

type R2Storage struct {
	Client    *s3.Client
	Bucket    string
	PublicURL string
}

func NewR2Storage(ctx context.Context) (*R2Storage, error) {
	accessKey := os.Getenv("R2_ACCESS_KEY")
	secretKey := os.Getenv("R2_SECRET_KEY")
	bucketName := os.Getenv("R2_BUCKET_NAME")
	endpoint := os.Getenv("R2_ENDPOINT")
	publicURL := os.Getenv("R2_PUBLIC_URL")

	if accessKey == "" {
		return nil, errors.New("R2_ACCESS_KEY is required")
	}

	if secretKey == "" {
		return nil, errors.New("R2_SECRET_KEY is required")
	}

	if bucketName == "" {
		return nil, errors.New("R2_BUCKET_NAME is required")
	}

	if endpoint == "" {
		return nil, errors.New("R2_ENDPOINT is required")
	}

	if publicURL == "" {
		return nil, errors.New("R2_PUBLIC_URL is required")
	}

	cfg, err := config.LoadDefaultConfig(
		ctx,
		config.WithRegion("auto"),
		config.WithCredentialsProvider(
			credentials.NewStaticCredentialsProvider(
				accessKey,
				secretKey,
				"",
			),
		),
	)
	if err != nil {
		return nil, err
	}

	client := s3.NewFromConfig(cfg, func(options *s3.Options) {
		options.BaseEndpoint = aws.String(endpoint)
		options.UsePathStyle = true
	})

	return &R2Storage{
		Client:    client,
		Bucket:    bucketName,
		PublicURL: strings.TrimRight(publicURL, "/"),
	}, nil
}

func (storage *R2Storage) Upload(
	ctx context.Context,
	key string,
	body io.Reader,
	contentType string,
) (string, error) {
	key = strings.TrimLeft(key, "/")

	if key == "" {
		return "", errors.New("R2 object key is required")
	}

	if body == nil {
		return "", errors.New("R2 upoad body is required")
	}

	_, err := storage.Client.PutObject(ctx, &s3.PutObjectInput{
		Bucket:      aws.String(storage.Bucket),
		Key:         aws.String(key),
		Body:        body,
		ContentType: aws.String(contentType),
	})

	if err != nil {
		return "", fmt.Errorf("uploadR2 object %q: %w", key, err)
	}

	publicURL, err := url.JoinPath(storage.PublicURL, key)
	if err != nil {
		return "", fmt.Errorf("build public URL for %q: %w", key, err)
	}

	return publicURL, nil
}

func (storage *R2Storage) Delete(ctx context.Context, key string) error {
	log.Println("- R2 Delete")

	if key == "" {
		return errors.New("R2 object key is required")
	}

	_, err := storage.Client.DeleteObject(
		ctx,
		&s3.DeleteObjectInput{
			Bucket: aws.String(storage.Bucket),
			Key:    aws.String(key),
		},
	)
	if err != nil {
		log.Println("   Unable to delete Song Artwork from R2: ", err)
		return errors.New("Unable to delete Song Artwork from R2")
	}

	log.Print("   Successfully deleted image artwork form R2")

	return nil
}
