package impl

import (
	pg "github.com/jackc/pgx/v5"
)

type ModeratorRepo struct {
	conn *pg.Conn
}

func NewModeratorRepo(conn *pg.Conn) *ModeratorRepo {
	return &ModeratorRepo{conn: conn}
}
