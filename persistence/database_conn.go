package persistence

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v4/pgxpool"
)

type pgConnection struct {
	db *pgxpool.Pool
}

var dbconn pgConnection
var table string

func CloseDb() {
	dbconn.db.Close()
}

func TableData(tablename string) {
	table = tablename
}
func MakeConnection(url string) error {
	fmt.Println("Creating the Configuration")

	pgxpoolconfig, err := pgxpool.ParseConfig(url)

	if err != nil {
		return fmt.Errorf("Cannot parse the Database_url")
	}

	pgxpoolconfig.MaxConns = 20
	pgxpoolconfig.MinConns = 10
	pgxpoolconfig.MaxConnLifetime = 2 * time.Hour
	pgxpoolconfig.MaxConnIdleTime = 1 * time.Hour
	pgxpoolconfig.ConnConfig.ConnectTimeout = 5 * time.Second

	conn, err := pgconnection(pgxpoolconfig)

	if err != nil {
		return err
	}
	dbconn.db = conn
	return nil
}

func pgconnection(config *pgxpool.Config) (*pgxpool.Pool, error) {

	connpool, err := pgxpool.ConnectConfig(context.Background(), config)
	if err != nil {
		return nil, err
	}

	return connpool, err

}
