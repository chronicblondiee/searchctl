package cluster

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"

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

func (c *client) CatShards(pattern string) ([]types.CatShardRow, error) {
	shardPattern := ""
	if pattern != "" {
		shardPattern = "/" + pattern
	}
	path := fmt.Sprintf("/_cat/shards%s?format=json&h=index,shard,prirep,state,docs,store,ip,node,unassigned.reason", shardPattern)
	resp, err := c.restClient.Get(path)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("error getting shards: %s", string(resp.Body))
	}
	var rows []types.CatShardRow
	if err := json.Unmarshal(resp.Body, &rows); err != nil {
		return nil, err
	}
	return rows, nil
}

func (c *client) ExplainAllocation(req types.AllocationExplainRequest, includeYes, includeDisk bool) (*types.AllocationExplainResponse, error) {
	v := url.Values{}
	if includeYes {
		v.Set("include_yes_decisions", "true")
	}
	if includeDisk {
		v.Set("include_disk_info", "true")
	}
	path := "/_cluster/allocation/explain"
	if len(v) > 0 {
		path += "?" + v.Encode()
	}
	resp, err := c.restClient.Post(path, req)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("error explaining allocation: %s", string(resp.Body))
	}
	var out types.AllocationExplainResponse
	if err := json.Unmarshal(resp.Body, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

func (c *client) Reroute(commands []types.RerouteCommand, opts types.RerouteOptions) (*types.RerouteResponse, error) {
	v := url.Values{}
	if opts.DryRun {
		v.Set("dry_run", "true")
	}
	if opts.Explain {
		v.Set("explain", "true")
	}
	if opts.RetryFailed {
		v.Set("retry_failed", "true")
	}
	path := "/_cluster/reroute"
	if len(v) > 0 {
		path += "?" + v.Encode()
	}
	body := map[string]interface{}{"commands": commands}
	resp, err := c.restClient.Post(path, body)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("error rerouting cluster: %s", string(resp.Body))
	}
	var out types.RerouteResponse
	if err := json.Unmarshal(resp.Body, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

func (c *client) GetSettings() (*types.ClusterSettings, error) {
	resp, err := c.restClient.Get("/_cluster/settings")
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("error getting cluster settings: %s", string(resp.Body))
	}
	var out types.ClusterSettings
	if err := json.Unmarshal(resp.Body, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

func (c *client) UpdateSettings(body map[string]interface{}) error {
	resp, err := c.restClient.Put("/_cluster/settings", body)
	if err != nil {
		return err
	}
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("error updating cluster settings: %s", string(resp.Body))
	}
	return nil
}
