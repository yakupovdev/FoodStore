package security

type CheckerAdminKey struct {
	SecretKey string
}

func NewCheckerAdminKey(SecretKey string) *CheckerAdminKey {
	return &CheckerAdminKey{
		SecretKey: SecretKey,
	}
}

func (c *CheckerAdminKey) CheckAdminKey(key string) bool {
	ok := key == c.SecretKey
	return ok
}
