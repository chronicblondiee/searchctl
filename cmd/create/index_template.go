package create

import (
	"fmt"
	"io"
	"os"

	"github.com/chronicblondiee/searchctl/pkg/client"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"gopkg.in/yaml.v3"
)

func NewCreateIndexTemplateCmd() *cobra.Command {
	var filename string

	cmd := &cobra.Command{
		Use:     "index-template TEMPLATE_NAME",
		Short:   "Create an index template",
		Long:    "Create a new index template in the search cluster.",
		Aliases: []string{"template", "it"},
		Args:    cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			templateName := args[0]

			if viper.GetBool("dry-run") {
				cmd.Printf("Would create index template: %s\n", templateName)
				return
			}

			c, err := client.NewClient()
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error creating client: %v\n", err)
				os.Exit(1)
			}

			var templateBody map[string]interface{}

			if filename != "" {
				// Read template from file
				templateBody, err = readTemplateFromFile(filename)
				if err != nil {
					fmt.Fprintf(os.Stderr, "Error reading template file: %v\n", err)
					os.Exit(1)
				}
			} else {
				// Use default template
				templateBody = map[string]interface{}{
					"index_patterns": []string{templateName + "-*"},
					"template": map[string]interface{}{
						"settings": map[string]interface{}{
							"number_of_shards":   1,
							"number_of_replicas": 0,
						},
						"mappings": map[string]interface{}{
							"properties": map[string]interface{}{
								"@timestamp": map[string]interface{}{
									"type": "date",
								},
								"message": map[string]interface{}{
									"type": "text",
								},
							},
						},
					},
				}
			}

			if err := c.CreateIndexTemplate(templateName, templateBody); err != nil {
				fmt.Fprintf(os.Stderr, "Error creating index template: %v\n", err)
				os.Exit(1)
			}

			cmd.Printf("Index template %s created successfully\n", templateName)
		},
	}

	cmd.Flags().StringVarP(&filename, "filename", "f", "", "Template definition file (YAML)")

	return cmd
}

func readTemplateFromFile(filename string) (map[string]interface{}, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	data, err := io.ReadAll(file)
	if err != nil {
		return nil, err
	}

	var template map[string]interface{}
	if err := yaml.Unmarshal(data, &template); err != nil {
		return nil, err
	}

	// Extract the spec if it's a Kubernetes-style resource
	if spec, ok := template["spec"].(map[interface{}]interface{}); ok {
		// Convert interface{} keys to string keys
		result := make(map[string]interface{})
		for k, v := range spec {
			if key, ok := k.(string); ok {
				result[key] = v
			}
		}
		return result, nil
	}

	return template, nil
}
