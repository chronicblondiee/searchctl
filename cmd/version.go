package cmd

import (
	"fmt"

	"github.com/chronicblondiee/searchctl/internal/version"
	"github.com/chronicblondiee/searchctl/pkg/output"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func NewVersionCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "version",
		Short: "Print version information",
		Long:  "Print version information for searchctl.",
		Run: func(cmd *cobra.Command, args []string) {
			info := version.Get()

			formatter := output.NewFormatter(viper.GetString("output"))
			if err := formatter.Format(info, cmd.OutOrStdout()); err != nil {
				fmt.Fprintf(cmd.ErrOrStderr(), "Error formatting output: %v\n", err)
			}
		},
	}

	return cmd
}
