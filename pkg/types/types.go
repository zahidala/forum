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
	ID            int       `json:"id"`
	Title         string    `json:"title"`
	Content       string    `json:"content"`
	AuthorID      int       `json:"authorId"`
	SubcategoryID int       `json:"subcategoryId"`
	CreatedAt     time.Time `json:"createdAt"`
}

type Comment struct {
	ID        int       `json:"id"`
	PostID    int       `json:"postId"`
	Content   string    `json:"content"`
	AuthorID  int       `json:"authorId"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}
