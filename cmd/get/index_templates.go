package get

import (
	"fmt"
	"os"

	"github.com/chronicblondiee/searchctl/pkg/client"
	"github.com/chronicblondiee/searchctl/pkg/output"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func NewGetIndexTemplatesCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "index-templates [PATTERN]",
		Short:   "Get index templates",
		Long:    "Get index templates from the search cluster.",
		Aliases: []string{"idx-templates", "template", "it", "index-template", "indextemplates", "indextemplate"},
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

			templates, err := c.GetIndexTemplates(pattern)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error getting index templates: %v\n", err)
				os.Exit(1)
			}

			// Convert to interface{} slice for formatting
			data := make([]interface{}, len(templates))
			for i, template := range templates {
				data[i] = map[string]interface{}{
					"NAME":     template.Name,
					"PATTERNS": template.IndexPattern,
					"PRIORITY": template.Priority,
					"VERSION":  template.Version,
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
