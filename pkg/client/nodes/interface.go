package nodes

import "github.com/chronicblondiee/searchctl/pkg/types"

type Interface interface {
	List() ([]types.Node, error)
	Get(nodeID string) (*types.Node, error)
}