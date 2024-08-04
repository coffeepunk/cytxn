package cytxn

import (
	"context"
	"encoding/json"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"log"
	"os"
)

// GetTestConfig reads the contents of the ".testconfig.json" file and unmarshals it into a TestConfig struct.
// If any error occurs during file reading or JSON parsing, a fatal log is outputted.
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

func GetTestDBConnection() DBConnection {
	config := GetTestConfig()
	return DBConnection{
		Target: config.Target,
		Name:   config.Name,
		Auth:   neo4j.BasicAuth(config.Username, config.Password, config.Realm),
	}
}

// CleanUp cleans up the Neo4j database by executing a DETACH DELETE all query.
// If any error occurs during driver creation, query execution, or driver closure, a fatal log is outputted.
// This function is typically used in test cases to clean up the database after running queries.
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

	defer func(driver neo4j.DriverWithContext, ctx context.Context) {
		err := driver.Close(ctx)
		if err != nil {
			log.Fatalf("Error closing driver: %v", err)
		}
	}(driver, ctx)

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
