package nodes

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

func (c *client) List() ([]types.Node, error) {
	resp, err := c.restClient.Get("/_cat/nodes?format=json&h=name,host,ip,heap.percent,ram.percent,cpu,load_1m,load_5m,load_15m,node.role,master")
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("error getting nodes: %s", string(resp.Body))
	}

	var nodes []types.Node
	if err := json.Unmarshal(resp.Body, &nodes); err != nil {
		return nil, err
	}

	return nodes, nil
}

func (c *client) Get(nodeID string) (*types.Node, error) {
	nodes, err := c.List()
	if err != nil {
		return nil, err
	}

	for _, node := range nodes {
		if node.Name == nodeID || node.IP == nodeID {
			return &node, nil
		}
	}

	return nil, fmt.Errorf("node %q not found", nodeID)
}