package delete

import (
	"testing"
)

func TestNewDeleteDataStreamCmd(t *testing.T) {
	cmd := NewDeleteDataStreamCmd()

	if cmd.Use != "datastream DATA_STREAM_NAME_OR_PATTERN" {
		t.Errorf("Expected Use to be 'datastream DATA_STREAM_NAME_OR_PATTERN', got %s", cmd.Use)
	}

	if cmd.Short != "Delete a data stream or data streams matching a pattern" {
		t.Errorf("Expected Short to be 'Delete a data stream or data streams matching a pattern', got %s", cmd.Short)
	}

	// Check aliases
	aliases := cmd.Aliases
	if len(aliases) != 1 || aliases[0] != "ds" {
		t.Errorf("Expected aliases to be ['ds'], got %v", aliases)
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
