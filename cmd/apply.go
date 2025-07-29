package cmd

import (
	"fmt"
	"io"
	"os"

	"github.com/chronicblondiee/searchctl/pkg/client"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"gopkg.in/yaml.v3"
)

func NewApplyCmd() *cobra.Command {
	var filename string

	cmd := &cobra.Command{
		Use:   "apply",
		Short: "Apply a configuration from a file",
		Long:  "Apply a configuration to resources by filename or stdin.",
		Run: func(cmd *cobra.Command, args []string) {
			if filename == "" {
				fmt.Fprintf(os.Stderr, "Error: must specify filename with -f flag\n")
				os.Exit(1)
			}

			if viper.GetBool("dry-run") {
				fmt.Printf("Would apply configuration from: %s\n", filename)
				return
			}

			c, err := client.NewClient()
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error creating client: %v\n", err)
				os.Exit(1)
			}

			if err := applyConfigurationFromFile(c, filename); err != nil {
				fmt.Fprintf(os.Stderr, "Error applying configuration: %v\n", err)
				os.Exit(1)
			}

			fmt.Printf("Configuration applied successfully from: %s\n", filename)
		},
	}

	cmd.Flags().StringVarP(&filename, "filename", "f", "", "Filename to apply")
	cmd.MarkFlagRequired("filename")

	return cmd
}

func applyConfigurationFromFile(c client.SearchClient, filename string) error {
	file, err := os.Open(filename)
	if err != nil {
		return fmt.Errorf("failed to open file %s: %w", filename, err)
	}
	defer file.Close()

	data, err := io.ReadAll(file)
	if err != nil {
		return fmt.Errorf("failed to read file %s: %w", filename, err)
	}

	var resource map[string]interface{}
	if err := yaml.Unmarshal(data, &resource); err != nil {
		return fmt.Errorf("failed to parse YAML: %w", err)
	}

	// Determine resource type and apply accordingly
	kind, ok := resource["kind"].(string)
	if !ok {
		return fmt.Errorf("resource kind not specified or invalid")
	}

	switch kind {
	case "IndexTemplate":
		return applyIndexTemplate(c, resource)
	default:
		return fmt.Errorf("unsupported resource kind: %s", kind)
	}
}

func applyIndexTemplate(c client.SearchClient, resource map[string]interface{}) error {
	// Handle both string and interface{} keys in metadata
	var metadata map[string]interface{}
	if meta, ok := resource["metadata"].(map[interface{}]interface{}); ok {
		metadata = make(map[string]interface{})
		for k, v := range meta {
			if key, ok := k.(string); ok {
				metadata[key] = v
			}
		}
	} else if meta, ok := resource["metadata"].(map[string]interface{}); ok {
		metadata = meta
	} else {
		return fmt.Errorf("metadata section missing or invalid")
	}

	name, ok := metadata["name"].(string)
	if !ok {
		return fmt.Errorf("template name missing or invalid")
	}

	// Handle both string and interface{} keys in spec
	var spec map[string]interface{}
	if s, ok := resource["spec"].(map[interface{}]interface{}); ok {
		spec = make(map[string]interface{})
		for k, v := range s {
			if key, ok := k.(string); ok {
				spec[key] = convertInterfaceKeys(v)
			}
		}
	} else if s, ok := resource["spec"].(map[string]interface{}); ok {
		spec = make(map[string]interface{})
		for k, v := range s {
			spec[k] = convertInterfaceKeys(v)
		}
	} else {
		return fmt.Errorf("spec section missing or invalid")
	}

	return c.CreateIndexTemplate(name, spec)
}

// Convert interface{} keys to string keys recursively
func convertInterfaceKeys(v interface{}) interface{} {
	switch val := v.(type) {
	case map[interface{}]interface{}:
		result := make(map[string]interface{})
		for k, v := range val {
			if key, ok := k.(string); ok {
				result[key] = convertInterfaceKeys(v)
			}
		}
		return result
	case []interface{}:
		result := make([]interface{}, len(val))
		for i, item := range val {
			result[i] = convertInterfaceKeys(item)
		}
		return result
	default:
		return val
	}
}
