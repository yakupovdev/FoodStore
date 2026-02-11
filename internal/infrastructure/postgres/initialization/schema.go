package initialization

import (
	"context"
	"errors"
	"log"

	pg "github.com/jackc/pgx/v5"
)

var (
	ErrUsersSchema                  = errors.New("users schema error")
	ErrRecoveryCodesSchema          = errors.New("recovery codes schema error")
	ErrWhitelistSchema              = errors.New("whitelist schema error")
	ErrBlacklistSchema              = errors.New("blacklist schema error")
	ErrCategoriesSchema             = errors.New("categories schema error")
	ErrProductsSchema               = errors.New("products schema error")
	ErrClientsSchema                = errors.New("clients schema error")
	ErrSellersSchema                = errors.New("sellers schema error")
	ErrSellerOffersSchema           = errors.New("seller offers schema error")
	ErrOrdersSchema                 = errors.New("orders schema error")
	ErrOrdersItemsSchema            = errors.New("orders items schema error")
	ErrModerationSellerOffersSchema = errors.New("moderation seller offers schema error")
	ErrCategoriesAdd                = errors.New("categories add error")
	ErrProductsAdd                  = errors.New("products add error")
	ErrCartSchema                   = errors.New("cart schema error")
	ErrCartItemsSchema              = errors.New("cart items schema error")
	ErrLogsSchema                   = errors.New("logs schema error")
	ErrSubscriptionSchema           = errors.New("subscription schema error")
)

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
	if err := ensureModerationSellerOffersSchema(ctx, conn); err != nil {
		return err
	}
	if err := ensureCartSchema(ctx, conn); err != nil {
		return err
	}
	if err := ensureCartItemsSchema(ctx, conn); err != nil {
		return err
	}
	if err := ensureCategoriesForTest(ctx, conn); err != nil {
		return err
	}
	if err := ensureProductsForTest(ctx, conn); err != nil {
		return err
	}
	if err := ensureLogsTransactionSchema(ctx, conn); err != nil {
		return err
	}
	if err := ensureSubscriptionSchema(ctx, conn); err != nil {
		return err
	}
	//if err := ensureCategoriesForTest(ctx, conn); err != nil {
	//	return err
	//}
	//if err := ensureProductsForTest(ctx, conn); err != nil {
	//	return err
	//}

	return nil
}

func ensureUsersSchema(ctx context.Context, conn *pg.Conn) error {
	_, err := conn.Exec(ctx, `
CREATE TABLE IF NOT EXISTS users (
	userid BIGSERIAL PRIMARY KEY,
	email VARCHAR(100) NOT NULL,
	password_hash VARCHAR(255) NOT NULL,
	type VARCHAR(10) NOT NULL,
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
	userid BIGINT NOT NULL REFERENCES users(userid) ON DELETE CASCADE,
	email VARCHAR(100) NOT NULL,
	code_hash VARCHAR(255) NOT NULL,
    type VARCHAR(10) NOT NULL,
	expired_at TIMESTAMPTZ NOT NULL,
	created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
	PRIMARY KEY (userid)
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
	userid BIGINT NOT NULL REFERENCES users(userid) ON DELETE CASCADE,
	access_token_hash VARCHAR(255) NOT NULL,
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
    	userid BIGINT NOT NULL REFERENCES users(userid) ON DELETE CASCADE,
    	access_token_hash VARCHAR(255) NOT NULL,
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
    name VARCHAR(255) NOT NULL,
    address VARCHAR(255),
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
    address VARCHAR(255),
    priority INT DEFAULT 0,
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
    client_id BIGINT NOT NULL REFERENCES clients(client_id) ON DELETE CASCADE,
    status VARCHAR(100) NOT NULL,
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
    seller_id BIGINT NOT NULL ,
    product_id BIGINT NOT NULL,
    quantity INT NOT NULL,
    price_at_purchase BIGINT NOT NULL,
    FOREIGN KEY (seller_id, product_id) REFERENCES seller_offers(seller_id, product_id) ON DELETE CASCADE
);`)
	if err != nil {
		log.Println(err)
		return ErrOrdersItemsSchema
	}
	return nil
}

func ensureModerationSellerOffersSchema(ctx context.Context, conn *pg.Conn) error {
	_, err := conn.Exec(ctx, `
CREATE TABLE IF NOT EXISTS moderation_seller_offers (
    seller_id BIGINT NOT NULL REFERENCES sellers(seller_id) ON DELETE CASCADE,
    seller_name VARCHAR(255) NOT NULL,
    seller_email VARCHAR(255) NOT NULL,
    category_id BIGINT NOT NULL REFERENCES categories(categories_id) ON DELETE CASCADE,
    category_name VARCHAR(255) NOT NULL,
    subcategory_id BIGINT NOT NULL REFERENCES categories(categories_id) ON DELETE CASCADE,
    subcategory_name VARCHAR(255) NOT NULL,
    product_id BIGSERIAL PRIMARY KEY,
    product_name VARCHAR(255) NOT NULL,
    description TEXT NOT NULL,
    image VARCHAR(255) NOT NULL,
    price BIGINT NOT NULL,
    quantity INT NOT NULL DEFAULT 0
);`)
	if err != nil {
		log.Println(err)
		return ErrModerationSellerOffersSchema
	}
	return nil
}

func ensureCategoriesForTest(ctx context.Context, conn *pg.Conn) error {
	stmt := []string{
		`INSERT INTO categories (name, parent_id) VALUES('Milk,cheese and eggs',NULL);`,
		`INSERT INTO categories (name, parent_id) VALUES('Meat and poultry',NULL);`,
		`INSERT INTO categories (name, parent_id)VALUES('Fish and products',NULL);`,

		`INSERT INTO categories (name, parent_id) VALUES('Milk, cream, condensed milk',1);`,
		`INSERT INTO categories (name, parent_id) VALUES('Kefir, cottage cheese, sour cream',1);`,
		`INSERT INTO categories (name, parent_id) VALUES('Yogurts, cottage cheese and desserts',1);`,
		`INSERT INTO categories (name, parent_id) VALUES('Eggs, butter, margarine',1);`,
		`INSERT INTO categories (name, parent_id) VALUES('Cheese',1);`,

		`INSERT INTO categories (name, parent_id) VALUES('Meat, steaks, minced meat',2);`,
		`INSERT INTO categories (name, parent_id) VALUES('Chicken, turkey, and poultry',2);`,
		`INSERT INTO categories (name, parent_id) VALUES('Semi-finished products and marinades',2);`,

		`INSERT INTO categories (name, parent_id) VALUES('Fish',3);`,
		`INSERT INTO categories (name, parent_id) VALUES('Seafood',3);`,
		`INSERT INTO categories (name, parent_id) VALUES('Caviar and snacks',3);`,
	}

	for _, s := range stmt {
		if _, err := conn.Exec(ctx, s); err != nil {
			return ErrCategoriesAdd
		}
	}

	return nil
}

func ensureProductsForTest(ctx context.Context, conn *pg.Conn) error {
	stmt := []string{
		`INSERT INTO products (name, description, categories_id, created_at, img) 
		 VALUES('Whole Milk', 'Fresh whole milk 1L', 1, NOW(), 'milk.jpg');`,
		`INSERT INTO products (name, description, categories_id, created_at, img) 
		 VALUES('Butter', 'Salted butter 200g', 1, NOW(), 'butter.jpg');`,

		`INSERT INTO products (name, description, categories_id, created_at, img) 
		 VALUES('Condensed Milk', 'Sweetened condensed milk 400g', 4, NOW(), 'condensed_milk.jpg');`,
		`INSERT INTO products (name, description, categories_id, created_at, img) 
		 VALUES('Cream 20%', 'Fresh cream 250ml', 4, NOW(), 'cream.jpg');`,

		`INSERT INTO products (name, description, categories_id, created_at, img) 
		 VALUES('Beef Steak', 'Premium beef steak 300g', 10, NOW(), 'beef_steak.jpg');`,
		`INSERT INTO products (name, description, categories_id, created_at, img) 
		 VALUES('Ground Beef', 'Minced beef 500g', 10, NOW(), 'ground_beef.jpg');`,

		`INSERT INTO products (name, description, categories_id, created_at, img) 
		 VALUES('Shrimps', 'Frozen shrimps 500g', 12, NOW(), 'shrimps.jpg');`,
		`INSERT INTO products (name, description, categories_id, created_at, img) 
		 VALUES('Salmon Fillet', 'Fresh salmon fillet 400g', 12, NOW(), 'salmon.jpg');`,
	}

	for _, s := range stmt {
		if _, err := conn.Exec(ctx, s); err != nil {
			return ErrProductsAdd
		}
	}

	return nil
}

func ensureCartSchema(ctx context.Context, conn *pg.Conn) error {
	_, err := conn.Exec(ctx, `
CREATE TABLE IF NOT EXISTS cart (
    cart_id BIGSERIAL PRIMARY KEY,
    client_id BIGINT NOT NULL REFERENCES clients(client_id) ON DELETE CASCADE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);`)
	if err != nil {
		return ErrCartSchema
	}
	return nil
}

func ensureCartItemsSchema(ctx context.Context, conn *pg.Conn) error {
	_, err := conn.Exec(ctx, `
CREATE TABLE IF NOT EXISTS cart_items (
    cart_item_id BIGSERIAL PRIMARY KEY,
    cart_id BIGINT NOT NULL REFERENCES cart(cart_id) ON DELETE CASCADE,
    seller_id BIGINT NOT NULL ,
    product_id BIGINT NOT NULL,
    quantity INT NOT NULL,
    price_at_purchase BIGINT NOT NULL,
    FOREIGN KEY (seller_id, product_id) REFERENCES seller_offers(seller_id, product_id) ON DELETE CASCADE
);`)
	if err != nil {
		return ErrCartItemsSchema
	}
	return nil
}

func ensureLogsTransactionSchema(ctx context.Context, conn *pg.Conn) error {
	_, err := conn.Exec(ctx, `
CREATE TABLE IF NOT EXISTS logs_transactions (
	log_id BIGSERIAL PRIMARY KEY,
	client_id BIGINT NOT NULL REFERENCES clients(client_id) ON DELETE CASCADE,
    seller_id BIGINT NOT NULL REFERENCES sellers(seller_id) ON DELETE CASCADE,
    total_amount BIGINT NOT NULL,
    commission_amount BIGINT NOT NULL,
	created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);`)
	if err != nil {
		return ErrLogsSchema
	}
	return nil
}

func ensureSubscriptionSchema(ctx context.Context, conn *pg.Conn) error {
	_, err := conn.Exec(ctx, `
CREATE TABLE IF NOT EXISTS subscriptions (
	subscription_id BIGSERIAL PRIMARY KEY,
	seller_id BIGINT NOT NULL REFERENCES sellers(seller_id) ON DELETE CASCADE,
	created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);`)
	if err != nil {
		return ErrSubscriptionSchema
	}
	return nil
}
