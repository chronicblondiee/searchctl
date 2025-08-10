package get

import (
	"fmt"
	"os"

	"github.com/chronicblondiee/searchctl/pkg/client"
	"github.com/chronicblondiee/searchctl/pkg/output"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func NewGetShardsCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "shards [INDEX_PATTERN]",
		Short:   "List shard allocations",
		Long:    "List shard allocations for the cluster or matching indices.",
		Aliases: []string{"shard"},
		Args:    cobra.MaximumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			pattern := ""
			if len(args) > 0 {
				pattern = args[0]
			}

			c, err := client.NewClient()
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error creating client: %v\n", err)
				os.Exit(1)
			}

			rows, err := c.GetShards(pattern)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error getting shards: %v\n", err)
				os.Exit(1)
			}

			// Convert rows to generic maps for table/formatters
			data := make([]interface{}, len(rows))
			for i, r := range rows {
				data[i] = map[string]interface{}{
					"INDEX":    r.Index,
					"SHARD":    r.Shard,
					"PRI/REP":  r.PrimaryOrReplica,
					"STATE":    r.State,
					"DOCS":     r.Docs,
					"STORE":    r.Store,
					"IP":       r.IP,
					"NODE":     r.Node,
					"UNASSIGN": r.UnassignedReason,
				}
			}

			formatter := output.NewFormatter(viper.GetString("output"))
			if err := formatter.Format(data, cmd.OutOrStdout()); err != nil {
				fmt.Fprintf(os.Stderr, "Error formatting output: %v\n", err)
				os.Exit(1)
			}
		},
	}

	return cmd
}
