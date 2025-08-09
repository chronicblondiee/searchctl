package datastreams

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/chronicblondiee/searchctl/pkg/client/rest"
	"github.com/chronicblondiee/searchctl/pkg/types"
)

type client struct {
	restClient *rest.Client
}

func New(restClient *rest.Client) Interface {
	return &client{
		restClient: restClient,
	}
}

func (c *client) List(pattern string) ([]types.DataStream, error) {
	dataStreamPattern := "*"
	if pattern != "" {
		dataStreamPattern = pattern
	}

	path := fmt.Sprintf("/_data_stream/%s", dataStreamPattern)
	resp, err := c.restClient.Get(path)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("error getting data streams: %s", string(resp.Body))
	}

	var response struct {
		DataStreams []types.DataStream `json:"data_streams"`
	}
	if err := json.Unmarshal(resp.Body, &response); err != nil {
		return nil, err
	}

	return response.DataStreams, nil
}

func (c *client) Create(name string) error {
	path := fmt.Sprintf("/_data_stream/%s", name)
	resp, err := c.restClient.Put(path, nil)
	if err != nil {
		return err
	}

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		return fmt.Errorf("error creating data stream: %s", string(resp.Body))
	}

	return nil
}

func (c *client) Delete(name string) error {
	path := fmt.Sprintf("/_data_stream/%s", name)
	resp, err := c.restClient.Delete(path)
	if err != nil {
		return err
	}

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("error deleting data stream: %s", string(resp.Body))
	}

	return nil
}

func (c *client) Rollover(name string, conditions map[string]interface{}, lazy bool) (*types.RolloverResponse, error) {
	requestBody := map[string]interface{}{}
	if conditions != nil {
		requestBody["conditions"] = conditions
	}

	path := fmt.Sprintf("/%s/_rollover", name)
	if lazy {
		path += "?lazy=true"
	}
	resp, err := c.restClient.Post(path, requestBody)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("error rolling over data stream: %s", string(resp.Body))
	}

	var rolloverResp types.RolloverResponse
	if err := json.Unmarshal(resp.Body, &rolloverResp); err != nil {
		return nil, err
	}

	return &rolloverResp, nil
}