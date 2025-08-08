package cluster

import "github.com/chronicblondiee/searchctl/pkg/types"

type Interface interface {
	Health() (*types.ClusterHealth, error)
	Info() (*types.ClusterInfo, error)
}