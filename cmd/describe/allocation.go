package describe

import (
	"fmt"
	"os"

	"github.com/chronicblondiee/searchctl/pkg/client"
	"github.com/chronicblondiee/searchctl/pkg/output"
	"github.com/chronicblondiee/searchctl/pkg/types"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func NewDescribeAllocationCmd() *cobra.Command {
	var index string
	var shard int
	var primary bool
	var includeYes bool
	var includeDisk bool

	cmd := &cobra.Command{
		Use:   "allocation",
		Short: "Explain shard allocation decisions",
		Long:  "Explain shard allocation decisions for a given shard using the cluster allocation explain API.",
		Run: func(cmd *cobra.Command, args []string) {
			c, err := client.NewClient()
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error creating client: %v\n", err)
				os.Exit(1)
			}

			req := types.AllocationExplainRequest{Index: index, Shard: shard, Primary: primary}
			resp, err := c.ExplainAllocation(req, includeYes, includeDisk)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error explaining allocation: %v\n", err)
				os.Exit(1)
			}

			outFmt := viper.GetString("output")
			formatter := output.NewFormatter(outFmt)
			if err := formatter.Format(resp, os.Stdout); err != nil {
				fmt.Fprintf(os.Stderr, "Error formatting output: %v\n", err)
				os.Exit(1)
			}
		},
	}

	cmd.Flags().StringVar(&index, "index", "", "index name")
	cmd.Flags().IntVar(&shard, "shard", 0, "shard number")
	cmd.Flags().BoolVar(&primary, "primary", false, "explain primary shard (default replica)")
	cmd.Flags().BoolVar(&includeYes, "include-yes", false, "include yes decisions")
	cmd.Flags().BoolVar(&includeDisk, "include-disk", false, "include disk info")

	return cmd
}
