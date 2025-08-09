package describe

import (
	"fmt"
	"os"

	"github.com/chronicblondiee/searchctl/pkg/client"
	"github.com/chronicblondiee/searchctl/pkg/output"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func NewDescribeNodeCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "node NODE_ID",
		Short:   "Describe a node",
		Long:    "Show detailed information about a specific node by name or IP.",
		Aliases: []string{"node", "no"},
		Args:    cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			nodeID := args[0]

			c, err := client.NewClient()
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error creating client: %v\n", err)
				os.Exit(1)
			}

			node, err := c.GetNode(nodeID)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error getting node: %v\n", err)
				os.Exit(1)
			}

			outFmt := viper.GetString("output")
			if outFmt == "json" || outFmt == "yaml" {
				formatter := output.NewFormatter(outFmt)
				if err := formatter.Format(node, os.Stdout); err != nil {
					fmt.Fprintf(os.Stderr, "Error formatting output: %v\n", err)
					os.Exit(1)
				}
				return
			}

			data := map[string]interface{}{
				"Name":        node.Name,
				"Host":        node.Host,
				"IP":          node.IP,
				"NodeRole":    node.NodeRole,
				"Master":      node.Master,
				"CPU":         node.CPU,
				"RAMPercent":  node.RAMPercent,
				"HeapPercent": node.HeapPercent,
				"Load1m":      node.Load1m,
				"Load5m":      node.Load5m,
				"Load15m":     node.Load15m,
			}

			formatter := output.NewFormatter(outFmt)
			if err := formatter.Format(data, os.Stdout); err != nil {
				fmt.Fprintf(os.Stderr, "Error formatting output: %v\n", err)
				os.Exit(1)
			}
		},
	}

	return cmd
}
