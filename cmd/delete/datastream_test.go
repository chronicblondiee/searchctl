package delete

import (
	"testing"
)

func TestNewDeleteDataStreamCmd(t *testing.T) {
	cmd := NewDeleteDataStreamCmd()

	if cmd.Use != "datastream DATA_STREAM_NAME" {
		t.Errorf("Expected Use to be 'datastream DATA_STREAM_NAME', got %s", cmd.Use)
	}

	if cmd.Short != "Delete a data stream" {
		t.Errorf("Expected Short to be 'Delete a data stream', got %s", cmd.Short)
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
}
