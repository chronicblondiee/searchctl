package ingest

import "github.com/chronicblondiee/searchctl/pkg/types"

type Interface interface {
	List(pattern string) ([]types.IngestPipeline, error)
	Get(name string) (*types.IngestPipeline, error)
	Create(name string, body map[string]interface{}) error
	Delete(name string) error
}
