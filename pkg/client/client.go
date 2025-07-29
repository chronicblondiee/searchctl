package client

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

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
	GetDataStreams(pattern string) ([]DataStream, error)
	CreateDataStream(name string) error
	DeleteDataStream(name string) error
	RolloverDataStream(name string, conditions map[string]interface{}, lazy bool) (*RolloverResponse, error)
}

type Client struct {
	httpClient *http.Client
	baseURL    string
	username   string
	password   string
	apiKey     string
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
	Name             string `json:"index"`
	Health           string `json:"health"`
	Status           string `json:"status"`
	UUID             string `json:"uuid"`
	Primary          string `json:"pri"`
	Replica          string `json:"rep"`
	DocsCount        string `json:"docs.count"`
	DocsDeleted      string `json:"docs.deleted"`
	StoreSize        string `json:"store.size"`
	PrimaryStoreSize string `json:"pri.store.size"`
}

type Node struct {
	Name        string `json:"name"`
	Host        string `json:"host"`
	IP          string `json:"ip"`
	HeapPercent string `json:"heap.percent"`
	RAMPercent  string `json:"ram.percent"`
	CPU         string `json:"cpu"`
	Load1m      string `json:"load_1m"`
	Load5m      string `json:"load_5m"`
	Load15m     string `json:"load_15m"`
	NodeRole    string `json:"node.role"`
	Master      string `json:"master"`
}

type DataStream struct {
	Name               string             `json:"name"`
	TimestampField     TimestampFieldType `json:"timestamp_field"`
	Indices            []DataStreamIndex  `json:"indices"`
	Generation         int                `json:"generation"`
	Status             string             `json:"status"`
	Template           string             `json:"template,omitempty"`
	IlmPolicy          string             `json:"ilm_policy,omitempty"`
	Hidden             bool               `json:"hidden,omitempty"`
	System             bool               `json:"system,omitempty"`
	AllowCustomRouting bool               `json:"allow_custom_routing,omitempty"`
}

type TimestampFieldType struct {
	Name string `json:"name"`
}

type DataStreamIndex struct {
	IndexName string `json:"index_name"`
	IndexUUID string `json:"index_uuid"`
	PreferILM bool   `json:"prefer_ilm,omitempty"`
	ManagedBy string `json:"managed_by,omitempty"`
}

type RolloverResponse struct {
	Acknowledged       bool            `json:"acknowledged"`
	ShardsAcknowledged bool            `json:"shards_acknowledged"`
	OldIndex           string          `json:"old_index"`
	NewIndex           string          `json:"new_index"`
	RolledOver         bool            `json:"rolled_over"`
	DryRun             bool            `json:"dry_run"`
	Conditions         map[string]bool `json:"conditions"`
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

	httpClient := &http.Client{}

	if cluster.Cluster.InsecureSkipTLSVerify {
		httpClient.Transport = &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true,
			},
		}
	}

	client := &Client{
		httpClient: httpClient,
		baseURL:    strings.TrimSuffix(cluster.Cluster.Server, "/"),
		username:   user.User.Username,
		password:   user.User.Password,
		apiKey:     user.User.APIKey,
	}

	return client, nil
}

func (c *Client) makeRequest(method, path string, body []byte) (*http.Response, error) {
	url := c.baseURL + path

	var reqBody io.Reader
	if body != nil {
		reqBody = bytes.NewBuffer(body)
	}

	req, err := http.NewRequest(method, url, reqBody)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")

	if c.apiKey != "" {
		req.Header.Set("Authorization", "ApiKey "+c.apiKey)
	} else if c.username != "" && c.password != "" {
		req.SetBasicAuth(c.username, c.password)
	}

	return c.httpClient.Do(req)
}

func (c *Client) ClusterHealth() (*ClusterHealth, error) {
	resp, err := c.makeRequest("GET", "/_cluster/health", nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("error getting cluster health: %s", string(body))
	}

	var health ClusterHealth
	if err := json.NewDecoder(resp.Body).Decode(&health); err != nil {
		return nil, err
	}

	return &health, nil
}

func (c *Client) ClusterInfo() (*ClusterInfo, error) {
	resp, err := c.makeRequest("GET", "/", nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("error getting cluster info: %s", string(body))
	}

	var info ClusterInfo
	if err := json.NewDecoder(resp.Body).Decode(&info); err != nil {
		return nil, err
	}

	return &info, nil
}

func (c *Client) GetIndices(pattern string) ([]Index, error) {
	indexPattern := "_all"
	if pattern != "" {
		indexPattern = pattern
	}

	path := fmt.Sprintf("/_cat/indices/%s?format=json&h=index,health,status,uuid,pri,rep,docs.count,docs.deleted,store.size,pri.store.size", indexPattern)
	resp, err := c.makeRequest("GET", path, nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("error getting indices: %s", string(body))
	}

	var indices []Index
	if err := json.NewDecoder(resp.Body).Decode(&indices); err != nil {
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
	var bodyBytes []byte
	if body != nil {
		var err error
		bodyBytes, err = json.Marshal(body)
		if err != nil {
			return err
		}
	}

	path := fmt.Sprintf("/%s", name)
	resp, err := c.makeRequest("PUT", path, bodyBytes)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("error creating index: %s", string(body))
	}

	return nil
}

func (c *Client) DeleteIndex(name string) error {
	path := fmt.Sprintf("/%s", name)
	resp, err := c.makeRequest("DELETE", path, nil)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("error deleting index: %s", string(body))
	}

	return nil
}

func (c *Client) GetNodes() ([]Node, error) {
	resp, err := c.makeRequest("GET", "/_cat/nodes?format=json&h=name,host,ip,heap.percent,ram.percent,cpu,load_1m,load_5m,load_15m,node.role,master", nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("error getting nodes: %s", string(body))
	}

	var nodes []Node
	if err := json.NewDecoder(resp.Body).Decode(&nodes); err != nil {
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

func (c *Client) GetDataStreams(pattern string) ([]DataStream, error) {
	dataStreamPattern := "*"
	if pattern != "" {
		dataStreamPattern = pattern
	}

	path := fmt.Sprintf("/_data_stream/%s", dataStreamPattern)
	resp, err := c.makeRequest("GET", path, nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("error getting data streams: %s", string(body))
	}

	var response struct {
		DataStreams []DataStream `json:"data_streams"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, err
	}

	return response.DataStreams, nil
}

func (c *Client) CreateDataStream(name string) error {
	path := fmt.Sprintf("/_data_stream/%s", name)
	resp, err := c.makeRequest("PUT", path, nil)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("error creating data stream: %s", string(body))
	}

	return nil
}

func (c *Client) DeleteDataStream(name string) error {
	path := fmt.Sprintf("/_data_stream/%s", name)
	resp, err := c.makeRequest("DELETE", path, nil)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("error deleting data stream: %s", string(body))
	}

	return nil
}

func (c *Client) RolloverDataStream(name string, conditions map[string]interface{}, lazy bool) (*RolloverResponse, error) {
	var bodyBytes []byte
	if conditions != nil || lazy {
		requestBody := map[string]interface{}{}
		if conditions != nil {
			requestBody["conditions"] = conditions
		}
		var err error
		bodyBytes, err = json.Marshal(requestBody)
		if err != nil {
			return nil, err
		}
	}

	path := fmt.Sprintf("/%s/_rollover", name)
	if lazy {
		path += "?lazy=true"
	}
	resp, err := c.makeRequest("POST", path, bodyBytes)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("error rolling over data stream: %s", string(body))
	}

	var rolloverResp RolloverResponse
	if err := json.NewDecoder(resp.Body).Decode(&rolloverResp); err != nil {
		return nil, err
	}

	return &rolloverResp, nil
}
