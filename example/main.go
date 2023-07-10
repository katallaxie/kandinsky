package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"net"
	"net/url"
	"os"
	"os/signal"
	"time"

	"github.com/gorilla/websocket"
)

var (
	addr = flag.String("addr", "localhost:8080", "server addr")
	svc  = flag.String("sv", "localhost:4222", "service addr")
)

// nolint:gocyclo
func main() {
	flag.Parse()
	log.SetFlags(0)

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)

	l, err := net.Listen("tcp4", *svc)
	if err != nil {
		log.Fatal("listener: ", err)
	}
	defer l.Close()

	go func() {
		c, err := l.Accept()
		if err != nil {
			log.Fatal("connection: ", err)

			return
		}
		go handleConnection(c)
	}()

	u := url.URL{Scheme: "ws", Host: *addr, Path: "/"}
	log.Printf("connecting to %s", u.String())

	c, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		log.Panic("dial:", err)
	}
	defer c.Close()

	done := make(chan struct{})

	go func() {
		defer close(done)
		for {
			_, message, err := c.ReadMessage()
			if err != nil {
				log.Println("read:", err)
				return
			}
			log.Printf("recv: %s", string(message))
		}
	}()

	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-done:
			return
		case t := <-ticker.C:
			err := c.WriteMessage(websocket.BinaryMessage, []byte(t.String()))
			if err != nil {
				log.Println("write:", err)
				return
			}
		case <-interrupt:
			log.Println("interrupt")

			// Cleanly close the connection by sending a close message and then
			// waiting (with timeout) for the server to close the connection.
			err := c.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
			if err != nil {
				log.Println("write close:", err)
				return
			}
			select {
			case <-done:
			case <-time.After(time.Second):
			}
			return
		}
	}
}

func handleConnection(c net.Conn) {
	fmt.Printf("Serving %s\n", c.RemoteAddr().String())
	for {
		buf := make([]byte, 65536)
		_, err := bufio.NewReader(c).Read(buf)
		if err != nil {
			fmt.Println(err)
			return
		}

		log.Println("recv:", string(buf))

		_, err = bufio.NewWriter(c).WriteString("got it")
		if err != nil {
			fmt.Println(err)
			return
		}
	}
}
