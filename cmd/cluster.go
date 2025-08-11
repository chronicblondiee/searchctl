package cmd

import (
	"fmt"
	"os"
	"strings"

	cluster "github.com/chronicblondiee/searchctl/cmd/cluster"
	"github.com/chronicblondiee/searchctl/pkg/client"
	"github.com/chronicblondiee/searchctl/pkg/output"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func NewClusterCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "cluster",
		Short: "Cluster operations",
		Long:  "Perform cluster-wide operations.",
	}

	cmd.AddCommand(NewClusterHealthCmd())
	cmd.AddCommand(NewClusterInfoCmd())
	cmd.AddCommand(NewClusterStatsCmd())
	cmd.AddCommand(NewClusterStateCmd())
	cmd.AddCommand(NewClusterPendingTasksCmd())
	cmd.AddCommand(cluster.NewAllocationSettingsCmd())

	return cmd
}

func NewClusterHealthCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "health",
		Short: "Show cluster health",
		Long:  "Display the health status of the cluster.",
		Run: func(cmd *cobra.Command, args []string) {
			c, err := client.NewClient()
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error creating client: %v\n", err)
				os.Exit(1)
			}

			health, err := c.ClusterHealth()
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error getting cluster health: %v\n", err)
				os.Exit(1)
			}

			data := map[string]interface{}{
				"Cluster Name":          health.ClusterName,
				"Status":                health.Status,
				"Timed Out":             health.TimedOut,
				"Number of Nodes":       health.NumberOfNodes,
				"Number of Data Nodes":  health.NumberOfDataNodes,
				"Active Primary Shards": health.ActivePrimaryShards,
				"Active Shards":         health.ActiveShards,
				"Relocating Shards":     health.RelocatingShards,
				"Initializing Shards":   health.InitializingShards,
				"Unassigned Shards":     health.UnassignedShards,
			}

			formatter := output.NewFormatter(viper.GetString("output"))
			if err := formatter.Format(data, os.Stdout); err != nil {
				fmt.Fprintf(os.Stderr, "Error formatting output: %v\n", err)
				os.Exit(1)
			}
		},
	}

	return cmd
}

func NewClusterInfoCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "info",
		Short: "Show cluster information",
		Long:  "Display general information about the cluster.",
		Run: func(cmd *cobra.Command, args []string) {
			c, err := client.NewClient()
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error creating client: %v\n", err)
				os.Exit(1)
			}

			info, err := c.ClusterInfo()
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error getting cluster info: %v\n", err)
				os.Exit(1)
			}

			data := map[string]interface{}{
				"Name":         info.Name,
				"Cluster Name": info.ClusterName,
				"Cluster UUID": info.ClusterUUID,
				"Version":      info.Version,
				"Tagline":      info.Tagline,
			}

			formatter := output.NewFormatter(viper.GetString("output"))
			if err := formatter.Format(data, os.Stdout); err != nil {
				fmt.Fprintf(os.Stderr, "Error formatting output: %v\n", err)
				os.Exit(1)
			}
		},
	}

	return cmd
}

func NewClusterStatsCmd() *cobra.Command {
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
			stats, err := c.GetClusterStats()
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error getting cluster stats: %v\n", err)
				os.Exit(1)
			}
			// Summarize useful fields if present
			summary := map[string]interface{}{
				"Cluster Name": stats.ClusterName,
			}
			if nodes, ok := stats.Nodes["count"].(map[string]interface{}); ok {
				if total, ok := nodes["total"]; ok {
					summary["Nodes"] = total
				}
				if data, ok := nodes["data"]; ok {
					summary["Data Nodes"] = data
				}
			}
			if indices, ok := stats.Indices["shards"].(map[string]interface{}); ok {
				if total, ok := indices["total"]; ok {
					summary["Shards"] = total
				}
			}
			formatter := output.NewFormatter(viper.GetString("output"))
			if err := formatter.Format(summary, os.Stdout); err != nil {
				fmt.Fprintf(os.Stderr, "Error formatting output: %v\n", err)
				os.Exit(1)
			}
		},
	}
	return cmd
}

func NewClusterStateCmd() *cobra.Command {
	var metrics string
	var indices string
	var masterTimeout string
	cmd := &cobra.Command{
		Use:   "state",
		Short: "Show cluster state",
		Long:  "Display cluster state with optional metric and index filtering.",
		Run: func(cmd *cobra.Command, args []string) {
			c, err := client.NewClient()
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error creating client: %v\n", err)
				os.Exit(1)
			}
			var metricList []string
			if metrics != "" {
				metricList = splitAndTrim(metrics)
			}
			st, err := c.GetClusterState(metricList, indices, masterTimeout)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error getting cluster state: %v\n", err)
				os.Exit(1)
			}
			data := map[string]interface{}{
				"Cluster Name": st.ClusterName,
				"State UUID":   st.StateUUID,
			}
			formatter := output.NewFormatter(viper.GetString("output"))
			if err := formatter.Format(data, os.Stdout); err != nil {
				fmt.Fprintf(os.Stderr, "Error formatting output: %v\n", err)
				os.Exit(1)
			}
		},
	}
	cmd.Flags().StringVar(&metrics, "metrics", "", "comma-separated metrics (e.g. metadata,routing_table,blocks,nodes)")
	cmd.Flags().StringVar(&indices, "indices", "", "indices filter for state")
	cmd.Flags().StringVar(&masterTimeout, "master-timeout", "", "timeout for connecting to master (e.g. 30s)")
	return cmd
}

func NewClusterPendingTasksCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "pending-tasks",
		Short: "Show cluster pending tasks",
		Long:  "Display cluster pending tasks.",
		Run: func(cmd *cobra.Command, args []string) {
			c, err := client.NewClient()
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error creating client: %v\n", err)
				os.Exit(1)
			}
			pt, err := c.GetPendingTasks()
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error getting pending tasks: %v\n", err)
				os.Exit(1)
			}
			data := map[string]interface{}{
				"Tasks": pt.Tasks,
				"Count": len(pt.Tasks),
			}
			formatter := output.NewFormatter(viper.GetString("output"))
			if err := formatter.Format(data, os.Stdout); err != nil {
				fmt.Fprintf(os.Stderr, "Error formatting output: %v\n", err)
				os.Exit(1)
			}
		},
	}
	return cmd
}

// splitAndTrim splits by comma and trims whitespace.
func splitAndTrim(s string) []string {
	var out []string
	for _, part := range strings.Split(s, ",") {
		p := strings.TrimSpace(part)
		if p != "" {
			out = append(out, p)
		}
	}
	return out
}
