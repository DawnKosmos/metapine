package psql

import (
	"context"
	"fmt"
	"github.com/DawnKosmos/metapine/backend/exchange/psql/gen"

	"github.com/jackc/pgx/v4/pgxpool"
)

var p *DB
var ctx = context.Background()

type DB struct {
	q  *pgxpool.Pool
	qq *gen.Queries
	//Exchange []exchange.Exchange Collection of exchange connection
}

func (d *DB) Ping() error {
	return d.q.Ping(ctx)
}

func SetPSQL(host, user, databaseName, password string, port int) {
	params := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, databaseName)
	//	pg, err := sql.Open("postgres", params)
	conn, err := pgxpool.Connect(ctx, params)
	if err != nil {
		panic(err)
	}

	p = &DB{q: conn, qq: gen.New(conn)}
}

//Being able to register indexes
