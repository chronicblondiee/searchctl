package rollover

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/chronicblondiee/searchctl/pkg/client"
	"github.com/chronicblondiee/searchctl/pkg/output"
	"github.com/chronicblondiee/searchctl/pkg/types"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	yaml "gopkg.in/yaml.v3"
)

func NewRolloverDataStreamCmd() *cobra.Command {
	var maxAge string
	var maxDocs int64
	var maxSize string
	var maxPrimaryShardSize string
	var maxPrimaryShardDocs int64
	var conditionsFile string
	var lazy bool

	cmd := &cobra.Command{
		Use:     "datastream DATA_STREAM_NAME",
		Short:   "Rollover a data stream",
		Long:    "Rollover a data stream to create a new backing index.",
		Aliases: []string{"ds"},
		Args:    cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			dataStreamName := args[0]

			// Build conditions
			conditions := make(map[string]interface{})

			if maxAge != "" {
				conditions["max_age"] = maxAge
			}
			if maxDocs > 0 {
				conditions["max_docs"] = maxDocs
			}
			if maxSize != "" {
				conditions["max_size"] = maxSize
			}
			if maxPrimaryShardSize != "" {
				conditions["max_primary_shard_size"] = maxPrimaryShardSize
			}
			if maxPrimaryShardDocs > 0 {
				conditions["max_primary_shard_docs"] = maxPrimaryShardDocs
			}

			// If conditions file is provided, read from it
			if conditionsFile != "" {
				fileConditions, err := readConditionsFromFile(conditionsFile)
				if err != nil {
					fmt.Fprintf(os.Stderr, "Error reading conditions file: %v\n", err)
					os.Exit(1)
				}
				for k, v := range fileConditions {
					conditions[k] = v
				}
			}

			if viper.GetBool("dry-run") {
				conditionsJSON, _ := json.MarshalIndent(conditions, "", "  ")
				cmd.Printf("Would rollover data stream: %s\nConditions:\n%s\n", dataStreamName, string(conditionsJSON))
				if lazy {
					cmd.Printf("Lazy rollover: true\n")
				}
				return
			}

			c, err := client.NewClient()
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error creating client: %v\n", err)
				os.Exit(1)
			}

			response, err := c.RolloverDataStream(dataStreamName, conditions, lazy)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error rolling over data stream: %v\n", err)
				os.Exit(1)
			}

			// Format and display response
			outputFormat := viper.GetString("output")
			if outputFormat == "json" || outputFormat == "yaml" {
				formatter := output.NewFormatter(outputFormat)
				if err := formatter.Format([]interface{}{response}, os.Stdout); err != nil {
					fmt.Fprintf(os.Stderr, "Error formatting output: %v\n", err)
					os.Exit(1)
				}
			} else {
				displayRolloverResult(response)
			}
		},
	}

	cmd.Flags().StringVar(&maxAge, "max-age", "", "Maximum age before rollover (e.g., 30d, 1h)")
	cmd.Flags().Int64Var(&maxDocs, "max-docs", 0, "Maximum number of documents before rollover")
	cmd.Flags().StringVar(&maxSize, "max-size", "", "Maximum index size before rollover (e.g., 50gb, 5gb)")
	cmd.Flags().StringVar(&maxPrimaryShardSize, "max-primary-shard-size", "", "Maximum primary shard size before rollover (e.g., 50gb)")
	cmd.Flags().Int64Var(&maxPrimaryShardDocs, "max-primary-shard-docs", 0, "Maximum number of documents in primary shard before rollover")
	cmd.Flags().StringVarP(&conditionsFile, "conditions-file", "f", "", "file containing rollover conditions (JSON or YAML)")
	cmd.Flags().BoolVar(&lazy, "lazy", false, "Only mark data stream for rollover at next write (data streams only)")

	return cmd
}

func readConditionsFromFile(filename string) (map[string]interface{}, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	// Try by extension first
	ext := filepath.Ext(filename)
	switch ext {
	case ".yaml", ".yml":
		// YAML parse with key conversion
		var raw interface{}
		if err := yaml.Unmarshal(data, &raw); err != nil {
			return nil, fmt.Errorf("failed to parse YAML: %w", err)
		}
		converted, ok := convertInterfaceKeys(raw).(map[string]interface{})
		if !ok {
			return nil, fmt.Errorf("invalid YAML structure for conditions; expected mapping at top level")
		}
		return converted, nil
	case ".json":
		var conditions map[string]interface{}
		if err := json.Unmarshal(data, &conditions); err != nil {
			return nil, fmt.Errorf("failed to parse JSON: %w", err)
		}
		return conditions, nil
	default:
		// Try JSON first, then YAML as fallback
		var conditions map[string]interface{}
		if err := json.Unmarshal(data, &conditions); err == nil {
			return conditions, nil
		}
		var raw interface{}
		if err := yaml.Unmarshal(data, &raw); err == nil {
			converted, ok := convertInterfaceKeys(raw).(map[string]interface{})
			if !ok {
				return nil, fmt.Errorf("invalid YAML structure for conditions; expected mapping at top level")
			}
			return converted, nil
		}
		return nil, fmt.Errorf("failed to parse conditions file as JSON or YAML")
	}
}

// convertInterfaceKeys converts map[interface{}]interface{} to map[string]interface{} recursively
func convertInterfaceKeys(value interface{}) interface{} {
	switch typed := value.(type) {
	case map[interface{}]interface{}:
		result := make(map[string]interface{}, len(typed))
		for k, v := range typed {
			key, ok := k.(string)
			if !ok {
				// skip non-string keys
				continue
			}
			result[key] = convertInterfaceKeys(v)
		}
		return result
	case map[string]interface{}:
		result := make(map[string]interface{}, len(typed))
		for k, v := range typed {
			result[k] = convertInterfaceKeys(v)
		}
		return result
	case []interface{}:
		for i, item := range typed {
			typed[i] = convertInterfaceKeys(item)
		}
		return typed
	default:
		return typed
	}
}

func displayRolloverResult(response *types.RolloverResponse) {
	fmt.Printf("Rollover Status: %s\n", getStatusMessage(response.RolledOver))
	if response.RolledOver {
		fmt.Printf("Old Index: %s\n", response.OldIndex)
		fmt.Printf("New Index: %s\n", response.NewIndex)
	}

	if response.DryRun {
		fmt.Printf("Dry Run: true\n")
	}

	if len(response.Conditions) > 0 {
		fmt.Printf("\nCondition Results:\n")
		for condition, met := range response.Conditions {
			status := "✗ Not met"
			if met {
				status = "✓ Met"
			}
			fmt.Printf("  %s: %s\n", condition, status)
		}
	}

	fmt.Printf("Acknowledged: %s\n", formatBool(response.Acknowledged))
	if response.ShardsAcknowledged {
		fmt.Printf("Shards Acknowledged: %s\n", formatBool(response.ShardsAcknowledged))
	}
}

func getStatusMessage(rolledOver bool) string {
	if rolledOver {
		return "SUCCESS"
	}
	return "NO ROLLOVER NEEDED"
}

func formatBool(b bool) string {
	if b {
		return "true"
	}
	return "false"
}
