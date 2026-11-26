package services

import (
	"context"
	"errors"
	"fmt"
	"log"
	"mime/multipart"
	"net/url"
	"strings"
)

func (s Service) ServicesUploadSongChart(ctx context.Context, file multipart.File, fileID, bandSlug, titleSlug, contentType string) (string, error) {
	log.Println("- ServicesUploadSongChart")

	if fileID == "" {
		return "", errors.New("ServicesUploadSongChart - file ID is empty")
	}
	bandSlug = strings.TrimSpace(bandSlug)
	if bandSlug == "" {
		return "", errors.New("ServicesUploadSongChart - band slug is empty")
	}
	titleSlug = strings.TrimSpace(titleSlug)
	if titleSlug == "" {
		return "", errors.New("ServicesUploadSongChart - title slug is empty")
	}
	if contentType == "" {
		contentType = "application/octet-stream"
	}

	destinationKey := fmt.Sprintf(
		"%s/charts/%s/%s",
		bandSlug,
		titleSlug,
		fileID,
	)

	publicURL, err := s.Storage.Upload(
		ctx,
		destinationKey,
		file,
		contentType,
	)

	if err != nil {
		return "", fmt.Errorf("ServicesUploadSongChart - R2 copy %s, %s image: %w", bandSlug, fileID, err)
	}

	return publicURL, nil
}

func (s Service) ServicesDeleteSongChart(ctx context.Context, chartURL string) error {
	chartURL = strings.TrimSpace(chartURL)
	if chartURL == "" {
		return nil
	}

	publicURL, err := url.Parse(s.Storage.PublicURL)
	if err != nil {
		return fmt.Errorf("parse R2 public URL: %w", err)
	}
	storedURL, err := url.Parse(chartURL)
	if err != nil {
		return fmt.Errorf("parse stored chart URL: %w", err)
	}
	if storedURL.Scheme != publicURL.Scheme || storedURL.Host != publicURL.Host {
		return errors.New("stored chart URL does not belong to configured R2 storage")
	}

	basePath := strings.TrimSuffix(publicURL.Path, "/")
	pathPrefix := basePath + "/"
	if !strings.HasPrefix(storedURL.Path, pathPrefix) {
		return errors.New("stored chart URL is outside configured R2 storage path")
	}

	key := strings.TrimPrefix(storedURL.Path, pathPrefix)
	if key == "" {
		return errors.New("stored chart URL has no object key")
	}

	if err := s.Storage.Delete(ctx, key); err != nil {
		return fmt.Errorf("delete song chart from R2: %w", err)
	}
	return nil
}
