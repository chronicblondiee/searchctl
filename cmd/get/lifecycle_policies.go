package get

import (
	"fmt"
	"os"

	"github.com/chronicblondiee/searchctl/pkg/client"
	"github.com/chronicblondiee/searchctl/pkg/output"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func NewGetLifecyclePoliciesCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "lifecycle-policies [PATTERN]",
		Short:   "Get lifecycle policies",
		Long:    "Get lifecycle policies from the search cluster (ILM for Elasticsearch, ISM for OpenSearch).",
		Aliases: []string{"lifecyclepolicies", "lifecycle-policy", "lifecyclepolicy", "ilm", "ism", "lp", "lifecycle"},
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

			policies, err := c.GetLifecyclePolicies(pattern)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error getting lifecycle policies: %v\n", err)
				os.Exit(1)
			}

			// Convert to interface{} slice for formatting
			data := make([]interface{}, len(policies))
			for i, policy := range policies {
				data[i] = map[string]interface{}{
					"NAME":          policy.Name,
					"VERSION":       policy.Version,
					"MODIFIED_DATE": policy.ModifiedDate,
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