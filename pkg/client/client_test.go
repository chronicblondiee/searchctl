package client_test

import (
	"testing"

	"github.com/chronicblondiee/searchctl/pkg/client"
)

func TestClusterHealthStruct(t *testing.T) {
	health := client.ClusterHealth{
		ClusterName:         "test-cluster",
		Status:              "green",
		NumberOfNodes:       3,
		NumberOfDataNodes:   3,
		ActivePrimaryShards: 5,
		ActiveShards:        10,
	}

	if health.ClusterName != "test-cluster" {
		t.Errorf("Expected cluster name 'test-cluster', got '%s'", health.ClusterName)
	}
	if health.Status != "green" {
		t.Errorf("Expected status 'green', got '%s'", health.Status)
	}
	if health.NumberOfNodes != 3 {
		t.Errorf("Expected 3 nodes, got %d", health.NumberOfNodes)
	}
}

func TestClusterInfoStruct(t *testing.T) {
	info := client.ClusterInfo{
		Name:        "test-node",
		ClusterName: "test-cluster",
		ClusterUUID: "test-uuid",
		Version: map[string]interface{}{
			"number": "8.0.0",
		},
		Tagline: "You Know, for Search",
	}

	if info.Name != "test-node" {
		t.Errorf("Expected name 'test-node', got '%s'", info.Name)
	}
	if info.ClusterName != "test-cluster" {
		t.Errorf("Expected cluster name 'test-cluster', got '%s'", info.ClusterName)
	}
}

func TestIndexStruct(t *testing.T) {
	index := client.Index{
		Name:             "test-index",
		Health:           "green",
		Status:           "open",
		UUID:             "test-uuid",
		Primary:          "1",
		Replica:          "1",
		DocsCount:        "100",
		StoreSize:        "1kb",
		PrimaryStoreSize: "500b",
	}

	if index.Name != "test-index" {
		t.Errorf("Expected index name 'test-index', got '%s'", index.Name)
	}
	if index.Health != "green" {
		t.Errorf("Expected health 'green', got '%s'", index.Health)
	}
}

func TestNodeStruct(t *testing.T) {
	node := client.Node{
		Name:        "node-1",
		Host:        "127.0.0.1",
		IP:          "127.0.0.1",
		HeapPercent: "50",
		RAMPercent:  "75",
		CPU:         "10",
		NodeRole:    "cdfhilmrstw",
		Master:      "*",
	}

	if node.Name != "node-1" {
		t.Errorf("Expected node name 'node-1', got '%s'", node.Name)
	}
	if node.IP != "127.0.0.1" {
		t.Errorf("Expected IP '127.0.0.1', got '%s'", node.IP)
	}
}
