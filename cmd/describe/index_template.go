package describe

import (
	"fmt"
	"os"

	"github.com/chronicblondiee/searchctl/pkg/client"
	"github.com/chronicblondiee/searchctl/pkg/output"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func NewDescribeIndexTemplateCmd() *cobra.Command {
	var showBody bool

	cmd := &cobra.Command{
		Use:     "index-template NAME",
		Short:   "Describe an index template",
		Long:    "Show detailed information about a specific composable index template.",
		Aliases: []string{"idx-templates", "template", "it", "index-template", "indextemplates", "indextemplate"},
		Args:    cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			name := args[0]

			c, err := client.NewClient()
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error creating client: %v\n", err)
				os.Exit(1)
			}

			tmpl, err := c.GetIndexTemplate(name)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error getting index template: %v\n", err)
				os.Exit(1)
			}

			outFmt := viper.GetString("output")
			if outFmt == "json" || outFmt == "yaml" {
				formatter := output.NewFormatter(outFmt)
				if err := formatter.Format(tmpl, os.Stdout); err != nil {
					fmt.Fprintf(os.Stderr, "Error formatting output: %v\n", err)
					os.Exit(1)
				}
				return
			}

			data := map[string]interface{}{
				"Name":          tmpl.Name,
				"IndexPatterns": tmpl.IndexPattern,
				"Priority":      tmpl.Priority,
				"Version":       tmpl.Version,
				"ComposedOf":    tmpl.ComposedOf,
			}

			if showBody {
				data["Template"] = tmpl.Template
				if len(tmpl.Meta) > 0 {
					data["Meta"] = tmpl.Meta
				}
				if len(tmpl.DataStream) > 0 {
					data["DataStream"] = tmpl.DataStream
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
