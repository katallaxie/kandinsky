package main

import (
	"log"

	"github.com/katallaxie/kandinsky/internal/cmd"
)

func main() {
	if err := cmd.RootCmd.Execute(); err != nil {
		log.Fatal(err)
	}
}
