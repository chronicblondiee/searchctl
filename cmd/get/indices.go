package get

import (
	"fmt"
	"os"

	"github.com/chronicblondiee/searchctl/pkg/client"
	"github.com/chronicblondiee/searchctl/pkg/output"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func NewGetIndicesCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "indices [INDEX_PATTERN]",
		Short:   "List indices",
		Long:    "List all indices or indices matching a pattern.",
		Aliases: []string{"index", "idx"},
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

			indices, err := c.GetIndices(pattern)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error getting indices: %v\n", err)
				os.Exit(1)
			}

			// Convert to interface{} slice for formatting
			data := make([]interface{}, len(indices))
			for i, idx := range indices {
				data[i] = map[string]interface{}{
					"NAME":       idx.Name,
					"HEALTH":     idx.Health,
					"STATUS":     idx.Status,
					"UUID":       idx.UUID,
					"PRI":        idx.Primary,
					"REP":        idx.Replica,
					"DOCS.COUNT": idx.DocsCount,
					"STORE.SIZE": idx.StoreSize,
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
