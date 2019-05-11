package cmd

import (
	"context"

	"github.com/andersnormal/kandinsky/relay"
	"github.com/andersnormal/pkg/server"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

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
	s := server.NewServer(ctx)

	// create relay
	r := relay.New(
		cfg.Addr,
		relay.WithAddr(cfg.Addr),
		relay.WithLog(root.logger),
	)
	s.Listen(r)

	// listen for the server and wait for it to fail,
	// or for sys interrupts
	if err := s.Wait(); err != nil {
		root.logger.Error(err)
	}

	// noop
	return nil
}
