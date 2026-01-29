package storage

import (
	"context"
	"fmt"

	pg "github.com/jackc/pgx/v5"
)

type Config struct {
	Database string `.env:"POSTGRES_DB"`
	Host     string `.env:"POSTGRES_URI"`
	Port     uint16 `.env:"POSTGRESQL_PORT"`
	Username string `.env:"POSTGRESQL_USERNAME"`
	Password string `.env:"POSTGRESQL_PASSWORD"`
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
		return nil, ErrDatabaseConnection
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
	if err := ensureWhitelistSchema(ctx, conn); err != nil {
		return err
	}
	if err := ensureBlacklistSchema(ctx, conn); err != nil {
		return err
	}
	if err := ensureCategoriesSchema(ctx, conn); err != nil {
		return err
	}
	if err := ensureProductsSchema(ctx, conn); err != nil {
		return err
	}
	if err := ensureClientsSchema(ctx, conn); err != nil {
		return err
	}
	if err := ensureSellersSchema(ctx, conn); err != nil {
		return err
	}
	if err := ensureSellerOffersSchema(ctx, conn); err != nil {
		return err
	}
	if err := ensureOrdersSchema(ctx, conn); err != nil {
		return err
	}
	if err := ensureOrdersItemsSchema(ctx, conn); err != nil {
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
		return ErrUsersSchema
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
			return ErrRecoveryCodesSchema
		}
	}

	return nil
}

func ensureWhitelistSchema(ctx context.Context, conn *pg.Conn) error {
	_, err := conn.Exec(ctx, `
CREATE TABLE IF NOT EXISTS token_whitelist (
	userid BIGINT NOT NULL REFERENCES users(userID) ON DELETE CASCADE,
	access_token_hash TEXT NOT NULL,
    expired_at TIMESTAMPTZ NOT NULL,
	created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
	PRIMARY KEY (userID, access_token_hash)
);
`)
	if err != nil {
		return ErrWhitelistSchema
	}
	return nil
}

func ensureBlacklistSchema(ctx context.Context, conn *pg.Conn) error {
	_, err := conn.Exec(ctx, `
CREATE TABLE IF NOT EXISTS token_blacklist (
    	userid BIGINT NOT NULL REFERENCES users(userID) ON DELETE CASCADE,
    	access_token_hash TEXT NOT NULL,
    expired_at TIMESTAMPTZ NOT NULL,
    PRIMARY KEY (userID, access_token_hash)
    );
`)
	if err != nil {
		return ErrBlacklistSchema
	}
	return nil
}

func ensureCategoriesSchema(ctx context.Context, conn *pg.Conn) error {
	_, err := conn.Exec(ctx, `
CREATE TABLE IF NOT EXISTS categories (
    categories_id BIGSERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    parent_id BIGINT REFERENCES categories(categories_id) ON DELETE SET NULL
);`)
	if err != nil {
		return ErrCategoriesSchema
	}
	return nil
}

func ensureProductsSchema(ctx context.Context, conn *pg.Conn) error {
	_, err := conn.Exec(ctx, `
CREATE TABLE IF NOT EXISTS products (
    product_id BIGSERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    description TEXT,
    categories_id BIGINT NOT NULL REFERENCES categories(categories_id),
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    img VARCHAR(255)
);`)
	if err != nil {
		return ErrProductsSchema
	}
	return nil
}

func ensureClientsSchema(ctx context.Context, conn *pg.Conn) error {
	_, err := conn.Exec(ctx, `
CREATE TABLE IF NOT EXISTS clients (
    client_id BIGINT PRIMARY KEY REFERENCES users(userid) ON DELETE CASCADE,
    rating NUMERIC(3,2) DEFAULT 0
);`)
	if err != nil {
		return ErrClientsSchema
	}
	return nil
}

func ensureSellersSchema(ctx context.Context, conn *pg.Conn) error {
	_, err := conn.Exec(ctx, `
CREATE TABLE IF NOT EXISTS sellers (
    seller_id BIGINT PRIMARY KEY REFERENCES users(userid) ON DELETE CASCADE,
    name VARCHAR(255) NOT NULL,
    rating NUMERIC(3,2) DEFAULT 0
);`)
	if err != nil {
		return ErrSellersSchema
	}
	return nil
}

func ensureSellerOffersSchema(ctx context.Context, conn *pg.Conn) error {
	_, err := conn.Exec(ctx, `
CREATE TABLE IF NOT EXISTS seller_offers (
    seller_id BIGINT NOT NULL REFERENCES sellers(seller_id) ON DELETE CASCADE,
    product_id BIGINT NOT NULL REFERENCES products(product_id) ON DELETE CASCADE,
    price BIGINT NOT NULL,
    quantity INT NOT NULL DEFAULT 0,
    PRIMARY KEY (seller_id, product_id)
);`)
	if err != nil {
		return ErrSellerOffersSchema
	}
	return nil
}

func ensureOrdersSchema(ctx context.Context, conn *pg.Conn) error {
	_, err := conn.Exec(ctx, `
CREATE TABLE IF NOT EXISTS orders (
    order_id BIGSERIAL PRIMARY KEY,
    client_id BIGINT NOT NULL REFERENCES clients(client_id),
    status VARCHAR(255) NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);`)
	if err != nil {
		return ErrOrdersSchema
	}
	return nil
}

func ensureOrdersItemsSchema(ctx context.Context, conn *pg.Conn) error {
	_, err := conn.Exec(ctx, `
CREATE TABLE IF NOT EXISTS orders_items (
    order_item_id BIGSERIAL PRIMARY KEY,
    order_id BIGINT NOT NULL REFERENCES orders(order_id) ON DELETE CASCADE,
    seller_id BIGINT NOT NULL,
    product_id BIGINT NOT NULL,
    quantity INT NOT NULL,
    price_at_purchase BIGINT NOT NULL,
    FOREIGN KEY (seller_id, product_id) REFERENCES seller_offers(seller_id, product_id)
);`)
	if err != nil {
		return ErrOrdersItemsSchema
	}
	return nil
}
