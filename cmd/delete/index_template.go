package delete

import (
	"fmt"
	"os"

	"github.com/chronicblondiee/searchctl/pkg/client"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func NewDeleteIndexTemplateCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "index-template TEMPLATE_NAME",
		Short:   "Delete an index template",
		Long:    "Delete an index template from the search cluster.",
		Aliases: []string{"idx-templates", "template", "it", "index-template", "indextemplates", "indextemplate"},
		Args:    cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			templateName := args[0]

			if viper.GetBool("dry-run") {
				cmd.Printf("Would delete index template: %s\n", templateName)
				return
			}

			c, err := client.NewClient()
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error creating client: %v\n", err)
				os.Exit(1)
			}

			if err := c.DeleteIndexTemplate(templateName); err != nil {
				fmt.Fprintf(os.Stderr, "Error deleting index template: %v\n", err)
				os.Exit(1)
			}

			cmd.Printf("Index template %s deleted successfully\n", templateName)
		},
	}

	return cmd
}
