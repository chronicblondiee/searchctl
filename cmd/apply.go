package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

func NewApplyCmd() *cobra.Command {
	var filename string

	cmd := &cobra.Command{
		Use:   "apply",
		Short: "Apply a configuration from a file",
		Long:  "Apply a configuration to resources by filename or stdin.",
		Run: func(cmd *cobra.Command, args []string) {
			if filename == "" {
				fmt.Println("Error: must specify filename with -f flag")
				return
			}

			fmt.Printf("Would apply configuration from: %s\n", filename)
			// Note: In a real implementation, this would parse and apply the configuration
		},
	}

	cmd.Flags().StringVarP(&filename, "filename", "f", "", "Filename to apply")
	cmd.MarkFlagRequired("filename")

	return cmd
}
