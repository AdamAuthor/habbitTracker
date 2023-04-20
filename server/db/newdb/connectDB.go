package newdb

import (
	"awesomeProject/server/db"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/jmoiron/sqlx"
	"log"
)

type DB struct {
	conn  *sqlx.DB
	habit *HabitRepository
	user  *UserRepository
}

func NewDB() db.Database {
	return &DB{}
}

func (D *DB) Connect(url string) error {
	conn, err := sqlx.Connect("pgx", url)
	if err != nil {
		return err
	}

	if err := conn.Ping(); err != nil {
		return err
	}

	D.conn = conn
	log.Println("DB has successful pinging")
	return nil
}

func (D *DB) Close() error {
	return D.conn.Close()
}
