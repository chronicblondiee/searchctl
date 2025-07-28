package delete

import (
	"bytes"
	"testing"

	"github.com/spf13/cobra"
)

func TestDeleteCommand(t *testing.T) {
	cmd := &cobra.Command{Use: "test"}
	cmd.AddCommand(NewDeleteCmd())

	// Test help output
	var buf bytes.Buffer
	cmd.SetOut(&buf)
	cmd.SetArgs([]string{"delete", "--help"})
	
	err := cmd.Execute()
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	output := buf.String()
	if !bytes.Contains([]byte(output), []byte("Delete a")) {
		t.Error("Expected help text to contain 'Delete a'")
	}
}

func TestDeleteCommandWithoutSubcommand(t *testing.T) {
	cmd := &cobra.Command{Use: "test"}
	cmd.AddCommand(NewDeleteCmd())

	var buf bytes.Buffer
	cmd.SetOut(&buf)
	cmd.SetErr(&buf)
	cmd.SetArgs([]string{"delete"})
	
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
