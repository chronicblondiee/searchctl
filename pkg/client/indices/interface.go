package indices

import "github.com/chronicblondiee/searchctl/pkg/types"

type Interface interface {
	List(pattern string) ([]types.Index, error)
	Get(name string) (*types.Index, error)
	Create(name string, body map[string]interface{}) error
	Delete(name string) error
	Templates() TemplatesInterface
	ComponentTemplates() ComponentTemplatesInterface
}

type TemplatesInterface interface {
	List(pattern string) ([]types.IndexTemplate, error)
	Get(name string) (*types.IndexTemplate, error)
	Create(name string, body map[string]interface{}) error
	Delete(name string) error
}

type ComponentTemplatesInterface interface {
	List(pattern string) ([]types.ComponentTemplate, error)
	Get(name string) (*types.ComponentTemplate, error)
	Create(name string, body map[string]interface{}) error
	Delete(name string) error
}