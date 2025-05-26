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
	return c.doRequest(http.MethodPost, url, request, headers)
}

func (c *GenericHttpClient) Get(url string, headers map[string]string) (*json.RawMessage, error) {
	return c.doRequest(http.MethodGet, url, nil, headers)
}

func (c *GenericHttpClient) Put(url string, request interface{}, headers map[string]string) (*json.RawMessage, error) {
	return c.doRequest(http.MethodPut, url, request, headers)
}

func (c *GenericHttpClient) Patch(url string, request interface{}, headers map[string]string) (*json.RawMessage, error) {
	return c.doRequest(http.MethodPatch, url, request, headers)
}

func (c *GenericHttpClient) doRequest(method, url string, request interface{}, headers map[string]string) (*json.RawMessage, error) {
	c.logger.Info("Making request", zap.String("method", method), zap.String("url", url))

	var body io.Reader
	if request != nil {
		jsonBytes, err := json.Marshal(request)
		if err != nil {
			return nil, err
		}
		body = bytes.NewBuffer(jsonBytes)
	}

	req, err := retryablehttp.NewRequest(method, url, body)
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
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		return nil, errors.New("HTTP error: " + resp.Status)
	}

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var raw json.RawMessage
	if err := json.Unmarshal(respBody, &raw); err != nil {
		return nil, err
	}

	formatted, err := json.MarshalIndent(raw, "", "  ")
	if err != nil {
		c.logger.Error("Failed to format response", zap.Error(err))
	} else {
		c.logger.Info("Received response", zap.String("response", string(formatted)))
	}

	return &raw, nil
}
