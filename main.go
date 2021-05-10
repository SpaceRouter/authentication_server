package main

import (
	"authentification_server/config"
	"authentification_server/server"
	"flag"
	"fmt"
	"log"
	"os"
)

func main() {

	environment := flag.String("e", "dev", "")
	flag.Usage = func() {
		fmt.Println("Usage: server -e {mode}")
		os.Exit(1)
	}
	flag.Parse()
	config.Init(*environment)
	err := server.Init()
	if err != nil {
		log.Fatal(err)
	}
}
