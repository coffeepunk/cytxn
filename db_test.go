package cytxn

import (
	"context"
	"encoding/json"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"log"
	"os"
	"testing"
)

type TestConfig struct {
	Target   string
	Name     string
	Username string
	Password string
	Realm    string
}

func GetTestConfig() TestConfig {
	confFile, errRead := os.ReadFile("./.testconfig.json")
	if errRead != nil {
		log.Fatalf("Error reading config file: %v", errRead)
	}

	var config TestConfig
	errJson := json.Unmarshal(confFile, &config)
	if errJson != nil {
		log.Fatalf("Error parsing config file: %v", errJson)
	}

	return config
}

func CleanUp() {
	s := Statement{
		Query:  `MATCH (n) DETACH DELETE n`,
		Params: map[string]interface{}{},
	}

	ctx := context.Background()
	c := GetTestConfig()

	driver, err := neo4j.NewDriverWithContext(c.Target, neo4j.BasicAuth(c.Username, c.Password, c.Realm))
	if err != nil {
		log.Fatalf("Error creating driver: %v", err)
	}
	defer driver.Close(ctx)

	ctx, cancel := context.WithCancel(ctx)
	defer cancel()
	_, errWrite := neo4j.ExecuteQuery(ctx, driver,
		s.Query,
		s.Params, neo4j.EagerResultTransformer,
		neo4j.ExecuteQueryWithDatabase(c.Name),
		neo4j.ExecuteQueryWithWritersRouting())

	if errWrite != nil {
		log.Fatalf("Error executing query: %v", errWrite)
	}
}

func DatabaseConnection() DatabaseService {
	config := GetTestConfig()

	details := DatabaseConnectionConfig{
		Target: config.Target,
		Name:   config.Name,
	}
	auth := BasicAuthCredentials{
		Username: config.Username,
		Password: config.Password,
		Realm:    config.Realm,
	}

	dbService, err := NewDBService(details, auth)
	if err != nil {
		log.Fatalf("Error creating database service: %v", err)
	}

	return dbService
}

func TestNewDBService(t *testing.T) {
	config := GetTestConfig()
	details := DatabaseConnectionConfig{
		Target: config.Target,
		Name:   config.Name,
	}

	auth := BasicAuthCredentials{
		Username: config.Username,
		Password: config.Password,
		Realm:    config.Realm,
	}

	ds, errService := NewDBService(details, auth)
	defer ds.Close()
	if errService != nil {
		t.Fatal(errService)
	}

	if err := ds.Driver().VerifyConnectivity(ds.Context()); err != nil {
		t.Errorf("Failed to connect to database: %v", err)
	}

	if ds.Name() != "cyphertest" {
		t.Errorf("got %v, want %v", ds.Name(), details.Name)
	}
}
