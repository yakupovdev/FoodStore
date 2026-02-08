package entity

import "time"

type RecoveryCode struct {
	UserID    int64
	Email     string
	UserType  string
	CodeHash  string
	ExpiredAt time.Time
}

func (rc *RecoveryCode) IsExpired() bool {
	return time.Now().After(rc.ExpiredAt)
}

func (rc *RecoveryCode) MatchesHash(hash string) bool {
	return rc.CodeHash == hash
}

func (rc *RecoveryCode) IsValid(hash string) bool {
	return !rc.IsExpired() && rc.MatchesHash(hash)
}
