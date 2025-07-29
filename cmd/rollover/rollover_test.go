package rollover

import (
	"testing"
)

func TestNewRolloverCmd(t *testing.T) {
	cmd := NewRolloverCmd()

	if cmd.Use != "rollover" {
		t.Errorf("Expected Use to be 'rollover', got %s", cmd.Use)
	}

	if cmd.Short != "Rollover a data stream" {
		t.Errorf("Expected Short to be 'Rollover a data stream', got %s", cmd.Short)
	}

	// Check that subcommands are added
	subCommands := cmd.Commands()
	if len(subCommands) == 0 {
		t.Error("Expected subcommands to be added")
	}

	// Check for datastream subcommand
	found := false
	for _, subCmd := range subCommands {
		if subCmd.Use == "datastream DATA_STREAM_NAME" {
			found = true
			break
		}
	}
	if !found {
		t.Error("Expected datastream subcommand to be added")
	}
}

func TestNewRolloverDataStreamCmd(t *testing.T) {
	cmd := NewRolloverDataStreamCmd()

	if cmd.Use != "datastream DATA_STREAM_NAME" {
		t.Errorf("Expected Use to be 'datastream DATA_STREAM_NAME', got %s", cmd.Use)
	}

	if cmd.Short != "Rollover a data stream" {
		t.Errorf("Expected Short to be 'Rollover a data stream', got %s", cmd.Short)
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

	// Check flags
	flags := cmd.Flags()

	expectedFlags := []string{"max-age", "max-docs", "max-size", "max-primary-shard-size", "max-primary-shard-docs", "conditions-file", "lazy"}
	for _, flagName := range expectedFlags {
		if flags.Lookup(flagName) == nil {
			t.Errorf("Expected flag %s to be defined", flagName)
		}
	}
}
