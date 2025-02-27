package utils

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"time"

	"github.com/hashicorp/go-retryablehttp"
	"go.uber.org/zap"
)

type GenericHttpClient struct {
	client *retryablehttp.Client
	logger *zap.Logger
}

func NewGenericHttpClient(logger *zap.Logger) *GenericHttpClient {
	client := retryablehttp.NewClient()
	client.RetryMax = 3
	client.RetryWaitMin = 10 * time.Second
	client.RetryWaitMax = 30 * time.Second
	client.Logger = nil // Suppress internal logging

	return &GenericHttpClient{
		client: client,
		logger: logger,
	}
}

func (c *GenericHttpClient) Post(url string, request interface{}, headers map[string]string) (*json.RawMessage, error) {
	c.logger.Info("Making POST request", zap.String("url", url))

	body, err := json.Marshal(request)
	if err != nil {
		return nil, err
	}

	req, err := retryablehttp.NewRequest(http.MethodPost, url, bytes.NewBuffer(body))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")
	for key, value := range headers {
		req.Header.Set(key, value)
	}

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {

		}
	}(resp.Body)

	if resp.StatusCode >= 400 {
		return nil, errors.New("HTTP error: " + resp.Status)
	}

	var response json.RawMessage
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if err := json.Unmarshal(respBody, &response); err != nil {
		return nil, err
	}

	formattedResponse, err := json.MarshalIndent(response, "", "  ")
	if err != nil {
		c.logger.Error("Failed to format response", zap.Error(err))
	} else {
		c.logger.Info("Received response", zap.String("response", string(formattedResponse)))
	}

	return &response, nil
}
