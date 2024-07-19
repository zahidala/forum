package types

import (
	"database/sql"
	"sync"
)

type Database struct {
	Conn *sql.DB
	Mu   sync.Mutex
}
type User struct {
	ID       int
	Name     string
	Username string
	Email    string
	Password string
}

type Session struct {
	ID        string
	UserID    int
	Data      string
	CreatedAt string
	ExpiresAt string
}
