package main

import (
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
)

func TestVersionStringStamped(t *testing.T) {
	original := version
	t.Cleanup(func() { version = original })

	version = "v9.9.9"
	if got := versionString(); got != "v9.9.9" {
		t.Errorf("versionString() = %q, want %q", got, "v9.9.9")
	}
}

func TestVersionStringFallback(t *testing.T) {
	original := version
	t.Cleanup(func() { version = original })

	version = "dev"
	got := versionString()
	if strings.TrimSpace(got) == "" {
		t.Error("versionString() returned an empty string")
	}
	if strings.Contains(got, "(devel)") {
		t.Errorf("versionString() = %q, should not contain %q", got, "(devel)")
	}
}

func TestVersionFlag(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping build-based test in short mode")
	}
	goBin, err := exec.LookPath("go")
	if err != nil {
		t.Skip("go toolchain not available")
	}

	binary := filepath.Join(t.TempDir(), "clidle")
	if runtime.GOOS == "windows" {
		binary += ".exe"
	}

	build := exec.Command(goBin, "build", "-ldflags", "-X main.version=v1.2.3", "-o", binary, ".")
	if out, err := build.CombinedOutput(); err != nil {
		t.Fatalf("could not build binary: %v\n%s", err, out)
	}

	for _, flag := range []string{"-version", "--version"} {
		out, err := exec.Command(binary, flag).Output()
		if err != nil {
			t.Fatalf("running %s: %v", flag, err)
		}
		if got := strings.TrimSpace(string(out)); got != "v1.2.3" {
			t.Errorf("%s printed %q, want %q", flag, got, "v1.2.3")
		}
	}

	// -version takes precedence over -serve, so no listener is bound.
	out, err := exec.Command(binary, "-serve", "127.0.0.1:0", "-version").Output()
	if err != nil {
		t.Fatalf("running -serve with -version: %v", err)
	}
	if got := strings.TrimSpace(string(out)); got != "v1.2.3" {
		t.Errorf("-serve with -version printed %q, want %q", got, "v1.2.3")
	}
}
