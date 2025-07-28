package version_test

import (
	"runtime"
	"testing"

	"github.com/chronicblondiee/searchctl/internal/version"
)

func TestVersionGet(t *testing.T) {
	info := version.Get()

	if info.Version == "" {
		t.Error("Expected version to be set")
	}
	if info.Commit == "" {
		t.Error("Expected commit to be set")
	}
	if info.Date == "" {
		t.Error("Expected date to be set")
	}
	if info.GoVersion == "" {
		t.Error("Expected goVersion to be set")
	}

	expectedGoVersion := runtime.Version()
	if info.GoVersion != expectedGoVersion {
		t.Errorf("Expected goVersion %s, got %s", expectedGoVersion, info.GoVersion)
	}
}

func TestVersionDefaults(t *testing.T) {
	info := version.Get()

	// These are the default values when not built with ldflags
	if info.Version != "dev" {
		t.Logf("Version is %s (may be set by build)", info.Version)
	}
	if info.Commit != "none" {
		t.Logf("Commit is %s (may be set by build)", info.Commit)
	}
	if info.Date != "unknown" {
		t.Logf("Date is %s (may be set by build)", info.Date)
	}
}
