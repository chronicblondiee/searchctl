package cmd

import (
	"fmt"
	"os"

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
