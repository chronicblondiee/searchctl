package delete

import (
	"fmt"
	"os"

	"github.com/chronicblondiee/searchctl/pkg/client"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func NewDeleteDataStreamCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "datastream DATA_STREAM_NAME",
		Short:   "Delete a data stream",
		Long:    "Delete a data stream and all its backing indices from the search cluster.",
		Aliases: []string{"ds"},
		Args:    cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			dataStreamName := args[0]

			if viper.GetBool("dry-run") {
				cmd.Printf("Would delete data stream: %s\n", dataStreamName)
				return
			}

			c, err := client.NewClient()
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error creating client: %v\n", err)
				os.Exit(1)
			}

			if err := c.DeleteDataStream(dataStreamName); err != nil {
				fmt.Fprintf(os.Stderr, "Error deleting data stream: %v\n", err)
				os.Exit(1)
			}

			cmd.Printf("Data stream %s deleted successfully\n", dataStreamName)
		},
	}

	return cmd
}
