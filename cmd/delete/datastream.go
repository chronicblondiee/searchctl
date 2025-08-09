package delete

import (
	"fmt"
	"os"
	"strings"

	"github.com/chronicblondiee/searchctl/pkg/client"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// getMatchingDataStreams returns a list of data streams matching the given pattern
func getMatchingDataStreams(c client.SearchClient, pattern string) ([]string, error) {
	// Get all data streams first
	dataStreams, err := c.GetDataStreams("*")
	if err != nil {
		return nil, err
	}

	var matches []string
	// Convert wildcard pattern to simple string matching
	// For now, we'll do basic prefix matching for patterns ending with *
	if strings.HasSuffix(pattern, "*") {
		prefix := strings.TrimSuffix(pattern, "*")
		for _, ds := range dataStreams {
			if strings.HasPrefix(ds.Name, prefix) {
				matches = append(matches, ds.Name)
			}
		}
	} else {
		// Exact match (shouldn't happen in wildcard path but handle it)
		for _, ds := range dataStreams {
			if ds.Name == pattern {
				matches = append(matches, ds.Name)
			}
		}
	}

	return matches, nil
}

func NewDeleteDataStreamCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "datastream DATA_STREAM_NAME_OR_PATTERN",
		Short:   "Delete a data stream or data streams matching a pattern",
		Long:    "Delete a data stream or data streams matching a pattern and all their backing indices from the search cluster. Supports wildcards like 'logs-*'.",
		Aliases: []string{"datastream", "datastreams", "ds"},
		Args:    cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			dataStreamPattern := args[0]

			// Handle dry-run mode
			if viper.GetBool("dry-run") {
				if strings.Contains(dataStreamPattern, "*") {
					cmd.Printf("Would delete data streams matching pattern: %s\n", dataStreamPattern)
				} else {
					cmd.Printf("Would delete data stream: %s\n", dataStreamPattern)
				}
				return
			}

			c, err := client.NewClient()
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error creating client: %v\n", err)
				os.Exit(1)
			}

			// Check if input contains wildcard
			if strings.Contains(dataStreamPattern, "*") {
				fmt.Printf("Wildcard pattern detected: %s\n", dataStreamPattern)

				// Get list of matching data streams
				dataStreams, err := getMatchingDataStreams(c, dataStreamPattern)
				if err != nil {
					fmt.Fprintf(os.Stderr, "Error listing matching data streams: %v\n", err)
					os.Exit(1)
				}

				if len(dataStreams) == 0 {
					fmt.Printf("No data streams match pattern: %s\n", dataStreamPattern)
					return
				}

				// Show what will be deleted
				fmt.Printf("Found %d matching data streams:\n", len(dataStreams))
				for _, ds := range dataStreams {
					fmt.Printf("  - %s\n", ds)
				}

				if !confirmAction(cmd, fmt.Sprintf("delete %d data streams matching pattern '%s'", len(dataStreams), dataStreamPattern)) {
					fmt.Println("Delete operation cancelled.")
					return
				}

				// Delete each data stream individually
				var errors []string
				for _, ds := range dataStreams {
					fmt.Printf("Deleting data stream: %s\n", ds)
					if err := c.DeleteDataStream(ds); err != nil {
						errors = append(errors, fmt.Sprintf("failed to delete %s: %v", ds, err))
					} else {
						fmt.Printf("Successfully deleted data stream: %s\n", ds)
					}
				}

				if len(errors) > 0 {
					fmt.Fprintf(os.Stderr, "Errors occurred during deletion:\n%s\n", strings.Join(errors, "\n"))
					os.Exit(1)
				}

				fmt.Printf("All matching data streams deleted successfully\n")
			} else {
				if !confirmAction(cmd, fmt.Sprintf("delete data stream '%s'", dataStreamPattern)) {
					fmt.Println("Delete operation cancelled.")
					return
				}

				if err := c.DeleteDataStream(dataStreamPattern); err != nil {
					fmt.Fprintf(os.Stderr, "Error deleting data stream: %v\n", err)
					os.Exit(1)
				}

				fmt.Printf("Data stream %s deleted successfully\n", dataStreamPattern)
			}
		},
	}

	cmd.Flags().BoolP("yes", "y", false, "automatically confirm deletion without prompting")

	return cmd
}
