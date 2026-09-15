// Package server serves a Bubble Tea application over SSH.
package server

import (
	"context"
	"net"
	"os"
	"time"

	"log/slog"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/ssh"
	"github.com/charmbracelet/wish"
	"github.com/charmbracelet/wish/activeterm"
	wtea "github.com/charmbracelet/wish/bubbletea"
	"github.com/charmbracelet/wish/logging"
	wrecover "github.com/charmbracelet/wish/recover"
	"github.com/muesli/termenv"
	"github.com/pkg/errors"
	gossh "golang.org/x/crypto/ssh"
)

const (
	// DefaultHost is the address the server binds to when nothing else is
	// configured. It is deliberately loopback-only: set --host/SSH_HOST to
	// 0.0.0.0 to expose the server publicly.
	DefaultHost = "127.0.0.1"
	// DefaultPort is the port the server listens on by default.
	DefaultPort = "23234"

	// EnvHost, EnvPort and EnvHostKeyPath are the environment variables used
	// as a fallback when the corresponding flag is not set.
	EnvHost        = "SSH_HOST"
	EnvPort        = "SSH_PORT"
	EnvHostKeyPath = "SSH_HOST_KEY_PATH"

	// idleTimeout is how long a session may stay idle before being closed.
	idleTimeout = 30 * time.Minute
	// shutdownTimeout is how long a graceful shutdown may take.
	shutdownTimeout = 30 * time.Second
)

// Config is the resolved configuration of the SSH server.
type Config struct {
	Host        string
	Port        string
	HostKeyPath string
}

// Address returns the address the server listens on.
func (c Config) Address() string {
	return net.JoinHostPort(c.Host, c.Port)
}

// Flags holds the raw command line values. An empty value means the flag was
// not set, in which case the environment variable and then the default value
// are used.
type Flags struct {
	Host        string
	Port        string
	HostKeyPath string
	// Address is a legacy alias that sets both host and port at once
	// (format: 0.0.0.0:1337). It takes precedence over Host and Port.
	Address string
}

// ResolveConfig resolves the server configuration, preferring flags over
// environment variables over defaults. defaultHostKeyPath is used when neither
// the flag nor the environment variable is set; wish generates an ed25519 key
// at that path on first run.
func ResolveConfig(flags Flags, defaultHostKeyPath string) (Config, error) {
	cfg := Config{
		Host:        firstNonEmpty(flags.Host, os.Getenv(EnvHost), DefaultHost),
		Port:        firstNonEmpty(flags.Port, os.Getenv(EnvPort), DefaultPort),
		HostKeyPath: firstNonEmpty(flags.HostKeyPath, os.Getenv(EnvHostKeyPath), defaultHostKeyPath),
	}

	if flags.Address != "" {
		host, port, err := net.SplitHostPort(flags.Address)
		if err != nil {
			return Config{}, errors.Wrapf(err, "invalid address: %s", flags.Address)
		}
		cfg.Host = host
		cfg.Port = port
	}

	return cfg, nil
}

// Handler builds a fresh model for an SSH session. It is called once per
// session, so the returned model must not share mutable state with other
// sessions. The renderer is derived from the client's terminal and should be
// used for all styling, so that colors match the client and not the server.
//
// TODO: the session's public key fingerprint (sess.PublicKey()) could be used
// to key persistent per-user state; remote sessions are ephemeral for now.
type Handler func(sess ssh.Session, renderer *lipgloss.Renderer) (tea.Model, error)

// New builds an SSH server that serves the Bubble Tea application returned by
// handler. Any public key (and keyboard-interactive) authentication is
// accepted: the game is meant to be publicly playable.
func New(cfg Config, handler Handler) (*ssh.Server, error) {
	server, err := wish.NewServer(
		wish.WithAddress(cfg.Address()),
		wish.WithHostKeyPath(cfg.HostKeyPath),
		wish.WithIdleTimeout(idleTimeout),
		wish.WithPublicKeyAuth(func(ssh.Context, ssh.PublicKey) bool { return true }),
		wish.WithKeyboardInteractiveAuth(func(ssh.Context, gossh.KeyboardInteractiveChallenge) bool { return true }),
		wish.WithMiddleware(
			// Middlewares run in reverse order: logging first, then the PTY
			// check, then the application, all wrapped in a panic handler so
			// that a broken session cannot take down the server.
			wrecover.Middleware(
				wtea.MiddlewareWithProgramHandler(programHandler(handler), termenv.ANSI256),
				activeterm.Middleware(),
				logging.Middleware(),
			),
		),
	)
	return server, errors.Wrap(err, "could not create server")
}

// ListenAndServe runs the server until ctx is cancelled or the server fails,
// then shuts it down gracefully.
func ListenAndServe(ctx context.Context, server *ssh.Server) error {
	errc := make(chan error, 1)

	slog.Info("starting SSH server", slog.String("address", server.Addr))
	go func() {
		if err := server.ListenAndServe(); err != nil && !errors.Is(err, ssh.ErrServerClosed) {
			errc <- errors.Wrap(err, "server returned an error")
			return
		}
		errc <- nil
	}()

	select {
	case err := <-errc:
		return err
	case <-ctx.Done():
	}

	slog.Info("stopping SSH server")
	shutdownCtx, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
	defer cancel()

	if err := server.Shutdown(shutdownCtx); err != nil && !errors.Is(err, ssh.ErrServerClosed) {
		return errors.Wrap(err, "could not shutdown server")
	}
	return nil
}

// programHandler creates a new Bubble Tea program for every session.
func programHandler(handler Handler) wtea.ProgramHandler {
	return func(sess ssh.Session) *tea.Program {
		model, err := handler(sess, wtea.MakeRenderer(sess))
		if err != nil {
			slog.Error("could not create model", slog.Any("error", err))
			wish.Fatalln(sess, "could not create model:", err)
			return nil
		}

		options := append(wtea.MakeOptions(sess), tea.WithAltScreen())
		return tea.NewProgram(model, options...)
	}
}

func firstNonEmpty(values ...string) string {
	for _, value := range values {
		if value != "" {
			return value
		}
	}
	return ""
}
