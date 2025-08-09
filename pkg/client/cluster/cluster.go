package cluster

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

func (c *client) Health() (*types.ClusterHealth, error) {
	resp, err := c.restClient.Get("/_cluster/health")
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("error getting cluster health: %s", string(resp.Body))
	}

	var health types.ClusterHealth
	if err := json.Unmarshal(resp.Body, &health); err != nil {
		return nil, err
	}

	return &health, nil
}

func (c *client) Info() (*types.ClusterInfo, error) {
	resp, err := c.restClient.Get("/")
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("error getting cluster info: %s", string(resp.Body))
	}

	var info types.ClusterInfo
	if err := json.Unmarshal(resp.Body, &info); err != nil {
		return nil, err
	}

	return &info, nil
}