package get

import (
	"fmt"
	"os"

	"github.com/chronicblondiee/searchctl/pkg/client"
	"github.com/chronicblondiee/searchctl/pkg/output"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func NewGetNodesCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "nodes [NODE_NAME]",
		Short:   "List nodes",
		Long:    "List all nodes in the cluster.",
		Aliases: []string{"node", "no"},
		Run: func(cmd *cobra.Command, args []string) {
			c, err := client.NewClient()
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error creating client: %v\n", err)
				os.Exit(1)
			}

			nodes, err := c.GetNodes()
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error getting nodes: %v\n", err)
				os.Exit(1)
			}

			// Convert to interface{} slice for formatting
			data := make([]interface{}, len(nodes))
			for i, node := range nodes {
				data[i] = map[string]interface{}{
					"NAME":         node.Name,
					"HOST":         node.Host,
					"IP":           node.IP,
					"HEAP.PERCENT": node.HeapPercent,
					"RAM.PERCENT":  node.RAMPercent,
					"CPU":          node.CPU,
					"LOAD_1M":      node.Load1m,
					"ROLE":         node.NodeRole,
					"MASTER":       node.Master,
				}
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
