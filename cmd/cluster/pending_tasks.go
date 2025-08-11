package cluster

import (
	"fmt"
	"os"

	"github.com/chronicblondiee/searchctl/pkg/client"
	"github.com/chronicblondiee/searchctl/pkg/output"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func NewPendingTasksCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "pending-tasks",
		Short: "Show cluster pending tasks",
		Long:  "Display cluster pending tasks.",
		Run: func(cmd *cobra.Command, args []string) {
			c, err := client.NewClient()
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error creating client: %v\n", err)
				os.Exit(1)
			}
			pt, err := c.ClusterPendingTasks()
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error getting pending tasks: %v\n", err)
				os.Exit(1)
			}
			data := map[string]interface{}{
				"Tasks": pt.Tasks,
				"Count": len(pt.Tasks),
			}
			formatter := output.NewFormatter(viper.GetString("output"))
			if err := formatter.Format(data, os.Stdout); err != nil {
				fmt.Fprintf(os.Stderr, "Error formatting output: %v\n", err)
				os.Exit(1)
			}
		},
	}
	return cmd
}
