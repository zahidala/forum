package types

import (
	"database/sql"
	"html/template"
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

type RegValidation struct {
	Name     string
	Username string
	Email    string
	Password string
	Errors   map[string]string
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
	Attachments   string    `json:"attachments"`
}

type Comment struct {
	ID          int           `json:"id"`
	PostID      int           `json:"postId"`
	Content     template.HTML `json:"content"`
	AuthorID    int           `json:"authorId"`
	CreatedAt   time.Time     `json:"createdAt"`
	UpdatedAt   time.Time     `json:"updatedAt"`
	Attachments string        `json:"attachments"`
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

type PostDisLike struct {
	ID        int       `json:"id"`
	PostID    int       `json:"postId"`
	UserID    int       `json:"userId"`
	IsDislike bool      `json:"isDislike"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

type CommentDisLike struct {
	ID        int       `json:"id"`
	CommentID int       `json:"commentId"`
	UserID    int       `json:"userId"`
	IsDislike bool      `json:"isDislike"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

// Image struct
type UploadedImage struct {
	StatusCode int     `json:"status_code"`
	Success    Success `json:"success"`
	Image      Image   `json:"image"`
	StatusTxt  string  `json:"status_txt"`
}

type Success struct {
	Message string `json:"message"`
	Code    int    `json:"code"`
}

type Image struct {
	Name               string  `json:"name"`
	Extension          string  `json:"extension"`
	Width              int     `json:"width"`
	Height             int     `json:"height"`
	Size               int     `json:"size"`
	Time               int64   `json:"time"`
	Expiration         int     `json:"expiration"`
	Likes              int     `json:"likes"`
	Description        *string `json:"description"`
	OriginalFilename   string  `json:"original_filename"`
	IsAnimated         int     `json:"is_animated"`
	NSFW               int     `json:"nsfw"`
	IDEncoded          string  `json:"id_encoded"`
	SizeFormatted      string  `json:"size_formatted"`
	Filename           string  `json:"filename"`
	URL                string  `json:"url"`
	URLShort           string  `json:"url_short"`
	URLSeo             string  `json:"url_seo"`
	URLViewer          string  `json:"url_viewer"`
	URLViewerPreview   string  `json:"url_viewer_preview"`
	URLViewerThumb     string  `json:"url_viewer_thumb"`
	ImageDetails       Details `json:"image"`
	Thumb              Details `json:"thumb"`
	Medium             Details `json:"medium"`
	DisplayURL         string  `json:"display_url"`
	DisplayWidth       int     `json:"display_width"`
	DisplayHeight      int     `json:"display_height"`
	ViewsLabel         string  `json:"views_label"`
	LikesLabel         string  `json:"likes_label"`
	HowLongAgo         string  `json:"how_long_ago"`
	DateFixedPeer      string  `json:"date_fixed_peer"`
	Title              string  `json:"title"`
	TitleTruncated     string  `json:"title_truncated"`
	TitleTruncatedHTML string  `json:"title_truncated_html"`
	IsUseLoader        bool    `json:"is_use_loader"`
}

type Details struct {
	Filename  string `json:"filename"`
	Name      string `json:"name"`
	Mime      string `json:"mime"`
	Extension string `json:"extension"`
	URL       string `json:"url"`
	Size      int    `json:"size,omitempty"`
}
