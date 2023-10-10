package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type SendoraCityClient struct {
	baseUri string
	client  *http.Client
}

func NewClient(baseUri string) *SendoraCityClient {
	return &SendoraCityClient{
		baseUri: baseUri,
		client:  &http.Client{},
	}
}

func (c *SendoraCityClient) doRequest(req *http.Request) (*http.Response, error) {
	resp, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode > http.StatusCreated && resp.StatusCode != http.StatusNotFound {
		return nil, handleError(req, resp)
	}
	return resp, nil
}

func handleError(req *http.Request, resp *http.Response) error {
	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if len(responseBody) != 0 {
		apiError := make(map[string]any)
		if err = json.Unmarshal(responseBody, &apiError); err != nil {
			return NewCustomError(resp.StatusCode, "Method %s to uri %s failed, reason : %s",
				req.Method, req.URL.String(), string(responseBody))
		}
		if errors, ok := apiError["errors"].(map[string]any); ok {
			errorMessage, err := json.Marshal(errors)
			if err != nil {
				return err
			}
			return NewCustomError(resp.StatusCode, "Method %s to uri %s failed, reason : %s",
				req.Method, req.URL.String(), errorMessage)
		}
	}

	return NewCustomError(resp.StatusCode, "Method %s to uri %s failed",
		req.Method, req.URL.String())
}

func (c *SendoraCityClient) DoCreate(url string, body []byte) (*http.Response, error) {
	req, err := http.NewRequest("POST", fmt.Sprintf("%s/%s", c.baseUri, url), bytes.NewBuffer(body))
	if err != nil {
		return nil, err
	}
	req.Header = http.Header{
		"Content-Type": {"application/json"},
		"Accept":       {"*/*"},
	}
	return c.doRequest(req)
}

func (c *SendoraCityClient) DoList(url string, filters map[string]string) (*http.Response, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/%s", c.baseUri, url), nil)
	if err != nil {
		return nil, err
	}
	querry := req.URL.Query()
	for key, value := range filters {
		querry.Add(key, value)
	}
	return c.doRequest(req)
}

func (c *SendoraCityClient) DoRead(url string) (*http.Response, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/%s", c.baseUri, url), nil)
	if err != nil {
		return nil, err
	}
	return c.doRequest(req)
}

func (c *SendoraCityClient) DoUpdate(url string, body []byte) (*http.Response, error) {
	req, err := http.NewRequest("PATCH", fmt.Sprintf("%s/%s", c.baseUri, url), bytes.NewBuffer(body))
	if err != nil {
		return nil, err
	}
	req.Header = http.Header{
		"Content-Type": {"application/json"},
		"Accept":       {"*/*"},
	}
	return c.doRequest(req)
}

func (c *SendoraCityClient) DoDelete(url string) (*http.Response, error) {
	req, err := http.NewRequest("DELETE", fmt.Sprintf("%s/%s", c.baseUri, url), nil)
	if err != nil {
		return nil, err
	}
	return c.doRequest(req)
}
