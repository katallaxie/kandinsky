package relay

import (
	"context"
	"io"
	"net"
	"net/http"
	"time"

	"github.com/katallaxie/pkg/server"

	ws "github.com/gorilla/websocket"
	log "github.com/sirupsen/logrus"
	"golang.org/x/sync/errgroup"
)

// Opt ...
type Opt func(*Opts)

// Opts ...
type Opts struct {
	Logger *log.Entry
	Relay  string
	Addr   string
}

type Relay interface {
	server.Listener
}

type relay struct {
	http *http.Server

	relay string
	addr  string

	log *log.Entry

	opts *Opts
}

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
func (r *relay) Start(ctx context.Context, ready server.ReadyFunc, run server.RunFunc) func() error {
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

func relayHandler(ll *log.Entry, relay string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		conn, err := net.Dial("tcp", relay)
		if err != nil {
			ll.Error(err)

			return
		}
		defer conn.Close()

		// upgrade connection from http to tcp
		c, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			ll.Error(err)

			return
		}
		defer c.Close()

		g, _ := errgroup.WithContext(r.Context())

		g.Go(readMessages(c, conn))
		g.Go(writeMessages(c, conn))

		time.Sleep(5 * time.Second)

		if err := g.Wait(); err != nil {
			// todo: log error
			ll.Error(err)
		}
	})
}

func writeMessages(w *ws.Conn, conn net.Conn) func() error {
	return func() error {
		for {

			writer, err := w.NextWriter(ws.BinaryMessage)
			if err != nil {
				return err
			}

			io.Copy(writer, conn)
		}
	}
}

func readMessages(w *ws.Conn, conn net.Conn) func() error {
	return func() error {
		for {
			// todo: check message type
			_, msg, err := w.NextReader()
			if err != nil {
				if ws.IsUnexpectedCloseError(err, ws.CloseGoingAway, ws.CloseAbnormalClosure) {
					return err
				}

				return err
			}

			io.Copy(conn, msg)
		}
	}
}

func configureHandler(r *relay) {
	r.http = &http.Server{
		Addr:              r.addr,
		Handler:           relayHandler(r.log, r.relay),
		ReadHeaderTimeout: 2 * time.Second,
	}
}

func configure(r *relay, opts ...Opt) {
	for _, o := range opts {
		o(r.opts)
	}

	r.log = r.opts.Logger
	r.addr = r.opts.Addr
}
