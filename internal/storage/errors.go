package storage

import "errors"

var (
	ErrDatabaseConnection  = errors.New("failed to connect to the database")
	ErrUsersSchema         = errors.New("users schema error")
	ErrRecoveryCodesSchema = errors.New("recovery codes schema error")
	ErrWhitelistSchema     = errors.New("whitelist schema error")
	ErrBlacklistSchema     = errors.New("blacklist schema error")
	ErrCategoriesSchema    = errors.New("categories schema error")
	ErrCategoriesAdd       = errors.New("failed to add category")
	ErrProductsSchema      = errors.New("products schema error")
	ErrClientsSchema       = errors.New("clients schema error")
	ErrSellersSchema       = errors.New("sellers schema error")
	ErrSellerOffersSchema  = errors.New("seller offers schema error")
	ErrOrdersSchema        = errors.New("orders schema error")
	ErrOrdersItemsSchema   = errors.New("orders items schema error")
)
