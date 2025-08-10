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

func NewAllocationSettingsCmd() *cobra.Command {
	var enable string
	var rebalance string
	var awareness string
	var file string

	cmd := &cobra.Command{
		Use:   "allocation-settings",
		Short: "Get or set cluster shard allocation settings",
		Long:  "Get or set cluster shard allocation settings like enable, rebalance, and awareness attributes.",
		Run: func(cmd *cobra.Command, args []string) {
			c, err := client.NewClient()
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error creating client: %v\n", err)
				os.Exit(1)
			}

			if enable == "" && rebalance == "" && awareness == "" && file == "" {
				// GET
				settings, err := c.GetClusterSettings()
				if err != nil {
					fmt.Fprintf(os.Stderr, "Error getting settings: %v\n", err)
					os.Exit(1)
				}
				formatter := output.NewFormatter(viper.GetString("output"))
				if err := formatter.Format(settings, os.Stdout); err != nil {
					fmt.Fprintf(os.Stderr, "Error formatting output: %v\n", err)
					os.Exit(1)
				}
				return
			}

			// SET
			body := map[string]interface{}{"transient": map[string]interface{}{}, "persistent": map[string]interface{}{}}
			t := body["transient"].(map[string]interface{})

			if enable != "" {
				t["cluster.routing.allocation.enable"] = enable
			}
			if rebalance != "" {
				t["cluster.routing.rebalance.enable"] = rebalance
			}
			if awareness != "" {
				t["cluster.routing.allocation.awareness.attributes"] = strings.TrimSpace(awareness)
			}

			if err := c.UpdateClusterSettings(body); err != nil {
				fmt.Fprintf(os.Stderr, "Error updating settings: %v\n", err)
				os.Exit(1)
			}
			fmt.Fprintln(os.Stdout, "Cluster allocation settings updated")
		},
	}

	cmd.Flags().StringVar(&enable, "enable", "", "allocation enable (all|primaries|new_primaries|none)")
	cmd.Flags().StringVar(&rebalance, "rebalance", "", "rebalance enable (all|primaries|replicas|none)")
	cmd.Flags().StringVar(&awareness, "awareness-attrs", "", "allocation awareness attributes (comma-separated)")
	cmd.Flags().StringVarP(&file, "filename", "f", "", "settings file (JSON/YAML) [TODO: not yet implemented]")

	return cmd
}
