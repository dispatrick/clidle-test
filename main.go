package main

import (
	"context"
	_ "embed"
	"flag"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"

	"database/sql"

	"github.com/adrg/xdg"
	"github.com/ajeetdsouza/clidle/server"
	"github.com/ajeetdsouza/clidle/store"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/ssh"

	"golang.org/x/exp/slog"
	_ "modernc.org/sqlite"
)

var (
	// pathClidle is the path to the local data directory.
	// This is usually set to ~/.local/share/clidle on most UNIX systems.
	pathClidle  string
	pathStore   string
	pathHostKey string

	//go:embed schema.sql
	schemaSQL string
)

func init() {
	pathClidle = os.Getenv("CLIDLE_DATA_DIR")
	if pathClidle == "" {
		pathClidle = filepath.Join(xdg.DataHome, "clidle")
	}

	pathStore = filepath.Join(pathClidle, "clidle.db")
	pathHostKey = filepath.Join(pathClidle, "hostkey")
}

func main() {
	var (
		flagSSH     = flag.Bool("ssh", false, "Spawns an SSH server instead of running the game locally")
		flagHost    = flag.String("host", "", "Address the SSH server binds to (env: "+server.EnvHost+", default: "+server.DefaultHost+")")
		flagPort    = flag.String("port", "", "Port the SSH server listens on (env: "+server.EnvPort+", default: "+server.DefaultPort+")")
		flagHostKey = flag.String("host-key", "", "Path to the SSH host key, generated if missing (env: "+server.EnvHostKeyPath+")")
		flagServe   = flag.String("serve", "", "Spawns an SSH server on the given address (format: 0.0.0.0:1337)")
	)

	// Allow `clidle serve` in addition to `clidle --ssh`.
	args := os.Args[1:]
	serve := false
	if len(args) > 0 && args[0] == "serve" {
		serve = true
		args = args[1:]
	}
	if err := flag.CommandLine.Parse(args); err != nil {
		os.Exit(2)
	}

	var err error
	if serve || *flagSSH || *flagServe != "" {
		err = runServer(server.Flags{
			Host:        *flagHost,
			Port:        *flagPort,
			HostKeyPath: *flagHostKey,
			Address:     *flagServe,
		})
	} else {
		err = runCLI()
	}
	if err != nil {
		slog.Error("error running application", "error", slog.Any("error", err))
		os.Exit(1)
	}
}

func runCLI() error {
	ctx := context.Background()
	model, err := getModel(ctx, nil)
	if err != nil {
		return err
	}
	program := tea.NewProgram(model, tea.WithAltScreen(), tea.WithOutput(os.Stderr))

	_, err = program.Run()
	return err
}

func runServer(flags server.Flags) error {
	cfg, err := server.ResolveConfig(flags, pathHostKey)
	if err != nil {
		return err
	}

	srv, err := server.New(cfg, func(session ssh.Session, renderer *lipgloss.Renderer) (tea.Model, error) {
		// Every session gets its own model, so that concurrent players do not
		// share any state.
		model, err := getModel(session.Context(), renderer)
		if err != nil {
			return nil, err
		}
		if pty, _, active := session.Pty(); active {
			model.setSize(pty.Window.Width, pty.Window.Height)
		}
		return model, nil
	})
	if err != nil {
		return err
	}

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	return server.ListenAndServe(ctx, srv)
}

func getModel(ctx context.Context, renderer *lipgloss.Renderer) (*model, error) {
	dictionary := EnglishDictionary
	store, err := getStore()
	if err != nil {
		return nil, err
	}
	return newModel(ctx, store, dictionary, renderer), nil
}

func getStore() (*store.Queries, error) {
	if err := os.MkdirAll(pathClidle, 0700); err != nil {
		return nil, err
	}

	db, err := sql.Open("sqlite", pathStore)
	if err != nil {
		return nil, err
	}
	db.SetMaxOpenConns(1) // SQLite does not support concurrent writes
	if _, err := db.Exec(schemaSQL); err != nil {
		return nil, err
	}
	return store.New(db), nil
}
