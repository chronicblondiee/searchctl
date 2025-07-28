package create

import (
	"bytes"
	"testing"

	"github.com/spf13/cobra"
)

func TestCreateCommand(t *testing.T) {
	cmd := &cobra.Command{Use: "test"}
	cmd.AddCommand(NewCreateCmd())

	// Test help output
	var buf bytes.Buffer
	cmd.SetOut(&buf)
	cmd.SetArgs([]string{"create", "--help"})
	
	err := cmd.Execute()
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	output := buf.String()
	if !bytes.Contains([]byte(output), []byte("Create a")) {
		t.Error("Expected help text to contain 'Create a'")
	}
}

func TestCreateCommandWithoutSubcommand(t *testing.T) {
	cmd := &cobra.Command{Use: "test"}
	cmd.AddCommand(NewCreateCmd())

	var buf bytes.Buffer
	cmd.SetOut(&buf)
	cmd.SetErr(&buf)
	cmd.SetArgs([]string{"create"})
	
	err := cmd.Execute()
	// Should show help when no subcommand is provided, not error
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
	
	output := buf.String()
	if !bytes.Contains([]byte(output), []byte("Available Commands")) {
		t.Error("Expected help output when no subcommand provided")
	}
}
