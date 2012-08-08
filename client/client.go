package main

import (
	"github.com/kisielk/gcfg3/gcfg3"
	"log"
)

func main() {
	client, err := gcfg3.NewGcfg3Client("localhost:1234")
	if err != nil {
		log.Fatal("dialing:", err)
	}
	log.Println("Connected")

	probes, err := client.GetProbes()
	if err != nil {
		log.Fatal("GetProbes:", err)
	}

	for _, probe := range probes {
		log.Printf("%s", probe)
	}
}
