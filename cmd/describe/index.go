package describe

import (
	"fmt"
	"os"

	"github.com/chronicblondiee/searchctl/pkg/client"
	"github.com/chronicblondiee/searchctl/pkg/output"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func NewDescribeIndexCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "index INDEX_NAME",
		Short:   "Describe an index",
		Long:    "Show detailed information about a specific index.",
		Aliases: []string{"idx"},
		Args:    cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			indexName := args[0]

			c, err := client.NewClient()
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error creating client: %v\n", err)
				os.Exit(1)
			}

			index, err := c.GetIndex(indexName)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error getting index: %v\n", err)
				os.Exit(1)
			}

			data := map[string]interface{}{
				"Name":               index.Name,
				"Health":             index.Health,
				"Status":             index.Status,
				"UUID":               index.UUID,
				"Primary Shards":     index.Primary,
				"Replica Shards":     index.Replica,
				"Documents Count":    index.DocsCount,
				"Documents Deleted":  index.DocsDeleted,
				"Store Size":         index.StoreSize,
				"Primary Store Size": index.PrimaryStoreSize,
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
