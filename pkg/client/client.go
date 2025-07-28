package client

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/elastic/go-elasticsearch/v8"
	"github.com/chronicblondiee/searchctl/pkg/config"
)

type SearchClient interface {
	ClusterHealth() (*ClusterHealth, error)
	ClusterInfo() (*ClusterInfo, error)
	GetIndices(pattern string) ([]Index, error)
	GetIndex(name string) (*Index, error)
	CreateIndex(name string, body map[string]interface{}) error
	DeleteIndex(name string) error
	GetNodes() ([]Node, error)
	GetNode(nodeID string) (*Node, error)
}

type Client struct {
	es *elasticsearch.Client
}

type ClusterHealth struct {
	ClusterName         string `json:"cluster_name"`
	Status              string `json:"status"`
	TimedOut            bool   `json:"timed_out"`
	NumberOfNodes       int    `json:"number_of_nodes"`
	NumberOfDataNodes   int    `json:"number_of_data_nodes"`
	ActivePrimaryShards int    `json:"active_primary_shards"`
	ActiveShards        int    `json:"active_shards"`
	RelocatingShards    int    `json:"relocating_shards"`
	InitializingShards  int    `json:"initializing_shards"`
	UnassignedShards    int    `json:"unassigned_shards"`
}

type ClusterInfo struct {
	Name        string                 `json:"name"`
	ClusterName string                 `json:"cluster_name"`
	ClusterUUID string                 `json:"cluster_uuid"`
	Version     map[string]interface{} `json:"version"`
	Tagline     string                 `json:"tagline"`
}

type Index struct {
	Name     string                 `json:"index"`
	Health   string                 `json:"health"`
	Status   string                 `json:"status"`
	UUID     string                 `json:"uuid"`
	Primary  string                 `json:"pri"`
	Replica  string                 `json:"rep"`
	DocsCount string                `json:"docs.count"`
	DocsDeleted string              `json:"docs.deleted"`
	StoreSize string                `json:"store.size"`
	PrimaryStoreSize string         `json:"pri.store.size"`
}

type Node struct {
	Name             string  `json:"name"`
	Host             string  `json:"host"`
	IP               string  `json:"ip"`
	HeapPercent      string  `json:"heap.percent"`
	RAMPercent       string  `json:"ram.percent"`
	CPU              string  `json:"cpu"`
	Load1m           string  `json:"load_1m"`
	Load5m           string  `json:"load_5m"`
	Load15m          string  `json:"load_15m"`
	NodeRole         string  `json:"node.role"`
	Master           string  `json:"master"`
}

func NewClient() (SearchClient, error) {
	ctx, err := config.GetCurrentContext()
	if err != nil {
		return nil, fmt.Errorf("failed to get current context: %w", err)
	}

	cluster, err := config.GetCluster(ctx.Context.Cluster)
	if err != nil {
		return nil, fmt.Errorf("failed to get cluster config: %w", err)
	}

	user, err := config.GetUser(ctx.Context.User)
	if err != nil {
		return nil, fmt.Errorf("failed to get user config: %w", err)
	}

	cfg := elasticsearch.Config{
		Addresses: []string{cluster.Cluster.Server},
	}

	if cluster.Cluster.InsecureSkipTLSVerify {
		cfg.Transport = &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true,
			},
		}
	}

	if user.User.Username != "" && user.User.Password != "" {
		cfg.Username = user.User.Username
		cfg.Password = user.User.Password
	}

	if user.User.APIKey != "" {
		cfg.APIKey = user.User.APIKey
	}

	es, err := elasticsearch.NewClient(cfg)
	if err != nil {
		return nil, fmt.Errorf("failed to create elasticsearch client: %w", err)
	}

	return &Client{es: es}, nil
}

func (c *Client) ClusterHealth() (*ClusterHealth, error) {
	res, err := c.es.Cluster.Health()
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.IsError() {
		return nil, fmt.Errorf("error getting cluster health: %s", res.String())
	}

	var health ClusterHealth
	if err := json.NewDecoder(res.Body).Decode(&health); err != nil {
		return nil, err
	}

	return &health, nil
}

func (c *Client) ClusterInfo() (*ClusterInfo, error) {
	res, err := c.es.Info()
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.IsError() {
		return nil, fmt.Errorf("error getting cluster info: %s", res.String())
	}

	var info ClusterInfo
	if err := json.NewDecoder(res.Body).Decode(&info); err != nil {
		return nil, err
	}

	return &info, nil
}

func (c *Client) GetIndices(pattern string) ([]Index, error) {
	indexPattern := "_all"
	if pattern != "" {
		indexPattern = pattern
	}

	res, err := c.es.Cat.Indices(
		c.es.Cat.Indices.WithIndex(indexPattern),
		c.es.Cat.Indices.WithFormat("json"),
		c.es.Cat.Indices.WithH("index,health,status,uuid,pri,rep,docs.count,docs.deleted,store.size,pri.store.size"),
	)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.IsError() {
		return nil, fmt.Errorf("error getting indices: %s", res.String())
	}

	var indices []Index
	if err := json.NewDecoder(res.Body).Decode(&indices); err != nil {
		return nil, err
	}

	return indices, nil
}

func (c *Client) GetIndex(name string) (*Index, error) {
	indices, err := c.GetIndices(name)
	if err != nil {
		return nil, err
	}

	for _, index := range indices {
		if index.Name == name {
			return &index, nil
		}
	}

	return nil, fmt.Errorf("index %q not found", name)
}

func (c *Client) CreateIndex(name string, body map[string]interface{}) error {
	var bodyReader *strings.Reader
	if body != nil {
		bodyBytes, err := json.Marshal(body)
		if err != nil {
			return err
		}
		bodyReader = strings.NewReader(string(bodyBytes))
	}

	res, err := c.es.Indices.Create(name, c.es.Indices.Create.WithBody(bodyReader))
	if err != nil {
		return err
	}
	defer res.Body.Close()

	if res.IsError() {
		return fmt.Errorf("error creating index: %s", res.String())
	}

	return nil
}

func (c *Client) DeleteIndex(name string) error {
	res, err := c.es.Indices.Delete([]string{name})
	if err != nil {
		return err
	}
	defer res.Body.Close()

	if res.IsError() {
		return fmt.Errorf("error deleting index: %s", res.String())
	}

	return nil
}

func (c *Client) GetNodes() ([]Node, error) {
	res, err := c.es.Cat.Nodes(
		c.es.Cat.Nodes.WithFormat("json"),
		c.es.Cat.Nodes.WithH("name,host,ip,heap.percent,ram.percent,cpu,load_1m,load_5m,load_15m,node.role,master"),
	)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.IsError() {
		return nil, fmt.Errorf("error getting nodes: %s", res.String())
	}

	var nodes []Node
	if err := json.NewDecoder(res.Body).Decode(&nodes); err != nil {
		return nil, err
	}

	return nodes, nil
}

func (c *Client) GetNode(nodeID string) (*Node, error) {
	nodes, err := c.GetNodes()
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
