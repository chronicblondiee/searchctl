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