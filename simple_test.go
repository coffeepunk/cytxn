package cytxn

import (
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"reflect"
	"testing"
)

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

	ds := GetTestDBService()
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

func TestQueryWrite(t *testing.T) {
	ds := GetTestDBService()
	defer ds.Close()
	defer CleanUp()

	s := Statement{
		Query:  `CREATE (t:TestNode {id: "1234", name: "This is a test node"}) RETURN t AS node`,
		Params: map[string]interface{}{},
	}

	result, err := QueryWrite(ds, s)
	if err != nil {
		t.Error(err)
	}

	want := neo4j.EagerResult{
		Keys:    []string{"node"},
		Records: []*neo4j.Record{},
		Summary: nil,
	}

	if want.Keys != nil && !reflect.DeepEqual(result.Keys, want.Keys) {
		t.Errorf("got %v, want %v", result.Keys, want.Keys)
	}

	if len(result.Records) != 1 {
		t.Errorf("got %v, want %v", len(result.Records), 1)
	}

	node, _ := result.Records[0].Get("node")
	props := node.(neo4j.Node).Props
	if props["id"] != "1234" {
		t.Errorf("got %v, want %v", props["id"], "1234")
	}

	if props["name"] != "This is a test node" {
		t.Errorf("got %v, want %v", props["name"], "This is a test node")
	}
}
