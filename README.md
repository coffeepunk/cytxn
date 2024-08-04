# cytxn - Cypher transactions
A convenience wrapper around the Neo4j Golang driver to make development using Cypher and Neo4j easier. 

Provides a consistent way of writing Cypher queries with or without parameters using specific write and read functions.

Early stage development. Do NOT use in production without knowing the risks.

## Installation
```
go get github.com/coffeepunk/cytxn
```

## Connect to the database
```
conn := DBConnection{
    Target: "Your connection target or URI",
    Name:   "Name of the database",
    Auth:   neo4j.BasicAuth(username, password, realm),
}

ds, errService := NewDBService(conn)
if errService != nil {
    log.Fatal(errService)
}


```

## Perform a read query
```
// Create a new statement.
var s cytxn.Statement
// Add the query
s.Query = `MATCH (u:User {
        id: $userId
    })
    RETURN u AS user`

// Add parameters or leave empty.
s.Params = map[string]interface{}{
    "userId":     "1234",
}
// Run the query and take care of the result and error.
res, err := cytxn.QueryWrite(ds, s)
```
