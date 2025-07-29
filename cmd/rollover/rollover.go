package rollover

import (
	"github.com/spf13/cobra"
)

func NewRolloverCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "rollover",
		Short: "Rollover a data stream",
		Long:  "Rollover a data stream in the search cluster.",
	}

	cmd.AddCommand(NewRolloverDataStreamCmd())

	return cmd
}
