package model

import "time"

type User struct {
	ID        int
	Login     string
	Password  string
	Type      string
	Created   time.Time
	LastEnter time.Time
	Balance   float64
}
