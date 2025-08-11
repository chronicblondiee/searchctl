package clone

import "github.com/spf13/cobra"

func NewCloneCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "clone",
		Short: "Export (clone) and import cluster configuration resources",
		Long:  "Clone cluster configuration (templates, policies, pipelines, settings) to/from the filesystem.",
	}
	cmd.AddCommand(NewExportCmd())
	cmd.AddCommand(NewImportCmd())
	return cmd
}
