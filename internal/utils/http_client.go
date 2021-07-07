package utils

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
)

type HTTPResponse struct {
	EntityGUID string `json:"entityGuid"`
	EntityID   string `json:"entityId"`
	EntityKey  string `json:"entityKey"`
}

type HTTPClient interface {
	Get(url string) (*HTTPResponse, error)
}

type ValidationClient struct{}

func NewValidationClient() *ValidationClient {
	return &ValidationClient{}
}

func (c *ValidationClient) Get(url string) (*HTTPResponse, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}

	if resp.Body != nil {
		defer resp.Body.Close()
	}

	data, _ := ioutil.ReadAll(resp.Body)

	response := HTTPResponse{}
	err = json.Unmarshal(data, &response)
	if err != nil {
		return nil, err
	}

	return &response, nil

	// return &HTTPResponse{
	// 	EntityGUID: "MTA5ODI2NzB8SU5GUkF8TkF8Nzc0NTc3NjI1NjYyNzI5NzYzNw",
	// 	EntityID:   "7745776256627297637",
	// 	EntityKey:  "i-00025a7a6b7bf0a4c",
	// }
}
