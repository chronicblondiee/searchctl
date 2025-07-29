package create

import (
	"fmt"
	"os"

	"github.com/chronicblondiee/searchctl/pkg/client"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func NewCreateDataStreamCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "datastream DATA_STREAM_NAME",
		Short:   "Create a data stream",
		Long:    "Create a new data stream in the search cluster. Note: A matching index template with data_stream configuration must exist before creating the data stream.",
		Aliases: []string{"ds"},
		Args:    cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			dataStreamName := args[0]

			if viper.GetBool("dry-run") {
				cmd.Printf("Would create data stream: %s\n", dataStreamName)
				return
			}

			c, err := client.NewClient()
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error creating client: %v\n", err)
				os.Exit(1)
			}

			if err := c.CreateDataStream(dataStreamName); err != nil {
				fmt.Fprintf(os.Stderr, "Error creating data stream: %v\n", err)
				os.Exit(1)
			}

			cmd.Printf("Data stream %s created successfully\n", dataStreamName)
		},
	}

	return cmd
}
