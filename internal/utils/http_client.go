package utils

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

// Should this interface live in the generic interfaces directory?
type HTTPClient interface {
	Get(ctx context.Context, url string) (*HTTPResponse, error)
}

// TODO: Rename this response per proper domain (e.g. AgentValidationResponse?)
type HTTPResponse struct {
	GUID string `json:"guid"`
}

// Should ValidationClient live in it's own file (e.g. validation_client.go)
type ValidationClient struct{}

func NewValidationClient() *ValidationClient {
	return &ValidationClient{}
}

func (c *ValidationClient) Get(ctx context.Context, url string) (*HTTPResponse, error) {
	httpClient := &http.Client{}

	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, err
	}

	resp, err := httpClient.Do(req)
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

	log.Print("\n\n **************************** \n")
	log.Printf("\n RESPONSE:  %+v \n", response)
	log.Print("\n **************************** \n\n")
	time.Sleep(7 * time.Second)

	return &response, nil

	// return &HTTPResponse{
	// 	EntityGUID: "MTA5ODI2NzB8SU5GUkF8TkF8Nzc0NTc3NjI1NjYyNzI5NzYzNw",
	// 	EntityID:   "7745776256627297637",
	// 	EntityKey:  "i-00025a7a6b7bf0a4c",
	// }
}
