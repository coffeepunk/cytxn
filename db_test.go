package cytxn

import (
	"log"
	"testing"
)

type TestConfig struct {
	Target   string
	Name     string
	Username string
	Password string
	Realm    string
}

func GetTestDBService() DatabaseService {
	conn := GetTestDBConnection()
	dbService, err := NewDBService(conn)
	if err != nil {
		log.Fatalf("Error creating database service: %v", err)
	}

	return dbService
}

func TestNewDBServiceBasicAuth(t *testing.T) {
	conn := GetTestDBConnection()
	ds, errService := NewDBService(conn)
	if errService != nil {
		t.Fatal(errService)
	}
	defer ds.Close()

	if err := ds.Driver().VerifyConnectivity(ds.Context()); err != nil {
		t.Errorf("Failed to connect to database: %v", err)
	}

	if ds.Name() != conn.Name {
		t.Errorf("got %v, want %v", ds.Name(), conn.Name)
	}
}
