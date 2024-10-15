package main

import (
	"flag"
	"log"
)

func init() {
	log.SetPrefix("Blockchain Node: ")
}

func main() {
	port := flag.Uint("port", 3333, "TCP port number for blockchain node")
	flag.Parse()

	app := NewBlockchianNode(uint16(*port))
	log.Default().Println("Starting blockchain node on port", *port)
	app.Run()

}
