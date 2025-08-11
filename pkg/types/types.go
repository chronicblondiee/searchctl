package types

type ClusterHealth struct {
	ClusterName         string `json:"cluster_name"`
	Status              string `json:"status"`
	TimedOut            bool   `json:"timed_out"`
	NumberOfNodes       int    `json:"number_of_nodes"`
	NumberOfDataNodes   int    `json:"number_of_data_nodes"`
	ActivePrimaryShards int    `json:"active_primary_shards"`
	ActiveShards        int    `json:"active_shards"`
	RelocatingShards    int    `json:"relocating_shards"`
	InitializingShards  int    `json:"initializing_shards"`
	UnassignedShards    int    `json:"unassigned_shards"`
}

type ClusterInfo struct {
	Name        string                 `json:"name"`
	ClusterName string                 `json:"cluster_name"`
	ClusterUUID string                 `json:"cluster_uuid"`
	Version     map[string]interface{} `json:"version"`
	Tagline     string                 `json:"tagline"`
}

// CatShardRow represents a row from _cat/shards
type CatShardRow struct {
	Index            string `json:"index"`
	Shard            string `json:"shard"`
	PrimaryOrReplica string `json:"prirep"`
	State            string `json:"state"`
	Docs             string `json:"docs"`
	Store            string `json:"store"`
	IP               string `json:"ip"`
	Node             string `json:"node"`
	UnassignedReason string `json:"unassigned.reason,omitempty"`
}

// AllocationExplainRequest describes a shard to explain
type AllocationExplainRequest struct {
	Index   string `json:"index,omitempty"`
	Shard   int    `json:"shard"`
	Primary bool   `json:"primary,omitempty"`
}

// AllocationExplainResponse is a simplified view of explain output
type AllocationExplainResponse struct {
	Index               string                   `json:"index"`
	Shard               int                      `json:"shard"`
	Primary             bool                     `json:"primary"`
	CurrentNode         map[string]interface{}   `json:"current_node,omitempty"`
	NodeExplanations    []map[string]interface{} `json:"node_explanations,omitempty"`
	CanAllocate         string                   `json:"can_allocate"`
	AllocateExplanation string                   `json:"allocate_explanation,omitempty"`
	UnassignedInfo      map[string]interface{}   `json:"unassigned_info,omitempty"`
}

// RerouteCommand supports multiple command forms
type RerouteCommand map[string]map[string]interface{}

type RerouteOptions struct {
	DryRun      bool
	Explain     bool
	RetryFailed bool
}

type RerouteResponse struct {
	State        map[string]interface{}   `json:"state,omitempty"`
	Explanations []map[string]interface{} `json:"explanations,omitempty"`
}

type ClusterSettings struct {
	Persistent map[string]interface{} `json:"persistent"`
	Transient  map[string]interface{} `json:"transient"`
}

// ClusterStats is a lightweight representation of /_cluster/stats
// Keep nested structures loosely typed to accommodate ES/OS differences without heavy typing.
type ClusterStats struct {
	ClusterName string                 `json:"cluster_name"`
	Indices     map[string]interface{} `json:"indices"`
	Nodes       map[string]interface{} `json:"nodes"`
}

// ClusterState represents /_cluster/state with optional metric filtering
type ClusterState struct {
	ClusterName  string                 `json:"cluster_name"`
	StateUUID    string                 `json:"state_uuid"`
	Metadata     map[string]interface{} `json:"metadata,omitempty"`
	RoutingTable map[string]interface{} `json:"routing_table,omitempty"`
	Blocks       map[string]interface{} `json:"blocks,omitempty"`
	Nodes        map[string]interface{} `json:"nodes,omitempty"`
}

// ClusterPendingTasks represents /_cluster/pending_tasks
type ClusterPendingTasks struct {
	Tasks []map[string]interface{} `json:"tasks"`
}

type Index struct {
	Name             string `json:"index"`
	Health           string `json:"health"`
	Status           string `json:"status"`
	UUID             string `json:"uuid"`
	Primary          string `json:"pri"`
	Replica          string `json:"rep"`
	DocsCount        string `json:"docs.count"`
	DocsDeleted      string `json:"docs.deleted"`
	StoreSize        string `json:"store.size"`
	PrimaryStoreSize string `json:"pri.store.size"`
}

type Node struct {
	Name        string `json:"name"`
	Host        string `json:"host"`
	IP          string `json:"ip"`
	HeapPercent string `json:"heap.percent"`
	RAMPercent  string `json:"ram.percent"`
	CPU         string `json:"cpu"`
	Load1m      string `json:"load_1m"`
	Load5m      string `json:"load_5m"`
	Load15m     string `json:"load_15m"`
	NodeRole    string `json:"node.role"`
	Master      string `json:"master"`
}

type DataStream struct {
	Name               string             `json:"name"`
	TimestampField     TimestampFieldType `json:"timestamp_field"`
	Indices            []DataStreamIndex  `json:"indices"`
	Generation         int                `json:"generation"`
	Status             string             `json:"status"`
	Template           string             `json:"template,omitempty"`
	IlmPolicy          string             `json:"ilm_policy,omitempty"`
	Hidden             bool               `json:"hidden,omitempty"`
	System             bool               `json:"system,omitempty"`
	AllowCustomRouting bool               `json:"allow_custom_routing,omitempty"`
}

type TimestampFieldType struct {
	Name string `json:"name"`
}

type DataStreamIndex struct {
	IndexName string `json:"index_name"`
	IndexUUID string `json:"index_uuid"`
	PreferILM bool   `json:"prefer_ilm,omitempty"`
	ManagedBy string `json:"managed_by,omitempty"`
}

type RolloverResponse struct {
	Acknowledged       bool            `json:"acknowledged"`
	ShardsAcknowledged bool            `json:"shards_acknowledged"`
	OldIndex           string          `json:"old_index"`
	NewIndex           string          `json:"new_index"`
	RolledOver         bool            `json:"rolled_over"`
	DryRun             bool            `json:"dry_run"`
	Conditions         map[string]bool `json:"conditions"`
}

type IndexTemplate struct {
	Name         string                 `json:"name"`
	IndexPattern []string               `json:"index_patterns"`
	Template     TemplateDefinition     `json:"template,omitempty"`
	ComposedOf   []string               `json:"composed_of,omitempty"`
	Priority     int                    `json:"priority,omitempty"`
	Version      int                    `json:"version,omitempty"`
	Meta         map[string]interface{} `json:"_meta,omitempty"`
	DataStream   map[string]interface{} `json:"data_stream,omitempty"`
}

type TemplateDefinition struct {
	Settings map[string]interface{} `json:"settings,omitempty"`
	Mappings map[string]interface{} `json:"mappings,omitempty"`
	Aliases  map[string]interface{} `json:"aliases,omitempty"`
}

type ComponentTemplate struct {
	Name     string                 `json:"name"`
	Template TemplateDefinition     `json:"template,omitempty"`
	Version  int                    `json:"version,omitempty"`
	Meta     map[string]interface{} `json:"_meta,omitempty"`
}

type LifecyclePolicy struct {
	Name         string                 `json:"name"`
	Policy       map[string]interface{} `json:"policy"`
	Version      int                    `json:"version,omitempty"`
	ModifiedDate string                 `json:"modified_date,omitempty"`
}
