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
	ID             int    `json:"id"`
	Name           string `json:"name"`
	Username       string `json:"username"`
	Email          string `json:"email"`
	Password       string `json:"password"`
	ProfilePicture string `json:"profilePicture"`
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
	ID          int      `json:"id"`
	Name        string   `json:"name"`
	Description string   `json:"description"`
	Icon        string   `json:"icon"`
	Category    Category `json:"category"`
	CategoryID  int      `json:"categoryId"`
}

type Post struct {
	ID            int       `json:"id"`
	Title         string    `json:"title"`
	Content       string    `json:"content"`
	AuthorID      int       `json:"authorId"`
	SubcategoryID int       `json:"subcategoryId"`
	CreatedAt     time.Time `json:"createdAt"`
	UpdatedAt     time.Time `json:"updatedAt"`
}

type Comment struct {
	ID        int       `json:"id"`
	PostID    int       `json:"postId"`
	Content   string    `json:"content"`
	AuthorID  int       `json:"authorId"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

type PostLike struct {
	ID        int       `json:"id"`
	PostID    int       `json:"postId"`
	UserID    int       `json:"userId"`
	IsLike    bool      `json:"isLike"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

type CommentLike struct {
	ID        int       `json:"id"`
	CommentID int       `json:"commentId"`
	UserID    int       `json:"userId"`
	IsLike    bool      `json:"isLike"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}
