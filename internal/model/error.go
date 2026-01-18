package model

const (
	ErrNoRecord                 = "no matching record found"
	ErrInvalidCredentials       = "invalid credentials"
	ErrDuplicateLogin           = "duplicate login"
	DatabaseConnectionError     = "failed to connect to the database"
	DatabaseQueryError          = "database query error"
	DatabaseNotInitializedError = "database not initialized"
	CouldNotGenerateTokenError  = "could not generate token"
)
