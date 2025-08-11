package ingest

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
	return &client{restClient: restClient}
}

func (c *client) List(pattern string) ([]types.IngestPipeline, error) {
	path := "/_ingest/pipeline"
	if pattern != "" {
		path = fmt.Sprintf("/_ingest/pipeline/%s", pattern)
	}
	resp, err := c.restClient.Get(path)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode == http.StatusNotFound || (resp.StatusCode != http.StatusOK && len(resp.Body) == 0) || (resp.StatusCode != http.StatusOK && string(resp.Body) == "{}") {
		// Treat missing endpoint or empty response as no pipelines
		return []types.IngestPipeline{}, nil
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("error getting ingest pipelines: %s", string(resp.Body))
	}
	// API returns an object keyed by pipeline id
	var body map[string]map[string]interface{}
	if err := json.Unmarshal(resp.Body, &body); err != nil {
		return nil, err
	}
	pipelines := make([]types.IngestPipeline, 0, len(body))
	for name, def := range body {
		pipelines = append(pipelines, types.IngestPipeline{
			Name: name,
			Body: def,
		})
	}
	return pipelines, nil
}

func (c *client) Get(name string) (*types.IngestPipeline, error) {
	resp, err := c.restClient.Get(fmt.Sprintf("/_ingest/pipeline/%s", name))
	if err != nil {
		return nil, err
	}
	if resp.StatusCode == http.StatusNotFound {
		return nil, fmt.Errorf("ingest pipeline %q not found", name)
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("error getting ingest pipeline: %s", string(resp.Body))
	}
	var body map[string]map[string]interface{}
	if err := json.Unmarshal(resp.Body, &body); err != nil {
		return nil, err
	}
	if def, ok := body[name]; ok {
		return &types.IngestPipeline{Name: name, Body: def}, nil
	}
	return nil, fmt.Errorf("ingest pipeline %q not found", name)
}

func (c *client) Create(name string, body map[string]interface{}) error {
	resp, err := c.restClient.Put(fmt.Sprintf("/_ingest/pipeline/%s", name), body)
	if err != nil {
		return err
	}
	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		return fmt.Errorf("error creating ingest pipeline: %s", string(resp.Body))
	}
	return nil
}

func (c *client) Delete(name string) error {
	resp, err := c.restClient.Delete(fmt.Sprintf("/_ingest/pipeline/%s", name))
	if err != nil {
		return err
	}
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("error deleting ingest pipeline: %s", string(resp.Body))
	}
	return nil
}
