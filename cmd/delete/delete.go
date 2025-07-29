package delete

import (
	"github.com/spf13/cobra"
)

func NewDeleteCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "delete",
		Short: "Delete a resource",
		Long:  "Delete a resource from the search cluster.",
	}

	cmd.AddCommand(NewDeleteIndexCmd())
	cmd.AddCommand(NewDeleteDataStreamCmd())

	return cmd
}
