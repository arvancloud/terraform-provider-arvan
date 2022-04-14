package api

import (
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"time"
)

const (
	BaseUrl        = "https://napi.arvancloud.com"
	DefaultTimeout = 1 * time.Minute
)

type ErrorResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

type SuccessResponse struct {
	Code int         `json:"code"`
	Data interface{} `json:"data"`
}

type Requester struct {
	ApiKey     string
	HttpClient *http.Client
}

func NewRequester(apiKey string) *Requester {
	return &Requester{
		HttpClient: &http.Client{
			Timeout: DefaultTimeout,
		},
		ApiKey: apiKey,
	}
}

func (r *Requester) DoRequest(method, endpoint string, data io.Reader) ([]byte, error) {
	return r.DoRequestWithQuery(method, endpoint, data, nil)
}

func (r *Requester) DoRequestWithQuery(method, endpoint string, data io.Reader, query map[string]string) ([]byte, error) {
	url := fmt.Sprintf("%v/%v", BaseUrl, endpoint)

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

	if !(http.StatusOK <= res.StatusCode && res.StatusCode <= http.StatusMultipleChoices) {
		return nil, fmt.Errorf("status: %d, body: %s", res.StatusCode, string(body))
	}

	return body, err
}
