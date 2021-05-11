package main

import (
	"flag"
	"fmt"
	"github.com/spacerouter/authentication_server/config"
	"github.com/spacerouter/authentication_server/server"
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
