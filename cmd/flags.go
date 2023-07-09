package cmd

import (
	c "github.com/katallaxie/kandinsky/config"

	"github.com/spf13/cobra"
)

func addFlags(cmd *cobra.Command, cfg *c.Config) {
	// enable verbose output
	cmd.Flags().BoolVar(&cfg.Verbose, "verbose", c.DefaultVerbose, "enable verbose output")

	// set addr to listen on
	cmd.Flags().StringVar(&cfg.Addr, "addr", c.DefaultAddr, "address to listen on")

	// relay is the server to relay to
	cmd.Flags().StringVar(&cfg.Relay, "relay", "", "address to relay (e.g. nats:4222")
}
