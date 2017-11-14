package main

import (
	"flag"
	"log"

	"github.com/vedhavyas/prog-image"
)

var addr = flag.String("addr", ":8080", "server address")

func main() {
	log.SetFlags(log.Lshortfile | log.LstdFlags)
	flag.Parse()

	if *addr == "" {
		log.Fatalf("invalid server adress: %s", *addr)
	}

	progimg.StartImageServer(*addr)
}
