package types

import (
	"database/sql"
	"sync"
	"time"
)

type Database struct {
	Conn *sql.DB
	Mu   sync.Mutex
}

type Error struct {
	Message string
	Code    int
}
type ErrorPageProps struct {
	Error Error
	Title string
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
	CreatedAt time.Time
	ExpiresAt time.Time
}
