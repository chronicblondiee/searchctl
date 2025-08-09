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

func NewDeleteLifecyclePolicyCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "lifecycle-policy POLICY_NAME",
		Short:   "Delete a lifecycle policy",
		Long:    "Delete a lifecycle policy from the search cluster (ILM for Elasticsearch, ISM for OpenSearch).",
		Aliases: []string{"lifecyclepolicy", "lifecycle-policies", "lifecyclepolicies", "ilm", "ism", "lp", "lifecycle"},
		Args:    cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			policyName := args[0]

			if viper.GetBool("dry-run") {
				cmd.Printf("Would delete lifecycle policy: %s\n", policyName)
				return
			}

			// Check for confirmation flag
			if yes, _ := cmd.Flags().GetBool("yes"); !yes {
				fmt.Printf("Are you sure you want to delete lifecycle policy '%s'? (y/N): ", policyName)
				reader := bufio.NewReader(os.Stdin)
				response, err := reader.ReadString('\n')
				if err != nil {
					fmt.Fprintf(os.Stderr, "Error reading input: %v\n", err)
					os.Exit(1)
				}
				
				response = strings.TrimSpace(strings.ToLower(response))
				if response != "y" && response != "yes" {
					fmt.Println("Operation cancelled")
					return
				}
			}

			c, err := client.NewClient()
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error creating client: %v\n", err)
				os.Exit(1)
			}

			if err := c.DeleteLifecyclePolicy(policyName); err != nil {
				fmt.Fprintf(os.Stderr, "Error deleting lifecycle policy: %v\n", err)
				os.Exit(1)
			}

			cmd.Printf("Lifecycle policy %s deleted successfully\n", policyName)
		},
	}

	cmd.Flags().BoolP("yes", "y", false, "automatically confirm deletion without prompting")
	
	return cmd
}