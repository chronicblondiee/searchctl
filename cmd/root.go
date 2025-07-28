package cmd

import (
	"fmt"
	"os"

	"github.com/chronicblondiee/searchctl/cmd/create"
	"github.com/chronicblondiee/searchctl/cmd/delete"
	"github.com/chronicblondiee/searchctl/cmd/describe"
	"github.com/chronicblondiee/searchctl/cmd/get"
	"github.com/chronicblondiee/searchctl/pkg/config"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	cfgFile    string
	context    string
	outputFlag string
	verbose    bool
)

var rootCmd = &cobra.Command{
	Use:   "searchctl",
	Short: "A kubectl-like CLI for OpenSearch and Elasticsearch management",
	Long: `searchctl is a command-line interface for managing OpenSearch and Elasticsearch clusters.
It provides familiar kubectl-like commands for cluster administration, index management, and more.`,
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		if err := config.InitConfig(cfgFile); err != nil {
			fmt.Fprintf(os.Stderr, "Error initializing config: %v\n", err)
			os.Exit(1)
		}
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}

func NewRootCmd() *cobra.Command {
	return rootCmd
}

func init() {
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.searchctl/config.yaml)")
	rootCmd.PersistentFlags().StringVar(&context, "context", "", "override current context")
	rootCmd.PersistentFlags().StringVarP(&outputFlag, "output", "o", "table", "output format (table|json|yaml|wide)")
	rootCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "verbose output")
	rootCmd.PersistentFlags().Bool("dry-run", false, "show what would be done without executing")

	viper.BindPFlag("context", rootCmd.PersistentFlags().Lookup("context"))
	viper.BindPFlag("output", rootCmd.PersistentFlags().Lookup("output"))
	viper.BindPFlag("verbose", rootCmd.PersistentFlags().Lookup("verbose"))
	viper.BindPFlag("dry-run", rootCmd.PersistentFlags().Lookup("dry-run"))

	// Add subcommands
	rootCmd.AddCommand(get.NewGetCmd())
	rootCmd.AddCommand(describe.NewDescribeCmd())
	rootCmd.AddCommand(create.NewCreateCmd())
	rootCmd.AddCommand(delete.NewDeleteCmd())
	rootCmd.AddCommand(NewApplyCmd())
	rootCmd.AddCommand(NewConfigCmd())
	rootCmd.AddCommand(NewClusterCmd())
	rootCmd.AddCommand(NewVersionCmd())
}
