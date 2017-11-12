package main

import (
	"log"

	"github.com/vedhavyas/prog-image"
)

func main() {
	log.SetFlags(log.Lshortfile | log.LstdFlags)
	progimg.StartImageServer(":8080")
}
