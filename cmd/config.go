package cmd

import (
	"fmt"
	"os"

	"github.com/chronicblondiee/searchctl/pkg/config"
	"github.com/chronicblondiee/searchctl/pkg/output"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func NewConfigCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "config",
		Short: "Modify searchctl configuration",
		Long:  "Display and modify searchctl configuration settings.",
	}

	cmd.AddCommand(NewConfigViewCmd())
	cmd.AddCommand(NewConfigUseContextCmd())

	return cmd
}

func NewConfigViewCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "view",
		Short: "Display current configuration",
		Long:  "Display the current searchctl configuration.",
		Run: func(cmd *cobra.Command, args []string) {
			cfg := config.GetConfig()
			if cfg == nil {
				fmt.Fprintf(os.Stderr, "No configuration found\n")
				os.Exit(1)
			}

			formatter := output.NewFormatter(viper.GetString("output"))
			if err := formatter.Format(cfg, os.Stdout); err != nil {
				fmt.Fprintf(os.Stderr, "Error formatting output: %v\n", err)
				os.Exit(1)
			}
		},
	}

	return cmd
}

func NewConfigUseContextCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "use-context CONTEXT_NAME",
		Short: "Set the current context",
		Long:  "Set the current context for searchctl operations.",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			contextName := args[0]
			fmt.Printf("Switched to context %q\n", contextName)
			// Note: In a real implementation, this would update the config file
		},
	}

	return cmd
}
