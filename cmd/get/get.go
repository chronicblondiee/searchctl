package get

import (
	"github.com/spf13/cobra"
)

func NewGetCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "get",
		Short: "Display one or many resources",
		Long:  "Get and display one or many resources from the search cluster.",
	}

	cmd.AddCommand(NewGetIndicesCmd())
	cmd.AddCommand(NewGetNodesCmd())
	cmd.AddCommand(NewGetDataStreamsCmd())
	cmd.AddCommand(NewGetIndexTemplatesCmd())

	return cmd
}
