package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"time"
)

const (
	BasePath         = "https://napi.arvancloud.com"
	AuthEndpoint     = "/resid/v1/wallets/me"
	DefaultTimeout   = 1 * time.Minute
	RequesterContext = "requesterContext"
)

type ErrorResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

type SuccessResponse struct {
	Message string `json:"message,omitempty"`
	Data    any    `json:"data,omitempty"`
}

type Requester struct {
	ApiKey     string
	HttpClient *http.Client
}

// NewRequester - using this for generating our requests
func NewRequester(apiKey string) *Requester {
	return &Requester{
		HttpClient: &http.Client{
			Timeout: DefaultTimeout,
		},
		ApiKey: apiKey,
	}
}

// CheckAuthenticate - to check authenticate
func (r *Requester) CheckAuthenticate() (err error) {
	_, err = r.DoRequest("GET", AuthEndpoint, nil)
	if err != nil {
		return err
	}
	return nil
}

// DoRequest - make a request without queries
func (r *Requester) DoRequest(method, endpoint string, data io.Reader) ([]byte, error) {
	return r.DoRequestWithQuery(method, endpoint, data, nil)
}

// DoRequestWithQuery - make a request via its queries
func (r *Requester) DoRequestWithQuery(method, endpoint string, data io.Reader,
	query map[string]string) ([]byte, error) {
	url := fmt.Sprintf("%v/%v", BasePath, endpoint)

	request, err := http.NewRequest(method, url, data)
	if err != nil {
		return nil, err
	}

	request.Header.Set("Authorization", r.ApiKey)
	request.Header.Set("Content-Type", "application/json")

	if query != nil {
		q := request.URL.Query()
		for key, value := range query {
			q.Add(key, value)
		}
		request.URL.RawQuery = q.Encode()
	}

	res, err := r.HttpClient.Do(request)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	if !(res.StatusCode <= http.StatusBadRequest) {
		return nil, fmt.Errorf("status: %d, body: %s", res.StatusCode, string(body))
	}

	return body, err
}

// Create - make a POST request
func (r *Requester) Create(endpoint string, opts any, queries map[string]string) (any, error) {
	return r.Custom("POST", endpoint, opts, queries)
}

// Read - make a GET Request
func (r *Requester) Read(endpoint string, queries map[string]string) (any, error) {
	return r.Custom("GET", endpoint, nil, queries)
}

// Patch - make a PATCH request
func (r *Requester) Patch(endpoint string, opts any, queries map[string]string) (any, error) {
	return r.Custom("PATCH", endpoint, opts, queries)
}

// Put - make PUT request
func (r *Requester) Put(endpoint string, opts any, queries map[string]string) (any, error) {
	return r.Custom("PUT", endpoint, opts, queries)
}

// Delete - make DELETE request
func (r *Requester) Delete(endpoint string, queries map[string]string) (err error) {
	if queries == nil {
		_, err = r.DoRequest("DELETE", endpoint, nil)
	} else {
		_, err = r.DoRequestWithQuery("DELETE", endpoint, nil, queries)
	}
	return err
}

// List - like Read but for List purpose
func (r *Requester) List(endpoint string, queries map[string]string) (any, error) {
	return r.Read(endpoint, queries)
}

// Custom - general purpose maker for request
func (r *Requester) Custom(method, endpoint string, opts any, queries map[string]string) (details any, err error) {
	var response, body []byte

	if opts != nil {
		body, err = json.Marshal(opts)
		if err != nil {
			return nil, err
		}
	}

	if queries != nil {
		response, err = r.DoRequestWithQuery(method, endpoint, bytes.NewBuffer(body), queries)
	} else {
		response, err = r.DoRequest(method, endpoint, bytes.NewBuffer(body))
	}
	if err != nil {
		return nil, err
	}

	var successResponse *SuccessResponse
	err = json.Unmarshal(response, &successResponse)
	if err != nil {
		return nil, err
	}

	data, err := json.Marshal(successResponse.Data)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(data, &details)
	if err != nil {
		return nil, err
	}

	return details, nil
}
