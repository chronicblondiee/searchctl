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

	// Check aliases
	aliases := cmd.Aliases
	if len(aliases) != 1 || aliases[0] != "idx" {
		t.Errorf("Expected aliases to be ['idx'], got %v", aliases)
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
