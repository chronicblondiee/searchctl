package output

import (
	"encoding/json"
	"fmt"
	"io"
	"text/tabwriter"

	"gopkg.in/yaml.v3"
)

type Formatter interface {
	Format(data interface{}, writer io.Writer) error
}

type TableFormatter struct{}
type JSONFormatter struct{}
type YAMLFormatter struct{}

func NewFormatter(format string) Formatter {
	switch format {
	case "json":
		return &JSONFormatter{}
	case "yaml":
		return &YAMLFormatter{}
	default:
		return &TableFormatter{}
	}
}

func (f *TableFormatter) Format(data interface{}, writer io.Writer) error {
	w := tabwriter.NewWriter(writer, 0, 0, 2, ' ', 0)
	defer w.Flush()

	switch v := data.(type) {
	case []interface{}:
		if len(v) == 0 {
			fmt.Fprintln(w, "No resources found")
			return nil
		}
		return f.formatSlice(v, w)
	default:
		return f.formatSingle(v, w)
	}
}

func (f *TableFormatter) formatSlice(data []interface{}, w *tabwriter.Writer) error {
	if len(data) == 0 {
		return nil
	}

	// Print headers based on first item
	first := data[0]
	switch first.(type) {
	case map[string]interface{}:
		// Extract headers from map keys
		headers := make([]string, 0)
		for k := range first.(map[string]interface{}) {
			headers = append(headers, k)
		}

		// Print headers
		for i, header := range headers {
			if i > 0 {
				fmt.Fprint(w, "\t")
			}
			fmt.Fprint(w, header)
		}
		fmt.Fprintln(w)

		// Print data
		for _, item := range data {
			if m, ok := item.(map[string]interface{}); ok {
				for i, header := range headers {
					if i > 0 {
						fmt.Fprint(w, "\t")
					}
					if val, exists := m[header]; exists {
						fmt.Fprint(w, val)
					}
				}
				fmt.Fprintln(w)
			}
		}
	}

	return nil
}

func (f *TableFormatter) formatSingle(data interface{}, w *tabwriter.Writer) error {
	if m, ok := data.(map[string]interface{}); ok {
		for k, v := range m {
			fmt.Fprintf(w, "%s:\t%v\n", k, v)
		}
	} else {
		// Handle structs by converting to map via JSON
		jsonData, err := json.Marshal(data)
		if err == nil {
			var m map[string]interface{}
			if json.Unmarshal(jsonData, &m) == nil {
				for k, v := range m {
					fmt.Fprintf(w, "%s:\t%v\n", k, v)
				}
			}
		}
	}
	return nil
}

func (f *JSONFormatter) Format(data interface{}, writer io.Writer) error {
	encoder := json.NewEncoder(writer)
	encoder.SetIndent("", "  ")
	return encoder.Encode(data)
}

func (f *YAMLFormatter) Format(data interface{}, writer io.Writer) error {
	encoder := yaml.NewEncoder(writer)
	defer encoder.Close()
	return encoder.Encode(data)
}
