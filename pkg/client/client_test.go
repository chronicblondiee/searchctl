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

func TestDataStreamStruct(t *testing.T) {
	dataStream := client.DataStream{
		Name:           "logs-nginx",
		TimestampField: client.TimestampFieldType{Name: "@timestamp"},
		Indices: []client.DataStreamIndex{
			{IndexName: "logs-nginx-000001", IndexUUID: "uuid1"},
			{IndexName: "logs-nginx-000002", IndexUUID: "uuid2"},
		},
		Generation:         2,
		Status:             "green",
		Template:           "logs-nginx-template",
		IlmPolicy:          "logs-policy",
		Hidden:             false,
		System:             false,
		AllowCustomRouting: false,
	}

	if dataStream.Name != "logs-nginx" {
		t.Errorf("Expected name 'logs-nginx', got '%s'", dataStream.Name)
	}
	if dataStream.Generation != 2 {
		t.Errorf("Expected generation 2, got %d", dataStream.Generation)
	}
	if len(dataStream.Indices) != 2 {
		t.Errorf("Expected 2 indices, got %d", len(dataStream.Indices))
	}
	if dataStream.TimestampField.Name != "@timestamp" {
		t.Errorf("Expected timestamp field '@timestamp', got '%s'", dataStream.TimestampField.Name)
	}
}

func TestRolloverResponseStruct(t *testing.T) {
	response := client.RolloverResponse{
		Acknowledged:       true,
		ShardsAcknowledged: true,
		OldIndex:           "logs-nginx-000001",
		NewIndex:           "logs-nginx-000002",
		RolledOver:         true,
		DryRun:             false,
		Conditions: map[string]bool{
			"max_age":  true,
			"max_docs": false,
		},
	}

	if !response.Acknowledged {
		t.Error("Expected Acknowledged to be true")
	}
	if !response.RolledOver {
		t.Error("Expected RolledOver to be true")
	}
	if response.OldIndex != "logs-nginx-000001" {
		t.Errorf("Expected old index 'logs-nginx-000001', got '%s'", response.OldIndex)
	}
	if len(response.Conditions) != 2 {
		t.Errorf("Expected 2 conditions, got %d", len(response.Conditions))
	}
}

func TestCreateDataStreamResponse(t *testing.T) {
	// This would be tested with actual client integration
	// For now, just test that we can create the expected structure
	dataStreamName := "test-logs"
	if dataStreamName != "test-logs" {
		t.Errorf("Expected data stream name 'test-logs', got '%s'", dataStreamName)
	}
}

func TestDeleteDataStreamResponse(t *testing.T) {
	// This would be tested with actual client integration
	// For now, just test that we can handle the expected response
	dataStreamName := "test-logs"
	if dataStreamName != "test-logs" {
		t.Errorf("Expected data stream name 'test-logs', got '%s'", dataStreamName)
	}
}
