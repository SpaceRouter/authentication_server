package main

import (
	"flag"
	"fmt"
	"github.com/spacerouter/authentication_server/config"
	"github.com/spacerouter/authentication_server/server"
	"log"
	"os"
	"os/user"
)

// @title SpaceRouter Authentication Server
// @version 0.1
// @description Authentication Server API.

// @contact.name ESIEESPACE Network
// @contact.url http://esieespace.fr
// @contact.email contact@esieespace.fr

// @license.name GPL-3.0
// @license.url https://github.com/SpaceRouter/authentication_server/blob/louis/LICENSE

// @host localhost:8080
// @BasePath /
// @query.collection.format multi

// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization
func main() {
	user, err := user.Current()
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("Run as %s \n", user.Username)

	environment := flag.String("e", "dev", "")
	flag.Usage = func() {
		fmt.Println("Usage: server -e {mode}")
		os.Exit(1)
	}

	flag.Parse()
	config.Init(*environment)
	err = server.Init()
	if err != nil {
		log.Fatal(err)
	}
}
