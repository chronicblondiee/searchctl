package cluster

import (
	"fmt"
	"os"
	"strings"

	"github.com/chronicblondiee/searchctl/pkg/client"
	"github.com/chronicblondiee/searchctl/pkg/output"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func NewStateCmd() *cobra.Command {
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
			st, err := c.ClusterState(metricList, indices, masterTimeout)
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
