package output_test

import (
	"bytes"
	"strings"
	"testing"

	"github.com/chronicblondiee/searchctl/pkg/output"
)

func TestTableFormatter(t *testing.T) {
	formatter := &output.TableFormatter{}
	var buf bytes.Buffer

	data := map[string]interface{}{
		"name":   "test-index",
		"health": "green",
		"status": "open",
	}

	err := formatter.Format(data, &buf)
	if err != nil {
		t.Errorf("TableFormatter.Format failed: %v", err)
	}

	result := buf.String()
	if !strings.Contains(result, "name:") || !strings.Contains(result, "test-index") {
		t.Error("Expected output to contain name and value")
	}
}

func TestJSONFormatter(t *testing.T) {
	formatter := &output.JSONFormatter{}
	var buf bytes.Buffer

	data := map[string]interface{}{
		"name":   "test-index",
		"health": "green",
	}

	err := formatter.Format(data, &buf)
	if err != nil {
		t.Errorf("JSONFormatter.Format failed: %v", err)
	}

	result := buf.String()
	if !strings.Contains(result, `"name": "test-index"`) {
		t.Error("Expected JSON output to contain name field")
	}
}

func TestYAMLFormatter(t *testing.T) {
	formatter := &output.YAMLFormatter{}
	var buf bytes.Buffer

	data := map[string]interface{}{
		"name":   "test-index",
		"health": "green",
	}

	err := formatter.Format(data, &buf)
	if err != nil {
		t.Errorf("YAMLFormatter.Format failed: %v", err)
	}

	result := buf.String()
	if !strings.Contains(result, "name: test-index") {
		t.Error("Expected YAML output to contain name field")
	}
}

func TestNewFormatter(t *testing.T) {
	tests := []struct {
		format   string
		expected string
	}{
		{"table", "*output.TableFormatter"},
		{"json", "*output.JSONFormatter"},
		{"yaml", "*output.YAMLFormatter"},
		{"unknown", "*output.TableFormatter"},
	}

	for _, tt := range tests {
		formatter := output.NewFormatter(tt.format)
		if formatter == nil {
			t.Errorf("NewFormatter returned nil for format %s", tt.format)
		}
	}
}
