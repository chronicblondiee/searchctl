package config_test

import (
	"os"
	"testing"

	"github.com/chronicblondiee/searchctl/pkg/config"
)

func TestInitConfigDefault(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "searchctl-test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	// Set HOME to tmpDir to isolate test
	oldHome := os.Getenv("HOME")
	os.Setenv("HOME", tmpDir)
	defer os.Setenv("HOME", oldHome)

	err = config.InitConfig("")
	if err != nil {
		t.Errorf("InitConfig failed: %v", err)
	}

	cfg := config.GetConfig()
	if cfg == nil {
		t.Fatal("Expected config to be initialized")
	}

	if cfg.Kind != "Config" {
		t.Errorf("Expected Kind Config, got %s", cfg.Kind)
	}
	if cfg.CurrentContext != "default" {
		t.Errorf("Expected CurrentContext default, got %s", cfg.CurrentContext)
	}
}

func TestGetCurrentContext(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "searchctl-test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	oldHome := os.Getenv("HOME")
	os.Setenv("HOME", tmpDir)
	defer os.Setenv("HOME", oldHome)

	err = config.InitConfig("")
	if err != nil {
		t.Fatalf("InitConfig failed: %v", err)
	}

	ctx, err := config.GetCurrentContext()
	if err != nil {
		t.Errorf("GetCurrentContext failed: %v", err)
	}

	if ctx.Name != "default" {
		t.Errorf("Expected context name default, got %s", ctx.Name)
	}
}

func TestGetCluster(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "searchctl-test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	oldHome := os.Getenv("HOME")
	os.Setenv("HOME", tmpDir)
	defer os.Setenv("HOME", oldHome)

	err = config.InitConfig("")
	if err != nil {
		t.Fatalf("InitConfig failed: %v", err)
	}

	cluster, err := config.GetCluster("default")
	if err != nil {
		t.Errorf("GetCluster failed: %v", err)
	}

	if cluster.Cluster.Server != "http://localhost:9200" {
		t.Errorf("Expected server http://localhost:9200, got %s", cluster.Cluster.Server)
	}

	// Test non-existent cluster
	_, err = config.GetCluster("nonexistent")
	if err == nil {
		t.Error("Expected error for non-existent cluster")
	}
}
