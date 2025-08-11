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
	var raw bool
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
			if raw {
				formatter := output.NewFormatter(viper.GetString("output"))
				if err := formatter.Format(st, os.Stdout); err != nil {
					fmt.Fprintf(os.Stderr, "Error formatting output: %v\n", err)
					os.Exit(1)
				}
				return
			}
			metricsOut := "all"
			if len(metricList) > 0 {
				metricsOut = strings.Join(metricList, ",")
			}
			// Build summary based on requested metrics (or common defaults)
			selected := map[string]bool{}
			if len(metricList) == 0 {
				selected["metadata"] = true
				selected["routing_table"] = true
				selected["nodes"] = true
				selected["blocks"] = true
			} else {
				for _, m := range metricList {
					selected[m] = true
				}
			}

			data := map[string]interface{}{
				"Cluster Name": st.ClusterName,
				"State UUID":   st.StateUUID,
				"Metrics":      metricsOut,
			}

			if selected["metadata"] {
				if mdIndices, ok := st.Metadata["indices"].(map[string]interface{}); ok {
					data["Metadata Indices"] = len(mdIndices)
				}
			}
			if selected["routing_table"] {
				if rtIndices, ok := st.RoutingTable["indices"].(map[string]interface{}); ok {
					data["Routing Indices"] = len(rtIndices)
					totalShards := 0
					primaries := 0
					replicas := 0
					for _, idxVal := range rtIndices {
						if idxMap, ok := idxVal.(map[string]interface{}); ok {
							if shards, ok := idxMap["shards"].(map[string]interface{}); ok {
								for _, shardVal := range shards {
									if allocs, ok := shardVal.([]interface{}); ok {
										totalShards += len(allocs)
										for _, alloc := range allocs {
											if allocMap, ok := alloc.(map[string]interface{}); ok {
												if p, ok := allocMap["primary"].(bool); ok && p {
													primaries++
												} else {
													replicas++
												}
											}
										}
									}
								}
							}
						}
					}
					data["Routing Shards Total"] = totalShards
					data["Routing Shards Primaries"] = primaries
					data["Routing Shards Replicas"] = replicas
				}
			}
			if selected["nodes"] {
				// The state nodes map contains node entries by ID
				data["Nodes In State"] = len(st.Nodes)
			}
			if selected["blocks"] {
				if indicesBlocks, ok := st.Blocks["indices"].(map[string]interface{}); ok {
					data["Blocked Indices"] = len(indicesBlocks)
					names := make([]string, 0, len(indicesBlocks))
					for name := range indicesBlocks {
						names = append(names, name)
					}
					if len(names) > 10 {
						names = names[:10]
					}
					if len(names) > 0 {
						data["Blocked Index Names"] = strings.Join(names, ",")
					}
				}
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
	cmd.Flags().BoolVar(&raw, "raw", false, "output full cluster state payload")
	return cmd
}

func splitAndTrim(s string) []string {
	var out []string
	for _, part := range strings.Split(s, ",") {
		p := strings.TrimSpace(part)
		if p != "" {
			p = strings.ReplaceAll(p, "-", "_")
		}
		if p != "" {
			out = append(out, p)
		}
	}
	return out
}
