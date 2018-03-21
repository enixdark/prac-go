package configuration


import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/enixdark/prac-go/rest/src/lib/persistence/dblayer"
)

var (
	DBTypeDefault = dblayer.DBTYPE("mongodb")
	DBConnectionDefault = "mongodb://root:root@127.0.0.1:27017/rest"
	RestfulDefault = "localhost:3000"
)

type ServiceConfig struct {
	Databasetype dblayer.DBTYPE `json:"databasetype"`
	DBConnection string `json:"dbconnection"`
	RestfulEndpoint string `json:"restfulapi_endpoint"`
}

func ExtractConfiguration(filename string) (ServiceConfig, error) {
	conf := ServiceConfig{
		DBTypeDefault,
		DBConnectionDefault,
		RestfulDefault,
	}

	file, err := os.Open(filename)

	if err != nil {
		fmt.Println("Configuration file not found. Continuing with default values.")
		return conf, err
	}

	err = json.NewDecoder(file).Decode(&conf)
	return conf, err
}