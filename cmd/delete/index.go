package delete

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/chronicblondiee/searchctl/pkg/client"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// getMatchingIndices returns a list of indices matching the given pattern
func getMatchingIndices(c client.SearchClient, pattern string) ([]string, error) {
	// Use the Get indices functionality to list all indices
	indices, err := c.GetIndices("*") // Get all indices first
	if err != nil {
		return nil, err
	}

	var matches []string
	// Convert wildcard pattern to simple string matching
	// For now, we'll do basic prefix matching for patterns ending with *
	if strings.HasSuffix(pattern, "*") {
		prefix := strings.TrimSuffix(pattern, "*")
		for _, index := range indices {
			if strings.HasPrefix(index.Name, prefix) {
				matches = append(matches, index.Name)
			}
		}
	} else {
		// Exact match (shouldn't happen in wildcard path but handle it)
		for _, index := range indices {
			if index.Name == pattern {
				matches = append(matches, index.Name)
			}
		}
	}

	return matches, nil
}

// confirmAction prompts the user for confirmation unless -y flag is set
func confirmAction(cmd *cobra.Command, action string) bool {
	// Check if -y flag is set
	if yes, _ := cmd.Flags().GetBool("yes"); yes {
		return true
	}

	fmt.Printf("Are you sure you want to %s? (y/N): ", action)
	reader := bufio.NewReader(os.Stdin)
	response, err := reader.ReadString('\n')
	if err != nil {
		return false
	}

	response = strings.TrimSpace(strings.ToLower(response))
	return response == "y" || response == "yes"
}

func NewDeleteIndexCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "index INDEX_NAME_OR_PATTERN",
		Short:   "Delete an index or indices matching a pattern",
		Long:    "Delete an index or indices matching a pattern from the search cluster. Supports wildcards like 'logs-*'.",
		Aliases: []string{"idx"},
		Args:    cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			indexPattern := args[0]

			// Handle dry-run mode
			if viper.GetBool("dry-run") {
				if strings.Contains(indexPattern, "*") {
					cmd.Printf("Would delete indices matching pattern: %s\n", indexPattern)
				} else {
					cmd.Printf("Would delete index: %s\n", indexPattern)
				}
				return
			}

			c, err := client.NewClient()
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error creating client: %v\n", err)
				os.Exit(1)
			}

			// Check if input contains wildcard
			if strings.Contains(indexPattern, "*") {
				fmt.Printf("Wildcard pattern detected: %s\n", indexPattern)

				// Get list of matching indices
				indices, err := getMatchingIndices(c, indexPattern)
				if err != nil {
					fmt.Fprintf(os.Stderr, "Error listing matching indices: %v\n", err)
					os.Exit(1)
				}

				if len(indices) == 0 {
					fmt.Printf("No indices match pattern: %s\n", indexPattern)
					return
				}

				// Show what will be deleted
				fmt.Printf("Found %d matching indices:\n", len(indices))
				for _, idx := range indices {
					fmt.Printf("  - %s\n", idx)
				}

				if !confirmAction(cmd, fmt.Sprintf("delete %d indices matching pattern '%s'", len(indices), indexPattern)) {
					fmt.Println("Delete operation cancelled.")
					return
				}

				// Delete each index individually
				var errors []string
				for _, idx := range indices {
					fmt.Printf("Deleting index: %s\n", idx)
					if err := c.DeleteIndex(idx); err != nil {
						errors = append(errors, fmt.Sprintf("failed to delete %s: %v", idx, err))
					} else {
						fmt.Printf("Successfully deleted index: %s\n", idx)
					}
				}

				if len(errors) > 0 {
					fmt.Fprintf(os.Stderr, "Errors occurred during deletion:\n%s\n", strings.Join(errors, "\n"))
					os.Exit(1)
				}

				fmt.Printf("All matching indices deleted successfully\n")
			} else {
				if !confirmAction(cmd, fmt.Sprintf("delete index '%s'", indexPattern)) {
					fmt.Println("Delete operation cancelled.")
					return
				}

				if err := c.DeleteIndex(indexPattern); err != nil {
					fmt.Fprintf(os.Stderr, "Error deleting index: %v\n", err)
					os.Exit(1)
				}

				fmt.Printf("Index %s deleted successfully\n", indexPattern)
			}
		},
	}

	cmd.Flags().BoolP("yes", "y", false, "automatically confirm deletion without prompting")

	return cmd
}
