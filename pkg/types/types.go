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

type Category struct {
	ID          int
	Name        string
	Description string
	Icon        string
}

type Subcategory struct {
	ID          int
	Name        string
	Description string
	Icon        string
	CategoryID  int
}

type Post struct {
	ID            int
	Title         string
	Content       string
	AuthorID      int
	SubcategoryID int
	CreatedAt     time.Time
}
