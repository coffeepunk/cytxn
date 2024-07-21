package cytxn

import (
	"context"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"log"
)

type DatabaseService interface {
	Close()
	Context() context.Context
	Driver() neo4j.DriverWithContext
	Name() string
}

type DatabaseServiceConfig struct {
	Ctx               context.Context
	DatabaseName      string
	DriverWithContext neo4j.DriverWithContext
}

func (dsc DatabaseServiceConfig) Close() {
	err := dsc.DriverWithContext.Close(dsc.Ctx)
	if err != nil {
		log.Printf("Error closing database driver: %v", err)
	}
}

func (dsc DatabaseServiceConfig) Context() context.Context {
	return dsc.Ctx
}

func (dsc DatabaseServiceConfig) Driver() neo4j.DriverWithContext {
	return dsc.DriverWithContext
}

func (dsc DatabaseServiceConfig) Name() string {
	return dsc.DatabaseName
}

type DatabaseConnectionConfig struct {
	Target string
	Name   string
}

type BasicAuthCredentials struct {
	Username string
	Password string
	Realm    string
}

type Statement struct {
	Query  string
	Params map[string]interface{}
}

func NewDBService(details DatabaseConnectionConfig, auth BasicAuthCredentials) (DatabaseService, error) {
	driver, err := neo4j.NewDriverWithContext(details.Target, neo4j.BasicAuth(auth.Username, auth.Password, auth.Realm))
	if err != nil {
		return nil, err
	}

	var dsc DatabaseServiceConfig
	dsc.Ctx = context.Background()
	dsc.DatabaseName = details.Name
	dsc.DriverWithContext = driver

	return dsc, nil
}
