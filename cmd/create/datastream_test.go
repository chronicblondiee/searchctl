package create

import (
	"testing"
)

func TestNewCreateDataStreamCmd(t *testing.T) {
	cmd := NewCreateDataStreamCmd()

	if cmd.Use != "datastream DATA_STREAM_NAME" {
		t.Errorf("Expected Use to be 'datastream DATA_STREAM_NAME', got %s", cmd.Use)
	}

	if cmd.Short != "Create a data stream" {
		t.Errorf("Expected Short to be 'Create a data stream', got %s", cmd.Short)
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
