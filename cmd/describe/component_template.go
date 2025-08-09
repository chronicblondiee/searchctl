package describe

import (
	"fmt"
	"os"

	"github.com/chronicblondiee/searchctl/pkg/client"
	"github.com/chronicblondiee/searchctl/pkg/output"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func NewDescribeComponentTemplateCmd() *cobra.Command {
	var showBody bool

	cmd := &cobra.Command{
		Use:     "component-template NAME",
		Short:   "Describe a component template",
		Long:    "Show detailed information about a specific component template.",
		Aliases: []string{"componenttemplates", "component-template", "componenttemplate", "ct", "comp-templates", "comp-template"},
		Args:    cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			name := args[0]

			c, err := client.NewClient()
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error creating client: %v\n", err)
				os.Exit(1)
			}

			ct, err := c.GetComponentTemplate(name)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error getting component template: %v\n", err)
				os.Exit(1)
			}

			outFmt := viper.GetString("output")
			if outFmt == "json" || outFmt == "yaml" {
				formatter := output.NewFormatter(outFmt)
				if err := formatter.Format(ct, os.Stdout); err != nil {
					fmt.Fprintf(os.Stderr, "Error formatting output: %v\n", err)
					os.Exit(1)
				}
				return
			}

			data := map[string]interface{}{
				"Name":    ct.Name,
				"Version": ct.Version,
			}

			if showBody {
				data["Template"] = ct.Template
				if len(ct.Meta) > 0 {
					data["Meta"] = ct.Meta
				}
			}

			formatter := output.NewFormatter(outFmt)
			if err := formatter.Format(data, os.Stdout); err != nil {
				fmt.Fprintf(os.Stderr, "Error formatting output: %v\n", err)
				os.Exit(1)
			}
		},
	}

	cmd.Flags().BoolVar(&showBody, "show-body", false, "include full template body in table output")

	return cmd
}
