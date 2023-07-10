package cmd

import (
	"context"

	"github.com/katallaxie/kandinsky/internal/config"
	"github.com/katallaxie/kandinsky/internal/relay"
	"github.com/katallaxie/pkg/server"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

type root struct {
	logger *log.Entry
}

var cfg = config.New()

// RootCmd represents the base command when called without any subcommands
var RootCmd = &cobra.Command{
	Use:  "kandinsky",
	RunE: runE,
}

func init() {
	// initialize cobra
	cobra.OnInitialize(initConfig)

	RootCmd.Flags().BoolVar(&cfg.Verbose, "verbose", cfg.Verbose, "enable verbose output")
	RootCmd.Flags().StringVar(&cfg.Addr, "addr", cfg.Addr, "address to listen on")
	RootCmd.Flags().StringVar(&cfg.Relay, "relay", cfg.Relay, "address to relay (e.g. nats:4222")
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	viper.AutomaticEnv() // read in environment variables that match

	// Log as JSON instead of the default ASCII formatter.
	log.SetFormatter(&log.JSONFormatter{})

	// Only log the warning severity or above.
	log.SetLevel(cfg.LogLevel)

	// if we should output verbose
	if cfg.Verbose {
		log.SetLevel(log.InfoLevel)
	}
}

func runE(cmd *cobra.Command, args []string) error {
	// create a new root
	root := new(root)

	// init logger
	root.logger = log.WithFields(log.Fields{
		"verbose": cfg.Verbose,
	})

	// create root context
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// create server
	s, _ := server.WithContext(ctx)

	// create relay
	r := relay.New(
		cfg.Relay,
		relay.WithAddr(cfg.Addr),
		relay.WithLog(root.logger),
	)
	s.Listen(r, false)

	// listen for the server and wait for it to fail,
	// or for sys interrupts
	if err := s.Wait(); err != nil {
		root.logger.Error(err)
	}

	// noop
	return nil
}
