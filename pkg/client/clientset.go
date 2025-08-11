package client

import (
	"github.com/chronicblondiee/searchctl/pkg/client/cluster"
	"github.com/chronicblondiee/searchctl/pkg/client/datastreams"
	"github.com/chronicblondiee/searchctl/pkg/client/indices"
	"github.com/chronicblondiee/searchctl/pkg/client/ingest"
	"github.com/chronicblondiee/searchctl/pkg/client/nodes"
	"github.com/chronicblondiee/searchctl/pkg/client/rest"
)

type Interface interface {
	Cluster() cluster.Interface
	Indices() indices.Interface
	DataStreams() datastreams.Interface
	Nodes() nodes.Interface
	Ingest() ingest.Interface
}

type Clientset struct {
	clusterClient     cluster.Interface
	indicesClient     indices.Interface
	dataStreamsClient datastreams.Interface
	nodesClient       nodes.Interface
	ingestClient      ingest.Interface
}

func NewClientset() (Interface, error) {
	factory, err := NewFactory()
	if err != nil {
		return nil, err
	}

	restClient := rest.NewClient(&rest.Config{
		HTTPClient: factory.HTTPClient(),
		BaseURL:    factory.BaseURL(),
		Username:   factory.Username(),
		Password:   factory.Password(),
		APIKey:     factory.APIKey(),
	})

	return &Clientset{
		clusterClient:     cluster.New(restClient),
		indicesClient:     indices.New(restClient),
		dataStreamsClient: datastreams.New(restClient),
		nodesClient:       nodes.New(restClient),
		ingestClient:      ingest.New(restClient),
	}, nil
}

func (c *Clientset) Cluster() cluster.Interface {
	return c.clusterClient
}

func (c *Clientset) Indices() indices.Interface {
	return c.indicesClient
}

func (c *Clientset) DataStreams() datastreams.Interface {
	return c.dataStreamsClient
}

func (c *Clientset) Nodes() nodes.Interface {
	return c.nodesClient
}

func (c *Clientset) Ingest() ingest.Interface {
	return c.ingestClient
}
