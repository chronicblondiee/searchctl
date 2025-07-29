package get

import (
	"fmt"
	"os"

	"github.com/chronicblondiee/searchctl/pkg/client"
	"github.com/chronicblondiee/searchctl/pkg/output"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func NewGetDataStreamsCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "datastreams [PATTERN]",
		Short:   "List data streams",
		Long:    "List all data streams or data streams matching a pattern.",
		Aliases: []string{"datastream", "ds"},
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

			dataStreams, err := c.GetDataStreams(pattern)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error getting data streams: %v\n", err)
				os.Exit(1)
			}

			// Convert to interface{} slice for formatting
			data := make([]interface{}, len(dataStreams))
			for i, ds := range dataStreams {
				indicesCount := len(ds.Indices)
				indicesNames := make([]string, len(ds.Indices))
				for j, idx := range ds.Indices {
					indicesNames[j] = idx.IndexName
				}
				
				data[i] = map[string]interface{}{
					"NAME":       ds.Name,
					"STATUS":     ds.Status,
					"INDICES":    indicesCount,
					"GENERATION": ds.Generation,
					"TEMPLATE":   ds.Template,
					"TIMESTAMP":  ds.TimestampField.Name,
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
