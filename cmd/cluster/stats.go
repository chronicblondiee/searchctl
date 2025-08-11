package cluster

import (
	"fmt"
	"os"

	"github.com/chronicblondiee/searchctl/pkg/client"
	"github.com/chronicblondiee/searchctl/pkg/output"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func NewStatsCmd() *cobra.Command {
	var raw bool
	cmd := &cobra.Command{
		Use:   "stats",
		Short: "Show cluster statistics",
		Long:  "Display cluster statistics summary.",
		Run: func(cmd *cobra.Command, args []string) {
			c, err := client.NewClient()
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error creating client: %v\n", err)
				os.Exit(1)
			}
			stats, err := c.ClusterStats()
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error getting cluster stats: %v\n", err)
				os.Exit(1)
			}
			if raw {
				formatter := output.NewFormatter(viper.GetString("output"))
				if err := formatter.Format(stats, os.Stdout); err != nil {
					fmt.Fprintf(os.Stderr, "Error formatting output: %v\n", err)
					os.Exit(1)
				}
				return
			}
			summary := map[string]interface{}{
				"Cluster Name": stats.ClusterName,
			}
			if nodesCount, ok := nested(stats.Nodes, "count").(map[string]interface{}); ok {
				if v, ok := nodesCount["total"]; ok {
					summary["Nodes"] = v
				}
				if v, ok := nodesCount["data"]; ok {
					summary["Data Nodes"] = v
				}
			}
			if v := nested(stats.Indices, "count"); v != nil {
				summary["Indices"] = v
			}
			if shards, ok := nested(stats.Indices, "shards").(map[string]interface{}); ok {
				if v, ok := shards["total"]; ok {
					summary["Shards Total"] = v
				}
				if v, ok := shards["primaries"]; ok {
					summary["Shards Primaries"] = v
				}
			}
			if docs, ok := nested(stats.Indices, "docs").(map[string]interface{}); ok {
				if v, ok := docs["count"]; ok {
					summary["Docs Count"] = v
				}
			}
			if store, ok := nested(stats.Indices, "store").(map[string]interface{}); ok {
				if v, ok := store["size_in_bytes"]; ok {
					summary["Store"] = formatBytes(v)
				}
			}
			if fs, ok := nested(stats.Nodes, "fs").(map[string]interface{}); ok {
				if v, ok := fs["total_in_bytes"]; ok {
					summary["FS Total"] = formatBytes(v)
				}
				if v, ok := fs["available_in_bytes"]; ok {
					summary["FS Available"] = formatBytes(v)
				}
			}
			if jvm, ok := nested(stats.Nodes, "jvm").(map[string]interface{}); ok {
				if mem, ok := jvm["mem"].(map[string]interface{}); ok {
					if v, ok := mem["heap_used_in_bytes"]; ok {
						summary["JVM Heap Used"] = formatBytes(v)
					}
					if v, ok := mem["heap_max_in_bytes"]; ok {
						summary["JVM Heap Max"] = formatBytes(v)
					}
				}
			}
			formatter := output.NewFormatter(viper.GetString("output"))
			if err := formatter.Format(summary, os.Stdout); err != nil {
				fmt.Fprintf(os.Stderr, "Error formatting output: %v\n", err)
				os.Exit(1)
			}
		},
	}
	cmd.Flags().BoolVar(&raw, "raw", false, "output full cluster stats payload")
	return cmd
}

func nested(m map[string]interface{}, key string) interface{} {
	if m == nil {
		return nil
	}
	if v, ok := m[key]; ok {
		return v
	}
	return nil
}

func formatBytes(v interface{}) string {
	bytesVal := toFloat64(v)
	units := []string{"B", "KB", "MB", "GB", "TB", "PB"}
	unitIdx := 0
	for bytesVal >= 1024.0 && unitIdx < len(units)-1 {
		bytesVal /= 1024.0
		unitIdx++
	}
	if unitIdx == 0 {
		return fmt.Sprintf("%.0f %s", bytesVal, units[unitIdx])
	}
	return fmt.Sprintf("%.1f %s", bytesVal, units[unitIdx])
}

func toFloat64(v interface{}) float64 {
	switch t := v.(type) {
	case int:
		return float64(t)
	case int32:
		return float64(t)
	case int64:
		return float64(t)
	case uint:
		return float64(t)
	case uint32:
		return float64(t)
	case uint64:
		return float64(t)
	case float32:
		return float64(t)
	case float64:
		return t
	case string:
		var f float64
		_, err := fmt.Sscanf(t, "%f", &f)
		if err == nil {
			return f
		}
		return 0
	default:
		return 0
	}
}
