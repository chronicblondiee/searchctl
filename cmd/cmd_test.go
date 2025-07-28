package cmd_test

import (
	"bytes"
	"os"
	"strings"
	"testing"

	"github.com/chronicblondiee/searchctl/cmd"
)

func setupTestEnv(t *testing.T) func() {
	tmpDir, err := os.MkdirTemp("", "searchctl-test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}

	oldHome := os.Getenv("HOME")
	os.Setenv("HOME", tmpDir)

	return func() {
		os.Setenv("HOME", oldHome)
		os.RemoveAll(tmpDir)
	}
}

func TestRootCommandHelp(t *testing.T) {
	cleanup := setupTestEnv(t)
	defer cleanup()

	var buf bytes.Buffer
	rootCmd := cmd.NewRootCmd()
	rootCmd.SetOut(&buf)
	rootCmd.SetArgs([]string{"--help"})

	err := rootCmd.Execute()
	if err != nil {
		t.Errorf("Root command help failed: %v", err)
	}

	output := buf.String()
	if !strings.Contains(output, "searchctl") {
		t.Error("Expected help output to contain 'searchctl'")
	}
}

func TestVersionCommand(t *testing.T) {
	cleanup := setupTestEnv(t)
	defer cleanup()

	var buf bytes.Buffer
	rootCmd := cmd.NewRootCmd()
	rootCmd.SetOut(&buf)
	rootCmd.SetArgs([]string{"version"})

	err := rootCmd.Execute()
	if err != nil {
		t.Errorf("Version command failed: %v", err)
	}

	output := buf.String()
	if !strings.Contains(output, "version:") {
		t.Error("Expected version output to contain 'version:'")
	}
}

func TestVersionCommandJSON(t *testing.T) {
	cleanup := setupTestEnv(t)
	defer cleanup()

	var buf bytes.Buffer
	rootCmd := cmd.NewRootCmd()
	rootCmd.SetOut(&buf)
	rootCmd.SetArgs([]string{"version", "-o", "json"})

	err := rootCmd.Execute()
	if err != nil {
		t.Errorf("Version command with JSON failed: %v", err)
	}

	output := buf.String()
	if !strings.Contains(output, `"version"`) {
		t.Error("Expected JSON output to contain version field")
	}
}

func TestGetCommandHelp(t *testing.T) {
	cleanup := setupTestEnv(t)
	defer cleanup()

	var buf bytes.Buffer
	rootCmd := cmd.NewRootCmd()
	rootCmd.SetOut(&buf)
	rootCmd.SetArgs([]string{"get", "--help"})

	err := rootCmd.Execute()
	if err != nil {
		t.Errorf("Get command help failed: %v", err)
	}

	output := buf.String()
	if !strings.Contains(output, "indices") {
		t.Error("Expected get help to contain 'indices'")
	}
}

func TestCreateIndexDryRun(t *testing.T) {
	cleanup := setupTestEnv(t)
	defer cleanup()

	var buf bytes.Buffer
	rootCmd := cmd.NewRootCmd()
	rootCmd.SetOut(&buf)
	rootCmd.SetErr(&buf)
	rootCmd.SetArgs([]string{"create", "index", "test-index", "--dry-run"})

	err := rootCmd.Execute()
	if err != nil {
		t.Errorf("Create index dry-run failed: %v", err)
	}

	output := buf.String()
	if !strings.Contains(output, "Would create") || !strings.Contains(output, "test-index") {
		t.Errorf("Expected dry-run output, got: %s", output)
	}
}

func TestDeleteIndexDryRun(t *testing.T) {
	cleanup := setupTestEnv(t)
	defer cleanup()

	var buf bytes.Buffer
	rootCmd := cmd.NewRootCmd()
	rootCmd.SetOut(&buf)
	rootCmd.SetErr(&buf)
	rootCmd.SetArgs([]string{"delete", "index", "test-index", "--dry-run"})

	err := rootCmd.Execute()
	if err != nil {
		t.Errorf("Delete index dry-run failed: %v", err)
	}

	output := buf.String()
	if !strings.Contains(output, "Would delete") || !strings.Contains(output, "test-index") {
		t.Errorf("Expected dry-run output, got: %s", output)
	}
}
