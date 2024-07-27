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
	ID          int    `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Icon        string `json:"icon"`
}

type Subcategory struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Icon        string `json:"icon"`
	CategoryID  int    `json:"categoryId"`
}

type Post struct {
	ID            int
	Title         string
	Content       string
	AuthorID      int
	SubcategoryID int
	CreatedAt     time.Time
}
