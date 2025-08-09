package client

import (
	"github.com/chronicblondiee/searchctl/pkg/types"
)

type SearchClient interface {
	ClusterHealth() (*types.ClusterHealth, error)
	ClusterInfo() (*types.ClusterInfo, error)
	GetIndices(pattern string) ([]types.Index, error)
	GetIndex(name string) (*types.Index, error)
	CreateIndex(name string, body map[string]interface{}) error
	DeleteIndex(name string) error
	GetNodes() ([]types.Node, error)
	GetNode(nodeID string) (*types.Node, error)
	GetDataStreams(pattern string) ([]types.DataStream, error)
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