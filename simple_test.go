package cytxn

import (
	"encoding/json"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"log"
	"os"
	"reflect"
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

func TestQueryRead(t *testing.T) {
	want := neo4j.EagerResult{
		Keys:    []string{"res"},
		Records: []*neo4j.Record{},
		Summary: nil,
	}

	s := Statement{
		Query:  "MATCH (n) RETURN n AS res LIMIT 10",
		Params: map[string]interface{}{},
	}

	ds := DatabaseConnection()
	defer ds.Close()

	result, err := QueryRead(ds, s)
	if err != nil {
		t.Error(err)
	}

	if want.Keys != nil && !reflect.DeepEqual(result.Keys, want.Keys) {
		t.Errorf("got %v, want %v", result, &want)
	}

	if want.Records != nil && len(want.Records) != len(result.Records) {
		t.Errorf("got %v, want %v", result, &want)
	}
}
