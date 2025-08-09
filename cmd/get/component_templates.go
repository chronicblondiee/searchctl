package get

import (
	"fmt"
	"os"

	"github.com/chronicblondiee/searchctl/pkg/client"
	"github.com/chronicblondiee/searchctl/pkg/output"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func NewGetComponentTemplatesCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "component-templates [PATTERN]",
		Short:   "Get component templates",
		Long:    "Get component templates from the search cluster.",
		Aliases: []string{"componenttemplates", "component-template", "componenttemplate", "ct", "comp-templates", "comp-template"},
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

			templates, err := c.GetComponentTemplates(pattern)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error getting component templates: %v\n", err)
				os.Exit(1)
			}

			// Convert to interface{} slice for formatting
			data := make([]interface{}, len(templates))
			for i, template := range templates {
				data[i] = map[string]interface{}{
					"NAME":    template.Name,
					"VERSION": template.Version,
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