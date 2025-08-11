package clone

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/chronicblondiee/searchctl/pkg/client"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"gopkg.in/yaml.v3"
)

type exportOptions struct {
	dir           string
	types         []string
	names         []string
	includeSystem bool
	all           bool
}

func NewExportCmd() *cobra.Command {
	var opts exportOptions
	cmd := &cobra.Command{
		Use:   "export",
		Short: "Export cluster configuration to a directory",
		RunE: func(cmd *cobra.Command, args []string) error {
			if opts.dir == "" {
				return fmt.Errorf("must provide --dir output directory")
			}
			return runExport(opts)
		},
	}
	cmd.Flags().StringVarP(&opts.dir, "dir", "d", "", "output directory")
	cmd.Flags().StringSliceVar(&opts.types, "types", []string{}, "resource types to export (index-templates,component-templates,lifecycle-policies,ingest-pipelines,cluster-settings)")
	cmd.Flags().StringSliceVar(&opts.names, "names", []string{}, "optional names/patterns to filter (comma-separated)")
	cmd.Flags().BoolVar(&opts.includeSystem, "include-system", false, "include system resources (e.g. names starting with .)")
	cmd.Flags().BoolVar(&opts.all, "all", false, "export all supported resource types")
	return cmd
}

func writeDoc(path string, doc interface{}) error {
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return err
	}
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()
	switch viper.GetString("output") {
	case "json":
		enc := json.NewEncoder(f)
		enc.SetIndent("", "  ")
		return enc.Encode(doc)
	default:
		enc := yaml.NewEncoder(f)
		defer enc.Close()
		return enc.Encode(doc)
	}
}

func runExport(opts exportOptions) error {
	c, err := client.NewClient()
	if err != nil {
		return fmt.Errorf("error creating client: %w", err)
	}

	// Manifest
	manifest := map[string]interface{}{
		"kind":     "CloneManifest",
		"metadata": map[string]interface{}{},
		"spec":     map[string]interface{}{},
	}
	if err := writeDoc(filepath.Join(opts.dir, "manifest.yaml"), manifest); err != nil {
		return err
	}

	selected := map[string]bool{}
	if opts.all {
		opts.types = append(opts.types, "all")
	}
	if len(opts.types) == 0 || contains(opts.types, "index-templates") || contains(opts.types, "all") {
		selected["index-templates"] = true
	}
	if len(opts.types) == 0 || contains(opts.types, "component-templates") || contains(opts.types, "all") {
		selected["component-templates"] = true
	}
	if len(opts.types) == 0 || contains(opts.types, "lifecycle-policies") || contains(opts.types, "ilm") || contains(opts.types, "all") {
		selected["lifecycle-policies"] = true
	}
	if len(opts.types) == 0 || contains(opts.types, "ingest-pipelines") || contains(opts.types, "all") {
		selected["ingest-pipelines"] = true
	}
	if len(opts.types) == 0 || contains(opts.types, "cluster-settings") || contains(opts.types, "all") {
		selected["cluster-settings"] = true
	}

	patterns := opts.names
	if len(patterns) == 0 {
		patterns = []string{""}
	}

	// component-templates first
	if selected["component-templates"] {
		for _, p := range patterns {
			items, err := c.GetComponentTemplates(p)
			if err != nil {
				// Treat missing endpoint or 404 payloads as no-op for this type
				if strings.Contains(err.Error(), "404") || strings.Contains(err.Error(), "not found") {
					continue
				}
				return err
			}
			for _, ct := range items {
				if !opts.includeSystem && strings.HasPrefix(ct.Name, ".") {
					continue
				}
				doc := map[string]interface{}{
					"kind": "ComponentTemplate",
					"metadata": map[string]interface{}{
						"name": ct.Name,
					},
					"spec": map[string]interface{}{
						"template": ct.Template,
					},
				}
				if ct.Version != 0 {
					doc["spec"].(map[string]interface{})["version"] = ct.Version
				}
				if len(ct.Meta) > 0 {
					doc["spec"].(map[string]interface{})["_meta"] = ct.Meta
				}
				path := filepath.Join(opts.dir, "component-templates", safeName(ct.Name)+ext())
				if err := writeDoc(path, doc); err != nil {
					return err
				}
			}
		}
	}

	if selected["index-templates"] {
		for _, p := range patterns {
			items, err := c.GetIndexTemplates(p)
			if err != nil {
				if strings.Contains(err.Error(), "404") || strings.Contains(err.Error(), "not found") {
					continue
				}
				return err
			}
			for _, it := range items {
				if !opts.includeSystem && strings.HasPrefix(it.Name, ".") {
					continue
				}
				spec := map[string]interface{}{
					"index_patterns": it.IndexPattern,
				}
				if len(it.Template.Settings) > 0 || len(it.Template.Mappings) > 0 || len(it.Template.Aliases) > 0 {
					spec["template"] = it.Template
				}
				if len(it.ComposedOf) > 0 {
					spec["composed_of"] = it.ComposedOf
				}
				if it.Priority != 0 {
					spec["priority"] = it.Priority
				}
				if it.Version != 0 {
					spec["version"] = it.Version
				}
				if len(it.Meta) > 0 {
					spec["_meta"] = it.Meta
				}
				if len(it.DataStream) > 0 {
					spec["data_stream"] = it.DataStream
				}
				doc := map[string]interface{}{
					"kind": "IndexTemplate",
					"metadata": map[string]interface{}{
						"name": it.Name,
					},
					"spec": spec,
				}
				path := filepath.Join(opts.dir, "index-templates", safeName(it.Name)+ext())
				if err := writeDoc(path, doc); err != nil {
					return err
				}
			}
		}
	}

	if selected["lifecycle-policies"] {
		for _, p := range patterns {
			items, err := c.GetLifecyclePolicies(p)
			if err != nil {
				if strings.Contains(err.Error(), "no handler found") || strings.Contains(err.Error(), "404") || strings.Contains(strings.ToLower(err.Error()), "not found") {
					continue
				}
				return err
			}
			for _, lp := range items {
				if !opts.includeSystem && strings.HasPrefix(lp.Name, ".") {
					continue
				}
				doc := map[string]interface{}{
					"kind": "LifecyclePolicy",
					"metadata": map[string]interface{}{
						"name": lp.Name,
					},
					"spec": lp.Policy,
				}
				path := filepath.Join(opts.dir, "lifecycle-policies", safeName(lp.Name)+ext())
				if err := writeDoc(path, doc); err != nil {
					return err
				}
			}
		}
	}

	if selected["ingest-pipelines"] {
		for _, p := range patterns {
			items, err := c.GetIngestPipelines(p)
			if err != nil {
				return err
			}
			for _, pl := range items {
				if !opts.includeSystem && strings.HasPrefix(pl.Name, ".") {
					continue
				}
				doc := map[string]interface{}{
					"kind": "IngestPipeline",
					"metadata": map[string]interface{}{
						"name": pl.Name,
					},
					"spec": pl.Body,
				}
				path := filepath.Join(opts.dir, "ingest-pipelines", safeName(pl.Name)+ext())
				if err := writeDoc(path, doc); err != nil {
					return err
				}
			}
		}
	}

	if selected["cluster-settings"] {
		settings, err := c.GetClusterSettings()
		if err != nil {
			return err
		}
		doc := map[string]interface{}{
			"kind":     "ClusterSettings",
			"metadata": map[string]interface{}{},
			"spec": map[string]interface{}{
				"persistent": settings.Persistent,
				"transient":  settings.Transient,
			},
		}
		if err := writeDoc(filepath.Join(opts.dir, "cluster-settings", "cluster-settings"+ext()), doc); err != nil {
			return err
		}
	}

	fmt.Printf("Exported resources to %s\n", opts.dir)
	return nil
}

func contains(slice []string, v string) bool {
	for _, s := range slice {
		if s == v {
			return true
		}
	}
	return false
}

func ext() string {
	if viper.GetString("output") == "json" {
		return ".json"
	}
	return ".yaml"
}

func safeName(s string) string {
	s = strings.ReplaceAll(s, string(os.PathSeparator), "_")
	return s
}
