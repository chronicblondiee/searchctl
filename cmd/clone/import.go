package clone

import (
	"encoding/json"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/chronicblondiee/searchctl/pkg/client"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
)

type importOptions struct {
	dir             string
	types           []string
	continueOnError bool
	dryRun          bool
}

func NewImportCmd() *cobra.Command {
	var opts importOptions
	cmd := &cobra.Command{
		Use:   "import",
		Short: "Import cluster configuration from a directory",
		RunE: func(cmd *cobra.Command, args []string) error {
			if opts.dir == "" {
				return fmt.Errorf("must provide --dir input directory")
			}
			return runImport(opts)
		},
	}
	cmd.Flags().StringVarP(&opts.dir, "dir", "d", "", "input directory")
	cmd.Flags().StringSliceVar(&opts.types, "types", []string{}, "resource types to import (component-templates,index-templates,lifecycle-policies,ingest-pipelines,cluster-settings)")
	cmd.Flags().BoolVar(&opts.continueOnError, "continue-on-error", false, "continue when a file fails")
	cmd.Flags().BoolVar(&opts.dryRun, "dry-run", false, "show planned operations without applying")
	return cmd
}

func runImport(opts importOptions) error {
	// Import order: component-templates -> index-templates -> lifecycle-policies -> ingest-pipelines -> cluster-settings
	order := []string{"component-templates", "index-templates", "lifecycle-policies", "ingest-pipelines", "cluster-settings"}
	selected := map[string]bool{}
	if len(opts.types) == 0 {
		for _, t := range order {
			selected[t] = true
		}
	} else {
		for _, t := range opts.types {
			selected[t] = true
		}
	}

	c, err := client.NewClient()
	if err != nil {
		return fmt.Errorf("error creating client: %w", err)
	}

	for _, t := range order {
		if !selected[t] {
			continue
		}
		dir := filepath.Join(opts.dir, t)
		entries := []string{}
		_ = filepath.WalkDir(dir, func(path string, d fs.DirEntry, err error) error {
			if err != nil || d == nil || d.IsDir() {
				return nil
			}
			if strings.HasSuffix(d.Name(), ".yaml") || strings.HasSuffix(d.Name(), ".yml") || strings.HasSuffix(d.Name(), ".json") {
				entries = append(entries, path)
			}
			return nil
		})
		sort.Strings(entries)
		for _, f := range entries {
			data, err := os.ReadFile(f)
			if err != nil {
				if opts.continueOnError {
					fmt.Fprintf(os.Stderr, "[WARN] read %s: %v\n", f, err)
					continue
				}
				return err
			}
			var obj map[string]interface{}
			if strings.HasSuffix(f, ".json") {
				if err := json.Unmarshal(data, &obj); err != nil {
					if opts.continueOnError {
						fmt.Fprintf(os.Stderr, "[WARN] parse %s: %v\n", f, err)
						continue
					}
					return err
				}
			} else {
				if err := yaml.Unmarshal(data, &obj); err != nil {
					if opts.continueOnError {
						fmt.Fprintf(os.Stderr, "[WARN] parse %s: %v\n", f, err)
						continue
					}
					return err
				}
			}
			kind := inferKind(t, obj)
			name := extractName(obj)
			spec := extractSpec(obj)
			if opts.dryRun {
				fmt.Printf("Would apply %s/%s from %s\n", kind, name, f)
				continue
			}
			if err := applyOne(c, kind, name, spec); err != nil {
				if opts.continueOnError {
					fmt.Fprintf(os.Stderr, "[WARN] apply %s: %v\n", f, err)
					continue
				}
				return fmt.Errorf("%s: %w", f, err)
			}
			fmt.Printf("Applied %s/%s\n", kind, name)
		}
	}
	return nil
}

func inferKind(folder string, obj map[string]interface{}) string {
	if k, ok := obj["kind"].(string); ok && k != "" {
		return k
	}
	switch folder {
	case "component-templates":
		return "ComponentTemplate"
	case "index-templates":
		return "IndexTemplate"
	case "lifecycle-policies":
		return "LifecyclePolicy"
	case "ingest-pipelines":
		return "IngestPipeline"
	case "cluster-settings":
		return "ClusterSettings"
	default:
		return ""
	}
}

func extractName(obj map[string]interface{}) string {
	if m, ok := obj["metadata"].(map[string]interface{}); ok {
		if n, ok2 := m["name"].(string); ok2 {
			return n
		}
	}
	return ""
}

func extractSpec(obj map[string]interface{}) map[string]interface{} {
	if s, ok := obj["spec"].(map[string]interface{}); ok {
		return s
	}
	return obj
}

func applyOne(c client.SearchClient, kind, name string, spec map[string]interface{}) error {
	switch kind {
	case "ComponentTemplate":
		return c.CreateComponentTemplate(name, spec)
	case "IndexTemplate":
		return c.CreateIndexTemplate(name, spec)
	case "LifecyclePolicy":
		return c.CreateLifecyclePolicy(name, spec)
	case "IngestPipeline":
		return c.CreateIngestPipeline(name, spec)
	case "ClusterSettings":
		return c.UpdateClusterSettings(spec)
	default:
		return fmt.Errorf("unsupported kind: %s", kind)
	}
}
