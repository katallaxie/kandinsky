package relay

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/andersnormal/pkg/server"

	ws "github.com/gorilla/websocket"
	log "github.com/sirupsen/logrus"
)

var _ server.Listener = (*relay)(nil)

// New ...
func New(addr string, opts ...Opt) Relay {
	options := new(Opts)

	r := new(relay)
	r.opts = options
	r.relay = addr

	configure(r, opts...)
	configureHandler(r)

	return r
}

// Start ...
func (r *relay) Start(ctx context.Context) func() error {
	return func() error {
		// todo: support TLS
		if err := r.http.ListenAndServe(); err != nil {
			return err
		}

		return nil
	}
}

// Stop ...
func (r *relay) Stop() error {
	if r.http != nil {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		if err := r.http.Shutdown(ctx); err != nil {
			return err
		}
	}

	return nil
}

// WithLog ...
func WithLog(ll *log.Entry) func(o *Opts) {
	return func(o *Opts) {
		o.Logger = ll
	}
}

// WithAddr ...
func WithAddr(addr string) func(o *Opts) {
	return func(o *Opts) {
		o.Addr = addr
	}
}

var upgrader = ws.Upgrader{}

func relayHandler(addr string, relay string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			return
		}
		defer conn.Close()
	})
}

func configureHandler(r *relay) error {
	fmt.Println(r.addr)
	r.http = &http.Server{
		Addr:    r.addr,
		Handler: relayHandler(r.addr, r.relay),
	}

	return nil
}

func configure(r *relay, opts ...Opt) error {
	for _, o := range opts {
		o(r.opts)
	}

	r.log = r.opts.Logger
	r.addr = r.opts.Addr

	return nil
}
