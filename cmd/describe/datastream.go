package describe

import (
	"fmt"
	"os"

	"github.com/chronicblondiee/searchctl/pkg/client"
	"github.com/chronicblondiee/searchctl/pkg/output"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func NewDescribeDataStreamCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "datastream NAME",
		Short:   "Describe a data stream",
		Long:    "Show detailed information about a specific data stream.",
		Aliases: []string{"datastreams", "ds"},
		Args:    cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			name := args[0]

			c, err := client.NewClient()
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error creating client: %v\n", err)
				os.Exit(1)
			}

			ds, err := c.GetDataStream(name)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error getting data stream: %v\n", err)
				os.Exit(1)
			}

			outFmt := viper.GetString("output")
			if outFmt == "json" || outFmt == "yaml" {
				formatter := output.NewFormatter(outFmt)
				if err := formatter.Format(ds, os.Stdout); err != nil {
					fmt.Fprintf(os.Stderr, "Error formatting output: %v\n", err)
					os.Exit(1)
				}
				return
			}

			indices := make([]map[string]interface{}, 0, len(ds.Indices))
			for _, idx := range ds.Indices {
				indices = append(indices, map[string]interface{}{
					"IndexName": idx.IndexName,
					"IndexUUID": idx.IndexUUID,
					"PreferILM": idx.PreferILM,
					"ManagedBy": idx.ManagedBy,
				})
			}

			data := map[string]interface{}{
				"Name":           ds.Name,
				"Generation":     ds.Generation,
				"Status":         ds.Status,
				"TimestampField": ds.TimestampField,
				"Template":       ds.Template,
				"Hidden":         ds.Hidden,
				"System":         ds.System,
				"ILMPolicy":      ds.IlmPolicy,
				"Indices":        indices,
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
