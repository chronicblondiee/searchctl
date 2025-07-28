package get

import (
	"bytes"
	"strings"
	"testing"

	"github.com/spf13/cobra"
)

func TestGetCommand(t *testing.T) {
	cmd := &cobra.Command{Use: "test"}
	cmd.AddCommand(NewGetCmd())

	// Test help output
	var buf bytes.Buffer
	cmd.SetOut(&buf)
	cmd.SetArgs([]string{"get", "--help"})
	
	err := cmd.Execute()
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	output := buf.String()
	if !bytes.Contains([]byte(output), []byte("display one or many")) {
		t.Error("Expected help text to contain 'display one or many'")
	}
}

func TestGetCommandWithoutSubcommand(t *testing.T) {
	cmd := &cobra.Command{Use: "test"}
	cmd.AddCommand(NewGetCmd())

	var buf bytes.Buffer
	cmd.SetOut(&buf)
	cmd.SetErr(&buf)
	cmd.SetArgs([]string{"get"})
	
	err := cmd.Execute()
	// Should show help when no subcommand is provided
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
	
	output := buf.String()
	if !bytes.Contains([]byte(output), []byte("Available Commands")) {
		t.Error("Expected help output when no subcommand provided")
	}
}

func TestGetCommandValidSubcommands(t *testing.T) {
	cmd := NewGetCmd()
	
	// Check that the command has the expected subcommands
	expectedSubcommands := map[string]bool{
		"indices": false,
		"nodes":   false,
	}
	
	for _, child := range cmd.Commands() {
		cmdName := child.Use
		// Extract just the command name (before any spaces)
		if spaceIdx := strings.Index(cmdName, " "); spaceIdx >= 0 {
			cmdName = cmdName[:spaceIdx]
		}
		if _, exists := expectedSubcommands[cmdName]; exists {
			expectedSubcommands[cmdName] = true
		}
	}
	
	for subcmd, found := range expectedSubcommands {
		if !found {
			t.Errorf("Expected subcommand '%s' not found", subcmd)
		}
	}
}
