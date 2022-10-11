package psql

import (
	"context"
	"fmt"
	"github.com/DawnKosmos/metapine/backend/exchange/psql/gen"
	"github.com/jackc/pgx/v4"
	"log"

	"github.com/jackc/pgx/v4/pgxpool"
)

var p *DB
var ctx = context.Background()

var ErrNoRows = pgx.ErrNoRows

type CustomLogger struct {
}

func (c CustomLogger) Write(p []byte) (n int, err error) {
	return fmt.Println(string(p))
}

type DB struct {
	q      *pgxpool.Pool
	qq     *gen.Queries
	loggin *log.Logger
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

	loggin := log.New(CustomLogger{}, "- ", log.Ltime)
	p = &DB{q: conn, qq: gen.New(conn), loggin: loggin}
}

// ReturnIndexList
// Although a single Ticker is also saved as Index. This Function returns only Index which are composits of more than 1 ticker
func ReturnIndexList() ([]string, error) {
	rows, err := p.qq.ReturnIndexList(ctx)
	if err != nil {
		return nil, err
	}
	var ss []string = []string{"Name \t ID \t composite_of"}
	for _, v := range rows {
		ss = append(ss, fmt.Sprintf("%s\t %d \t %d", v.Name, v.IndexID, v.CompositeOf))
	}
	return ss, nil
}

func ReturnIndex(id int32) (out Index, err error) {
	row, err := p.qq.ReturnIndex(ctx, id)
	if err != nil {
		return
	}
	out.name = row[0].Name
	for _, v := range row {
		out.Tickers = append(out.Tickers, Ticker{
			Exchange:      v.Exchange,
			Ticker:        v.Ticker,
			Weight:        v.Weight,
			ExcludeVolume: v.Excludevolume,
		})
	}
	return out, nil
}
