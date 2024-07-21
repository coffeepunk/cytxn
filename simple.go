package cytxn

import (
	"context"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
)

func QueryRead(ds DatabaseService, statement Statement) (*neo4j.EagerResult, error) {
	ctx, cancel := context.WithCancel(ds.Context())
	defer cancel()
	return neo4j.ExecuteQuery(ctx, ds.Driver(),
		statement.Query,
		statement.Params, neo4j.EagerResultTransformer,
		neo4j.ExecuteQueryWithDatabase(ds.Name()),
		neo4j.ExecuteQueryWithReadersRouting())
}

func QueryWrite(ds DatabaseService, statement Statement) (*neo4j.EagerResult, error) {
	ctx, cancel := context.WithCancel(ds.Context())
	defer cancel()
	return neo4j.ExecuteQuery(ctx, ds.Driver(),
		statement.Query,
		statement.Params, neo4j.EagerResultTransformer,
		neo4j.ExecuteQueryWithDatabase(ds.Name()),
		neo4j.ExecuteQueryWithWritersRouting())
}
