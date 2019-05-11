package relay

import (
	"net/http"

	"github.com/andersnormal/pkg/server"

	log "github.com/sirupsen/logrus"
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
	http  *http.Server
	https *http.Server

	relay string
	addr  string

	log *log.Entry

	opts *Opts
}
