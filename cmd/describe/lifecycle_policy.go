package describe

import (
	"fmt"
	"os"

	"github.com/chronicblondiee/searchctl/pkg/client"
	"github.com/chronicblondiee/searchctl/pkg/output"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func NewDescribeLifecyclePolicyCmd() *cobra.Command {
	var showBody bool

	cmd := &cobra.Command{
		Use:     "lifecycle-policy NAME",
		Short:   "Describe a lifecycle policy",
		Long:    "Show detailed information about a specific lifecycle policy (ILM or ISM).",
		Aliases: []string{"lifecyclepolicy", "lifecycle-policies", "lifecyclepolicies", "ilm", "ism", "lp", "lifecycle"},
		Args:    cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			name := args[0]

			c, err := client.NewClient()
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error creating client: %v\n", err)
				os.Exit(1)
			}

			policy, err := c.GetLifecyclePolicy(name)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error getting lifecycle policy: %v\n", err)
				os.Exit(1)
			}

			outFmt := viper.GetString("output")
			if outFmt == "json" || outFmt == "yaml" {
				formatter := output.NewFormatter(outFmt)
				if err := formatter.Format(policy, os.Stdout); err != nil {
					fmt.Fprintf(os.Stderr, "Error formatting output: %v\n", err)
					os.Exit(1)
				}
				return
			}

			data := map[string]interface{}{
				"Name":         policy.Name,
				"Version":      policy.Version,
				"ModifiedDate": policy.ModifiedDate,
			}

			if showBody {
				data["Policy"] = policy.Policy
			}

			formatter := output.NewFormatter(outFmt)
			if err := formatter.Format(data, os.Stdout); err != nil {
				fmt.Fprintf(os.Stderr, "Error formatting output: %v\n", err)
				os.Exit(1)
			}
		},
	}

	cmd.Flags().BoolVar(&showBody, "show-body", false, "include full policy body in table output")

	return cmd
}
