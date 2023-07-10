package config

import (
	"syscall"
	"time"

	log "github.com/sirupsen/logrus"
)

const (
	// DefaultLogLevel is the default logging level.
	DefaultLogLevel = log.WarnLevel

	// DefaultTermSignal is the signal to term the agent.
	DefaultTermSignal = syscall.SIGTERM

	// DefaultReloadSignal is the default signal for reload.
	DefaultReloadSignal = syscall.SIGHUP

	// DefaultKillSignal is the default signal for termination.
	DefaultKillSignal = syscall.SIGINT

	// DefaultVerbose is the default verbosity.
	DefaultVerbose = false

	// DefaultStatusAddr is the default addrs for debug listener
	DefaultStatusAddr = ":8443"

	// DefaultAddr is s the default addrs to listen on
	DefaultAddr = ":8080"
)

// Config contains a configuration for Voskhod
type Config struct {
	// Verbose toggles the verbosity
	Verbose bool

	// LogLevel is the level with with to log for this config
	LogLevel log.Level

	// ReloadSignal
	ReloadSignal syscall.Signal

	// TermSignal
	TermSignal syscall.Signal

	// KillSignal
	KillSignal syscall.Signal

	// Timeout of the runtime
	Timeout time.Duration

	// StatusAddr is the addr of the debug listener
	StatusAddr string

	// Addr is the addr to listen for
	Addr string

	// Relay is the addr to relay to
	Relay string
}

// New returns a new Config
func New() *Config {
	return &Config{
		Verbose:      DefaultVerbose,
		LogLevel:     DefaultLogLevel,
		ReloadSignal: DefaultReloadSignal,
		TermSignal:   DefaultTermSignal,
		KillSignal:   DefaultKillSignal,
		StatusAddr:   DefaultStatusAddr,
		Addr:         DefaultAddr,
	}
}
