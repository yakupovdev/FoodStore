package initialization

import (
	"context"
	"fmt"

	pg "github.com/jackc/pgx/v5"
)

type Config struct {
	Database string
	Host     string
	Port     uint16
	User     string
	Password string
}

func NewConnection(ctx context.Context, cfg Config) (*pg.Conn, error) {
	opts, err := pg.ParseConfig("")
	if err != nil {
		return nil, fmt.Errorf("failed to parse config: %w", err)
	}

	opts.User = cfg.User
	opts.Password = cfg.Password
	opts.Host = cfg.Host
	opts.Port = cfg.Port
	opts.Database = cfg.Database

	conn, err := pg.ConnectConfig(ctx, opts)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to the database: %w", err)
	}

	return conn, nil
}
