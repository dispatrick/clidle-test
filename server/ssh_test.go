package server

import (
	"context"
	"path/filepath"
	"testing"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/ssh"
)

func TestResolveConfig(t *testing.T) {
	tests := []struct {
		name     string
		flags    Flags
		env      map[string]string
		wantAddr string
	}{
		{
			name:     "defaults",
			wantAddr: DefaultHost + ":" + DefaultPort,
		},
		{
			name:     "env overrides defaults",
			env:      map[string]string{EnvHost: "0.0.0.0", EnvPort: "2222"},
			wantAddr: "0.0.0.0:2222",
		},
		{
			name:     "flags override env",
			flags:    Flags{Host: "127.0.0.2", Port: "2323"},
			env:      map[string]string{EnvHost: "0.0.0.0", EnvPort: "2222"},
			wantAddr: "127.0.0.2:2323",
		},
		{
			name:     "address overrides host and port",
			flags:    Flags{Host: "127.0.0.2", Port: "2323", Address: "0.0.0.0:22"},
			wantAddr: "0.0.0.0:22",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			for key, value := range test.env {
				t.Setenv(key, value)
			}

			cfg, err := ResolveConfig(test.flags, "hostkey")
			if err != nil {
				t.Fatalf("ResolveConfig() returned an error: %v", err)
			}
			if got := cfg.Address(); got != test.wantAddr {
				t.Errorf("Address() = %q, want %q", got, test.wantAddr)
			}
			if cfg.HostKeyPath != "hostkey" {
				t.Errorf("HostKeyPath = %q, want %q", cfg.HostKeyPath, "hostkey")
			}
		})
	}
}

func TestResolveConfigHostKeyPath(t *testing.T) {
	t.Setenv(EnvHostKeyPath, "/env/hostkey")

	cfg, err := ResolveConfig(Flags{}, "/default/hostkey")
	if err != nil {
		t.Fatalf("ResolveConfig() returned an error: %v", err)
	}
	if cfg.HostKeyPath != "/env/hostkey" {
		t.Errorf("HostKeyPath = %q, want %q", cfg.HostKeyPath, "/env/hostkey")
	}

	cfg, err = ResolveConfig(Flags{HostKeyPath: "/flag/hostkey"}, "/default/hostkey")
	if err != nil {
		t.Fatalf("ResolveConfig() returned an error: %v", err)
	}
	if cfg.HostKeyPath != "/flag/hostkey" {
		t.Errorf("HostKeyPath = %q, want %q", cfg.HostKeyPath, "/flag/hostkey")
	}
}

func TestResolveConfigInvalidAddress(t *testing.T) {
	if _, err := ResolveConfig(Flags{Address: "not-an-address"}, "hostkey"); err == nil {
		t.Error("ResolveConfig() with an invalid address did not return an error")
	}
}

func TestNew(t *testing.T) {
	cfg := Config{
		Host:        "127.0.0.1",
		Port:        "0",
		HostKeyPath: filepath.Join(t.TempDir(), "hostkey"),
	}

	server, err := New(cfg, func(ssh.Session, *lipgloss.Renderer) (tea.Model, error) {
		return nil, nil
	})
	if err != nil {
		t.Fatalf("New() returned an error: %v", err)
	}
	if server.Addr != cfg.Address() {
		t.Errorf("Addr = %q, want %q", server.Addr, cfg.Address())
	}
	if server.PublicKeyHandler == nil {
		t.Error("PublicKeyHandler is not set")
	}
	if server.KeyboardInteractiveHandler == nil {
		t.Error("KeyboardInteractiveHandler is not set")
	}
}

func TestListenAndServeShutsDownOnContextCancel(t *testing.T) {
	cfg := Config{
		Host:        "127.0.0.1",
		Port:        "0",
		HostKeyPath: filepath.Join(t.TempDir(), "hostkey"),
	}

	server, err := New(cfg, func(ssh.Session, *lipgloss.Renderer) (tea.Model, error) {
		return nil, nil
	})
	if err != nil {
		t.Fatalf("New() returned an error: %v", err)
	}

	ctx, cancel := context.WithCancel(context.Background())
	errc := make(chan error, 1)
	go func() { errc <- ListenAndServe(ctx, server) }()

	// Give the server a moment to start listening, then shut it down.
	time.Sleep(100 * time.Millisecond)
	cancel()

	select {
	case err := <-errc:
		if err != nil {
			t.Errorf("ListenAndServe() returned an error: %v", err)
		}
	case <-time.After(10 * time.Second):
		t.Error("ListenAndServe() did not return after the context was cancelled")
	}
}
