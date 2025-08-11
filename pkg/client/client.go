package client

import (
	"github.com/chronicblondiee/searchctl/pkg/types"
)

type SearchClient interface {
	ClusterHealth() (*types.ClusterHealth, error)
	ClusterInfo() (*types.ClusterInfo, error)
	ClusterStats() (*types.ClusterStats, error)
	ClusterState(metrics []string, indices, masterTimeout string) (*types.ClusterState, error)
	ClusterPendingTasks() (*types.ClusterPendingTasks, error)
	GetIndices(pattern string) ([]types.Index, error)
	GetIndex(name string) (*types.Index, error)
	CreateIndex(name string, body map[string]interface{}) error
	DeleteIndex(name string) error
	GetNodes() ([]types.Node, error)
	GetNode(nodeID string) (*types.Node, error)
	GetDataStreams(pattern string) ([]types.DataStream, error)
	GetDataStream(name string) (*types.DataStream, error)
	CreateDataStream(name string) error
	DeleteDataStream(name string) error
	RolloverDataStream(name string, conditions map[string]interface{}, lazy bool) (*types.RolloverResponse, error)
	GetIndexTemplates(pattern string) ([]types.IndexTemplate, error)
	GetIndexTemplate(name string) (*types.IndexTemplate, error)
	CreateIndexTemplate(name string, body map[string]interface{}) error
	DeleteIndexTemplate(name string) error
	GetComponentTemplates(pattern string) ([]types.ComponentTemplate, error)
	GetComponentTemplate(name string) (*types.ComponentTemplate, error)
	CreateComponentTemplate(name string, body map[string]interface{}) error
	DeleteComponentTemplate(name string) error
	GetLifecyclePolicies(pattern string) ([]types.LifecyclePolicy, error)
	GetLifecyclePolicy(name string) (*types.LifecyclePolicy, error)
	CreateLifecyclePolicy(name string, body map[string]interface{}) error
	DeleteLifecyclePolicy(name string) error
	GetShards(pattern string) ([]types.CatShardRow, error)
	ExplainAllocation(req types.AllocationExplainRequest, includeYes, includeDisk bool) (*types.AllocationExplainResponse, error)
	Reroute(commands []types.RerouteCommand, opts types.RerouteOptions) (*types.RerouteResponse, error)
	GetClusterSettings() (*types.ClusterSettings, error)
	UpdateClusterSettings(body map[string]interface{}) error
	GetIngestPipelines(pattern string) ([]types.IngestPipeline, error)
	GetIngestPipeline(name string) (*types.IngestPipeline, error)
	CreateIngestPipeline(name string, body map[string]interface{}) error
	DeleteIngestPipeline(name string) error
}

type Client struct {
	clientset Interface
}

func NewClient() (SearchClient, error) {
	clientset, err := NewClientset()
	if err != nil {
		return nil, err
	}

	return &Client{
		clientset: clientset,
	}, nil
}

func (c *Client) ClusterHealth() (*types.ClusterHealth, error) {
	return c.clientset.Cluster().Health()
}

func (c *Client) ClusterInfo() (*types.ClusterInfo, error) {
	return c.clientset.Cluster().Info()
}

func (c *Client) ClusterStats() (*types.ClusterStats, error) {
	return c.clientset.Cluster().Stats()
}

func (c *Client) ClusterState(metrics []string, indices, masterTimeout string) (*types.ClusterState, error) {
	return c.clientset.Cluster().State(metrics, indices, masterTimeout)
}

func (c *Client) ClusterPendingTasks() (*types.ClusterPendingTasks, error) {
	return c.clientset.Cluster().PendingTasks()
}

func (c *Client) GetIndices(pattern string) ([]types.Index, error) {
	return c.clientset.Indices().List(pattern)
}

func (c *Client) GetIndex(name string) (*types.Index, error) {
	return c.clientset.Indices().Get(name)
}

func (c *Client) CreateIndex(name string, body map[string]interface{}) error {
	return c.clientset.Indices().Create(name, body)
}

func (c *Client) DeleteIndex(name string) error {
	return c.clientset.Indices().Delete(name)
}

func (c *Client) GetNodes() ([]types.Node, error) {
	return c.clientset.Nodes().List()
}

func (c *Client) GetNode(nodeID string) (*types.Node, error) {
	return c.clientset.Nodes().Get(nodeID)
}

func (c *Client) GetDataStreams(pattern string) ([]types.DataStream, error) {
	return c.clientset.DataStreams().List(pattern)
}

func (c *Client) GetDataStream(name string) (*types.DataStream, error) {
	return c.clientset.DataStreams().Get(name)
}

func (c *Client) CreateDataStream(name string) error {
	return c.clientset.DataStreams().Create(name)
}

func (c *Client) DeleteDataStream(name string) error {
	return c.clientset.DataStreams().Delete(name)
}

func (c *Client) RolloverDataStream(name string, conditions map[string]interface{}, lazy bool) (*types.RolloverResponse, error) {
	return c.clientset.DataStreams().Rollover(name, conditions, lazy)
}

func (c *Client) GetIndexTemplates(pattern string) ([]types.IndexTemplate, error) {
	return c.clientset.Indices().Templates().List(pattern)
}

func (c *Client) GetIndexTemplate(name string) (*types.IndexTemplate, error) {
	return c.clientset.Indices().Templates().Get(name)
}

func (c *Client) CreateIndexTemplate(name string, body map[string]interface{}) error {
	return c.clientset.Indices().Templates().Create(name, body)
}

func (c *Client) DeleteIndexTemplate(name string) error {
	return c.clientset.Indices().Templates().Delete(name)
}

func (c *Client) GetComponentTemplates(pattern string) ([]types.ComponentTemplate, error) {
	return c.clientset.Indices().ComponentTemplates().List(pattern)
}

func (c *Client) GetComponentTemplate(name string) (*types.ComponentTemplate, error) {
	return c.clientset.Indices().ComponentTemplates().Get(name)
}

func (c *Client) CreateComponentTemplate(name string, body map[string]interface{}) error {
	return c.clientset.Indices().ComponentTemplates().Create(name, body)
}

func (c *Client) DeleteComponentTemplate(name string) error {
	return c.clientset.Indices().ComponentTemplates().Delete(name)
}

func (c *Client) GetLifecyclePolicies(pattern string) ([]types.LifecyclePolicy, error) {
	return c.clientset.Indices().LifecyclePolicies().List(pattern)
}

func (c *Client) GetLifecyclePolicy(name string) (*types.LifecyclePolicy, error) {
	return c.clientset.Indices().LifecyclePolicies().Get(name)
}

func (c *Client) CreateLifecyclePolicy(name string, body map[string]interface{}) error {
	return c.clientset.Indices().LifecyclePolicies().Create(name, body)
}

func (c *Client) DeleteLifecyclePolicy(name string) error {
	return c.clientset.Indices().LifecyclePolicies().Delete(name)
}

func (c *Client) GetShards(pattern string) ([]types.CatShardRow, error) {
	return c.clientset.Cluster().CatShards(pattern)
}

func (c *Client) ExplainAllocation(req types.AllocationExplainRequest, includeYes, includeDisk bool) (*types.AllocationExplainResponse, error) {
	return c.clientset.Cluster().ExplainAllocation(req, includeYes, includeDisk)
}

func (c *Client) Reroute(commands []types.RerouteCommand, opts types.RerouteOptions) (*types.RerouteResponse, error) {
	return c.clientset.Cluster().Reroute(commands, opts)
}

func (c *Client) GetClusterSettings() (*types.ClusterSettings, error) {
	return c.clientset.Cluster().GetSettings()
}

func (c *Client) UpdateClusterSettings(body map[string]interface{}) error {
	return c.clientset.Cluster().UpdateSettings(body)
}

func (c *Client) GetIngestPipelines(pattern string) ([]types.IngestPipeline, error) {
	return c.clientset.Ingest().List(pattern)
}

func (c *Client) GetIngestPipeline(name string) (*types.IngestPipeline, error) {
	return c.clientset.Ingest().Get(name)
}

func (c *Client) CreateIngestPipeline(name string, body map[string]interface{}) error {
	return c.clientset.Ingest().Create(name, body)
}

func (c *Client) DeleteIngestPipeline(name string) error {
	return c.clientset.Ingest().Delete(name)
}
