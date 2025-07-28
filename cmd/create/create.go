package create

import (
	"github.com/spf13/cobra"
)

func NewCreateCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "create",
		Short: "Create a resource",
		Long:  "Create a new resource in the search cluster.",
	}

	cmd.AddCommand(NewCreateIndexCmd())

	return cmd
}
