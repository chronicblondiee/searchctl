package describe

import (
	"github.com/spf13/cobra"
)

func NewDescribeCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "describe",
		Short: "Show details of a specific resource",
		Long:  "Show detailed information about a specific resource.",
	}

	cmd.AddCommand(NewDescribeIndexCmd())

	return cmd
}
