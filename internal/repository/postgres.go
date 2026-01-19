package repository

import (
	pg "github.com/jackc/pgx/v5"
)

type Postgres struct {
	Conn *pg.Conn
}

func NewPostgres(conn *pg.Conn) *Postgres {
	return &Postgres{
		Conn: conn,
	}
}
