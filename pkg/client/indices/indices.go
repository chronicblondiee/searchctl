package indices

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/chronicblondiee/searchctl/pkg/client/rest"
	"github.com/chronicblondiee/searchctl/pkg/types"
)

type client struct {
	restClient *rest.Client
}

type templatesClient struct {
	restClient *rest.Client
}

type componentTemplatesClient struct {
	restClient *rest.Client
}

type lifecyclePoliciesClient struct {
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

func (c *client) ComponentTemplates() ComponentTemplatesInterface {
	return &componentTemplatesClient{restClient: c.restClient}
}

func (c *client) LifecyclePolicies() LifecyclePoliciesInterface {
	return &lifecyclePoliciesClient{restClient: c.restClient}
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

func (c *componentTemplatesClient) List(pattern string) ([]types.ComponentTemplate, error) {
	templatePattern := "*"
	if pattern != "" {
		templatePattern = pattern
	}

	path := fmt.Sprintf("/_component_template/%s", templatePattern)
	resp, err := c.restClient.Get(path)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("error getting component templates: %s", string(resp.Body))
	}

	var response struct {
		ComponentTemplates []struct {
			Name              string                 `json:"name"`
			ComponentTemplate types.ComponentTemplate `json:"component_template"`
		} `json:"component_templates"`
	}
	if err := json.Unmarshal(resp.Body, &response); err != nil {
		return nil, err
	}

	var templates []types.ComponentTemplate
	for _, template := range response.ComponentTemplates {
		template.ComponentTemplate.Name = template.Name
		templates = append(templates, template.ComponentTemplate)
	}

	return templates, nil
}

func (c *componentTemplatesClient) Get(name string) (*types.ComponentTemplate, error) {
	path := fmt.Sprintf("/_component_template/%s", name)
	resp, err := c.restClient.Get(path)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("error getting component template: %s", string(resp.Body))
	}

	var response struct {
		ComponentTemplates []struct {
			Name              string                 `json:"name"`
			ComponentTemplate types.ComponentTemplate `json:"component_template"`
		} `json:"component_templates"`
	}
	if err := json.Unmarshal(resp.Body, &response); err != nil {
		return nil, err
	}

	if len(response.ComponentTemplates) == 0 {
		return nil, fmt.Errorf("component template %q not found", name)
	}

	template := response.ComponentTemplates[0].ComponentTemplate
	template.Name = response.ComponentTemplates[0].Name
	return &template, nil
}

func (c *componentTemplatesClient) Create(name string, body map[string]interface{}) error {
	path := fmt.Sprintf("/_component_template/%s", name)
	resp, err := c.restClient.Put(path, body)
	if err != nil {
		return err
	}

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		return fmt.Errorf("error creating component template: %s", string(resp.Body))
	}

	return nil
}

func (c *componentTemplatesClient) Delete(name string) error {
	path := fmt.Sprintf("/_component_template/%s", name)
	resp, err := c.restClient.Delete(path)
	if err != nil {
		return err
	}

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("error deleting component template: %s", string(resp.Body))
	}

	return nil
}

func (c *lifecyclePoliciesClient) List(pattern string) ([]types.LifecyclePolicy, error) {
	policyPattern := "*"
	if pattern != "" {
		policyPattern = pattern
	}

	// Determine if this is Elasticsearch or OpenSearch based on cluster info
	// For now, try Elasticsearch first, then OpenSearch
	path := fmt.Sprintf("/_ilm/policy/%s", policyPattern)
	resp, err := c.restClient.Get(path)
	
	if err != nil || resp.StatusCode == http.StatusNotFound {
		// Try OpenSearch ISM API
		path = fmt.Sprintf("/_plugins/_ism/policies/%s", policyPattern)
		resp, err = c.restClient.Get(path)
		if err != nil {
			return nil, err
		}
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("error getting lifecycle policies: %s", string(resp.Body))
	}

	// Handle both Elasticsearch and OpenSearch response formats
	if strings.Contains(path, "_ilm") {
		// Elasticsearch ILM response format
		var response map[string]struct {
			Version      int                    `json:"version"`
			ModifiedDate string                 `json:"modified_date"`
			Policy       map[string]interface{} `json:"policy"`
		}
		if err := json.Unmarshal(resp.Body, &response); err != nil {
			return nil, err
		}

		var policies []types.LifecyclePolicy
		for name, policy := range response {
			policies = append(policies, types.LifecyclePolicy{
				Name:         name,
				Policy:       policy.Policy,
				Version:      policy.Version,
				ModifiedDate: policy.ModifiedDate,
			})
		}
		return policies, nil
	} else {
		// OpenSearch ISM response format
		var response struct {
			Policies []struct {
				ID     string                 `json:"_id"`
				Policy map[string]interface{} `json:"policy"`
			} `json:"policies"`
		}
		if err := json.Unmarshal(resp.Body, &response); err != nil {
			return nil, err
		}

		var policies []types.LifecyclePolicy
		for _, policy := range response.Policies {
			policies = append(policies, types.LifecyclePolicy{
				Name:   policy.ID,
				Policy: policy.Policy,
			})
		}
		return policies, nil
	}
}

func (c *lifecyclePoliciesClient) Get(name string) (*types.LifecyclePolicy, error) {
	// Try Elasticsearch first
	path := fmt.Sprintf("/_ilm/policy/%s", name)
	resp, err := c.restClient.Get(path)
	
	if err != nil || resp.StatusCode == http.StatusNotFound {
		// Try OpenSearch ISM API
		path = fmt.Sprintf("/_plugins/_ism/policies/%s", name)
		resp, err = c.restClient.Get(path)
		if err != nil {
			return nil, err
		}
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("error getting lifecycle policy: %s", string(resp.Body))
	}

	if strings.Contains(path, "_ilm") {
		// Elasticsearch ILM response format
		var response map[string]struct {
			Version      int                    `json:"version"`
			ModifiedDate string                 `json:"modified_date"`
			Policy       map[string]interface{} `json:"policy"`
		}
		if err := json.Unmarshal(resp.Body, &response); err != nil {
			return nil, err
		}

		if policy, exists := response[name]; exists {
			return &types.LifecyclePolicy{
				Name:         name,
				Policy:       policy.Policy,
				Version:      policy.Version,
				ModifiedDate: policy.ModifiedDate,
			}, nil
		}
		return nil, fmt.Errorf("lifecycle policy %q not found", name)
	} else {
		// OpenSearch ISM response format
		var response struct {
			Policy map[string]interface{} `json:"policy"`
		}
		if err := json.Unmarshal(resp.Body, &response); err != nil {
			return nil, err
		}

		return &types.LifecyclePolicy{
			Name:   name,
			Policy: response.Policy,
		}, nil
	}
}

func (c *lifecyclePoliciesClient) Create(name string, body map[string]interface{}) error {
	// Try Elasticsearch first
	path := fmt.Sprintf("/_ilm/policy/%s", name)
	resp, err := c.restClient.Put(path, body)
	
	// Only fallback to OpenSearch if Elasticsearch endpoint is not found or returns method not allowed
	if err != nil || resp.StatusCode == http.StatusNotFound || resp.StatusCode == http.StatusMethodNotAllowed {
		// Try OpenSearch ISM API
		path = fmt.Sprintf("/_plugins/_ism/policies/%s", name)
		resp, err = c.restClient.Put(path, body)
		if err != nil {
			return err
		}
	}

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		return fmt.Errorf("error creating lifecycle policy: %s", string(resp.Body))
	}

	return nil
}

func (c *lifecyclePoliciesClient) Delete(name string) error {
	// Try Elasticsearch first
	path := fmt.Sprintf("/_ilm/policy/%s", name)
	resp, err := c.restClient.Delete(path)
	
	if err != nil || resp.StatusCode == http.StatusNotFound {
		// Try OpenSearch ISM API
		path = fmt.Sprintf("/_plugins/_ism/policies/%s", name)
		resp, err = c.restClient.Delete(path)
		if err != nil {
			return err
		}
	}

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("error deleting lifecycle policy: %s", string(resp.Body))
	}

	return nil
}