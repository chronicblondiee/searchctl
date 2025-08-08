package indices

import "github.com/chronicblondiee/searchctl/pkg/types"

type Interface interface {
	List(pattern string) ([]types.Index, error)
	Get(name string) (*types.Index, error)
	Create(name string, body map[string]interface{}) error
	Delete(name string) error
	Templates() TemplatesInterface
}

type TemplatesInterface interface {
	List(pattern string) ([]types.IndexTemplate, error)
	Get(name string) (*types.IndexTemplate, error)
	Create(name string, body map[string]interface{}) error
	Delete(name string) error
}