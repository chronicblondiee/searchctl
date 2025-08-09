package indices

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

type templatesClient struct {
	restClient *rest.Client
}

func New(restClient *rest.Client) Interface {
	return &client{
		restClient: restClient,
	}
}

func (c *client) List(pattern string) ([]types.Index, error) {
	indexPattern := "_all"
	if pattern != "" {
		indexPattern = pattern
	}

	path := fmt.Sprintf("/_cat/indices/%s?format=json&h=index,health,status,uuid,pri,rep,docs.count,docs.deleted,store.size,pri.store.size", indexPattern)
	resp, err := c.restClient.Get(path)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("error getting indices: %s", string(resp.Body))
	}

	var indices []types.Index
	if err := json.Unmarshal(resp.Body, &indices); err != nil {
		return nil, err
	}

	return indices, nil
}

func (c *client) Get(name string) (*types.Index, error) {
	indices, err := c.List(name)
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

func (c *client) Create(name string, body map[string]interface{}) error {
	path := fmt.Sprintf("/%s", name)
	resp, err := c.restClient.Put(path, body)
	if err != nil {
		return err
	}

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		return fmt.Errorf("error creating index: %s", string(resp.Body))
	}

	return nil
}

func (c *client) Delete(name string) error {
	path := fmt.Sprintf("/%s", name)
	resp, err := c.restClient.Delete(path)
	if err != nil {
		return err
	}

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("error deleting index: %s", string(resp.Body))
	}

	return nil
}

func (c *client) Templates() TemplatesInterface {
	return &templatesClient{restClient: c.restClient}
}

func (c *templatesClient) List(pattern string) ([]types.IndexTemplate, error) {
	templatePattern := "*"
	if pattern != "" {
		templatePattern = pattern
	}

	path := fmt.Sprintf("/_index_template/%s", templatePattern)
	resp, err := c.restClient.Get(path)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("error getting index templates: %s", string(resp.Body))
	}

	var response struct {
		IndexTemplates []struct {
			Name          string             `json:"name"`
			IndexTemplate types.IndexTemplate `json:"index_template"`
		} `json:"index_templates"`
	}
	if err := json.Unmarshal(resp.Body, &response); err != nil {
		return nil, err
	}

	var templates []types.IndexTemplate
	for _, template := range response.IndexTemplates {
		template.IndexTemplate.Name = template.Name
		templates = append(templates, template.IndexTemplate)
	}

	return templates, nil
}

func (c *templatesClient) Get(name string) (*types.IndexTemplate, error) {
	path := fmt.Sprintf("/_index_template/%s", name)
	resp, err := c.restClient.Get(path)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("error getting index template: %s", string(resp.Body))
	}

	var response struct {
		IndexTemplates []struct {
			Name          string             `json:"name"`
			IndexTemplate types.IndexTemplate `json:"index_template"`
		} `json:"index_templates"`
	}
	if err := json.Unmarshal(resp.Body, &response); err != nil {
		return nil, err
	}

	if len(response.IndexTemplates) == 0 {
		return nil, fmt.Errorf("index template %q not found", name)
	}

	template := response.IndexTemplates[0].IndexTemplate
	template.Name = response.IndexTemplates[0].Name
	return &template, nil
}

func (c *templatesClient) Create(name string, body map[string]interface{}) error {
	path := fmt.Sprintf("/_index_template/%s", name)
	resp, err := c.restClient.Put(path, body)
	if err != nil {
		return err
	}

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		return fmt.Errorf("error creating index template: %s", string(resp.Body))
	}

	return nil
}

func (c *templatesClient) Delete(name string) error {
	path := fmt.Sprintf("/_index_template/%s", name)
	resp, err := c.restClient.Delete(path)
	if err != nil {
		return err
	}

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("error deleting index template: %s", string(resp.Body))
	}

	return nil
}