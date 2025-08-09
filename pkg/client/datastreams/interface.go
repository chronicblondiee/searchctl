package datastreams

import "github.com/chronicblondiee/searchctl/pkg/types"

type Interface interface {
	List(pattern string) ([]types.DataStream, error)
	Create(name string) error
	Delete(name string) error
	Rollover(name string, conditions map[string]interface{}, lazy bool) (*types.RolloverResponse, error)
}