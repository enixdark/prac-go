package main;

import (
	"flag"
	"fmt"
	"log"

	"github.com/enixdark/prac-go/rest/src/rest"
	"github.com/enixdark/prac-go/rest/src/lib/configuration"
	"github.com/enixdark/prac-go/rest/src/lib/persistence/dblayer"
)

func main() {
	confPath := flag.String("conf", `./configuration/config.json`, "flag to set the path to the configuration json file")
	flag.Parse()

	config, _ := configuration.ExtractConfiguration(*confPath)

	fmt.Println("Connecting to database")
	dbhandler, _ := dblayer.NewPersistenceLayer(config.Databasetype, config.DBConnection)

	log.Fatal(rest.ServeAPI(config.RestfulEndpoint, dbhandler))
}