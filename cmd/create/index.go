package create

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/chronicblondiee/searchctl/pkg/client"
)

func NewCreateIndexCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "index INDEX_NAME",
		Short:   "Create an index",
		Long:    "Create a new index in the search cluster.",
		Aliases: []string{"idx"},
		Args:    cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			indexName := args[0]

			if viper.GetBool("dry-run") {
				fmt.Printf("Would create index: %s\n", indexName)
				return
			}

			c, err := client.NewClient()
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error creating client: %v\n", err)
				os.Exit(1)
			}

			if err := c.CreateIndex(indexName, nil); err != nil {
				fmt.Fprintf(os.Stderr, "Error creating index: %v\n", err)
				os.Exit(1)
			}

			fmt.Printf("Index %s created successfully\n", indexName)
		},
	}

	return cmd
}
