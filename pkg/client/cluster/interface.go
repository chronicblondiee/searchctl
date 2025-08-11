package cluster

import "github.com/chronicblondiee/searchctl/pkg/types"

type Interface interface {
	Health() (*types.ClusterHealth, error)
	Info() (*types.ClusterInfo, error)
	CatShards(pattern string) ([]types.CatShardRow, error)
	ExplainAllocation(req types.AllocationExplainRequest, includeYes, includeDisk bool) (*types.AllocationExplainResponse, error)
	Reroute(commands []types.RerouteCommand, opts types.RerouteOptions) (*types.RerouteResponse, error)
	GetSettings() (*types.ClusterSettings, error)
	UpdateSettings(body map[string]interface{}) error

	// New operations
	Stats() (*types.ClusterStats, error)
	State(metrics []string, indices string, masterTimeout string) (*types.ClusterState, error)
	PendingTasks() (*types.PendingTasks, error)
}
