package storage

import (
	"context"
	"fmt"

	pg "github.com/jackc/pgx/v5"
)

type Config struct {
	Database string `env:"POSTGRES_DB"`
	Host     string `env:"POSTGRES_URI"`
	Port     uint16 `env:"POSTGRESQL_PORT"`
	Username string `env:"POSTGRESQL_USERNAME"`
	Password string `env:"POSTGRESQL_PASSWORD"`
}

func NewPostgresDB(ctx context.Context, cfg Config) (*pg.Conn, error) {
	opts, err := pg.ParseConfig("")
	if err != nil {
		return nil, fmt.Errorf("failed to parse config: %w", err)
	}

	opts.User = cfg.Username
	opts.Password = cfg.Password
	opts.Host = cfg.Host
	opts.Port = cfg.Port
	opts.Database = cfg.Database

	conn, err := pg.ConnectConfig(ctx, opts)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	return conn, nil
}

func InitSchema(ctx context.Context, conn *pg.Conn) error {
	if err := ensureUsersSchema(ctx, conn); err != nil {
		return err
	}
	if err := ensureRecoveryCodesSchema(ctx, conn); err != nil {
		return err
	}
	return nil
}

func ensureUsersSchema(ctx context.Context, conn *pg.Conn) error {
	_, err := conn.Exec(ctx, `
CREATE TABLE IF NOT EXISTS users (
	userid BIGSERIAL PRIMARY KEY,
	email TEXT NOT NULL UNIQUE,
	password_hash TEXT NOT NULL,
	type TEXT NOT NULL,
	balance BIGINT NOT NULL DEFAULT 0,
	created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
	last_enter TIMESTAMPTZ NOT NULL DEFAULT now()
);
`)
	if err != nil {
		return fmt.Errorf("ensure users schema: %w", err)
	}
	return nil
}

func ensureRecoveryCodesSchema(ctx context.Context, conn *pg.Conn) error {
	stmts := []string{
		`
CREATE TABLE IF NOT EXISTS password_recovery_codes (
	userid BIGINT NOT NULL REFERENCES users(userID) ON DELETE CASCADE,
	email TEXT NOT NULL,
	code_hash TEXT NOT NULL,
	expired_at TIMESTAMPTZ NOT NULL,
	created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
	PRIMARY KEY (userID, email)
);`,
		`CREATE INDEX IF NOT EXISTS idx_prc_email ON password_recovery_codes (email);`,
		`CREATE INDEX IF NOT EXISTS idx_prc_expired_at ON password_recovery_codes (expired_at);`,
	}

	for _, stmt := range stmts {
		if _, err := conn.Exec(ctx, stmt); err != nil {
			return fmt.Errorf("ensure recovery codes schema: %w", err)
		}
	}

	return nil
}
