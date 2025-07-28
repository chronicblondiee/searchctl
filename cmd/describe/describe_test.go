package describe

import (
	"bytes"
	"testing"

	"github.com/spf13/cobra"
)

func TestDescribeCommand(t *testing.T) {
	cmd := &cobra.Command{Use: "test"}
	cmd.AddCommand(NewDescribeCmd())

	// Test help output
	var buf bytes.Buffer
	cmd.SetOut(&buf)
	cmd.SetArgs([]string{"describe", "--help"})
	
	err := cmd.Execute()
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	output := buf.String()
	if !bytes.Contains([]byte(output), []byte("detailed information")) {
		t.Error("Expected help text to contain 'detailed information'")
	}
}

func TestDescribeCommandWithoutSubcommand(t *testing.T) {
	cmd := &cobra.Command{Use: "test"}
	cmd.AddCommand(NewDescribeCmd())

	var buf bytes.Buffer
	cmd.SetOut(&buf)
	cmd.SetErr(&buf)
	cmd.SetArgs([]string{"describe"})
	
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
