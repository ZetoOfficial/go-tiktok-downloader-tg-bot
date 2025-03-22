package clients

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/ZetoOfficial/go-tiktok-downloader-tg-bot/config"
	"github.com/ZetoOfficial/go-tiktok-downloader-tg-bot/internal/models"
	"github.com/go-resty/resty/v2"
)

type DouyinClient struct {
	client *resty.Client
	apiURL *url.URL
}

func NewDouyinClient(cfg *config.Config) (*DouyinClient, error) {
	parsedURL, err := url.Parse(cfg.DouyinAPI)
	if err != nil {
		return nil, fmt.Errorf("parse url: %w", err)
	}

	client := resty.New().
		SetBaseURL(parsedURL.String()).
		SetTimeout(30 * time.Second).
		SetRetryCount(3)

	return &DouyinClient{
		client: client,
		apiURL: parsedURL,
	}, nil
}

func (c *DouyinClient) Download(ctx context.Context, link string, options ...models.DownloadOption) (*models.DownloadResponse, error) {
	if link == "" {
		return nil, fmt.Errorf("empty link")
	}
	params := map[string]string{
		"url":            link,
		"prefix":         "false",
		"with_watermark": "false",
	}
	for _, opt := range options {
		params[opt.Key] = fmt.Sprintf("%t", opt.Value)
	}
	resp, err := c.client.R().
		SetContext(ctx).
		SetQueryParams(params).
		Get("/download")
	if err != nil {
		return nil, fmt.Errorf("request error: %w", err)
	}
	if resp == nil {
		return nil, fmt.Errorf("empty response from server")
	}
	if resp.IsError() {
		return nil, fmt.Errorf("api error: %s - %s", resp.Status(), string(resp.Body()))
	}
	if resp.StatusCode() != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d - %s", resp.StatusCode(), string(resp.Body()))
	}

	contentType := resp.Header().Get("Content-Type")
	body := resp.Body()
	if len(body) == 0 {
		return nil, fmt.Errorf("empty response body")
	}

	filename := extractFileName(resp.Header().Get("Content-Disposition"))
	if filename == "" {
		filename = "downloaded_file"
	}

	contentLength := resp.Header().Get("Content-Length")
	if contentLength != "" {
		expectedSize, parseErr := parseContentLength(contentLength)
		if parseErr != nil {
			log.Printf("Warning: failed to parse Content-Length: %v", parseErr)
		} else if len(body) < expectedSize {
			return nil, fmt.Errorf("incomplete download: expected %d bytes, got %d", expectedSize, len(body))
		}
	}

	if contentType == "application/json" {
		var apiErr models.APIError
		if err := json.Unmarshal(body, &apiErr); err != nil {
			return nil, fmt.Errorf("invalid JSON error response: %w", err)
		}
		return nil, fmt.Errorf("server error: %s", apiErr.Message)
	}

	if contentType == "application/zip" || strings.HasSuffix(filename, ".zip") {
		photos, err := extractPhotosFromZip(body)
		if err != nil {
			return nil, fmt.Errorf("extract zip: %w", err)
		}
		if len(photos) == 0 {
			return nil, fmt.Errorf("zip contains no image files")
		}
		return &models.DownloadResponse{
			FileName: filename,
			Photos:   photos,
		}, nil
	}

	return &models.DownloadResponse{
		FileName: filename,
		Data:     body,
	}, nil
}
