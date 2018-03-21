package dblayer

import (
	"github.com/enixdark/prac-go/rest/src/eventsservice/lib/persistence"
	"github.com/enixdark/prac-go/rest/src/eventsservice/lib/persistence/mongolayer"
)

type DBTYPE string

const (
	MONGODB DBTYPE = "mongodb"
	DYNAMODB DBTYPE = "dynamodb"
)

func NewPersistenceLayer(options DBTYPE, connection string) (persistence.DatabaseHandler, error) {
	switch expression {
	case MONGODB:
		return mongolayer.NewMongoDBLayer(connection)
	}
	return nil, nil
}