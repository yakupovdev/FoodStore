package pkg

import (
	"context"
	"fmt"

	pg "github.com/jackc/pgx/v5"
)

type Config struct {
	Database string `env:"POSTGRES_DB"`
	HOST     string `env:"POSTGRES_URI"`
	Port     uint16 `env:"POSTGRESQL_PORT"`
	Username string `env:"POSTGRESQL_USERNAME"`
	Password string `env:"POSTGRESQL_PASSWORD"`
}

type DB struct {
	Conn *pg.Conn
}

func NewPostgresDB(ctx context.Context, cfg Config) (*DB, error) {
	opts, err := pg.ParseConfig("")
	if err != nil {
		return nil, fmt.Errorf("failed to parse config: %w", err)
	}

	opts.User = cfg.Username
	opts.Password = cfg.Password
	opts.Host = cfg.HOST
	opts.Port = cfg.Port
	opts.Database = cfg.Database

	db, err := pg.ConnectConfig(ctx, opts)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}
	EnsureSchema(ctx, db)
	EnsureSchemaRecoveryCodes(ctx, db)

	return &DB{Conn: db}, nil
}

func EnsureSchema(ctx context.Context, conn *pg.Conn) error {

	_, err := conn.Exec(ctx, `
CREATE TABLE IF NOT EXISTS users (
	id BIGSERIAL PRIMARY KEY,
	email TEXT NOT NULL UNIQUE,
	password TEXT NOT NULL,
	type TEXT NOT NULL,
	balance BIGINT NOT NULL DEFAULT 0,
	created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    last_enter TIMESTAMPTZ NOT NULL DEFAULT now()
);
`)
	if err != nil {
		return fmt.Errorf("ensure schema: %w", err)
	}
	return nil
}

func EnsureSchemaRecoveryCodes(ctx context.Context, conn *pg.Conn) error {
	stmts := []string{
		`
CREATE TABLE IF NOT EXISTS password_recovery_codes (
  user_id    BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
  email      TEXT   NOT NULL,
  code_hash  TEXT   NOT NULL,
  expired_at TIMESTAMPTZ NOT NULL,
  created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
  PRIMARY KEY (user_id, email)
);`,
		`CREATE INDEX IF NOT EXISTS idx_password_recovery_codes_email ON password_recovery_codes (email);`,
		`CREATE INDEX IF NOT EXISTS idx_password_recovery_codes_expired_at ON password_recovery_codes (expired_at);`,
	}

	for _, q := range stmts {
		if _, err := conn.Exec(ctx, q); err != nil {
			return fmt.Errorf("ensure schema recovery codes: %w", err)
		}
	}
	return nil
}
