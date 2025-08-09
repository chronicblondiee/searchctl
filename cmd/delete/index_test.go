package delete

import (
	"testing"
)

func TestNewDeleteIndexCmd(t *testing.T) {
	cmd := NewDeleteIndexCmd()

	if cmd.Use != "index INDEX_NAME_OR_PATTERN" {
		t.Errorf("Expected Use to be 'index INDEX_NAME_OR_PATTERN', got %s", cmd.Use)
	}

	if cmd.Short != "Delete an index or indices matching a pattern" {
		t.Errorf("Expected Short to be 'Delete an index or indices matching a pattern', got %s", cmd.Short)
	}

	// Check aliases include short alias
	aliases := cmd.Aliases
	hasIdx := false
	for _, a := range aliases {
		if a == "idx" {
			hasIdx = true
			break
		}
	}
	if !hasIdx {
		t.Errorf("Expected aliases to include 'idx', got %v", aliases)
	}

	// Check args requirement
	if cmd.Args == nil {
		t.Error("Expected Args to be set")
	}

	// Check that -y flag is present
	yFlag := cmd.Flag("yes")
	if yFlag == nil {
		t.Error("Expected -y/--yes flag to be present")
	}
}
