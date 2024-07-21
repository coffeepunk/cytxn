package cytxn

import "testing"

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
