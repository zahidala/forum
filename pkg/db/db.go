package db

import (
	"database/sql"
	Types "forum/pkg/types"
	"log"
	"os"
	"sync"

	_ "github.com/mattn/go-sqlite3"
)

var instance *Types.Database
var once sync.Once

// Init initializes the database connection
func Init() {
	once.Do(func() {
		dbFile := "./forum-v2.db"
		_, err := os.Stat(dbFile)
		dbExists := !os.IsNotExist(err)

		conn, err := sql.Open("sqlite3", "./forum-v2.db")
		if err != nil {
			log.Fatalf("Error opening the database: %s", err)
			return
		}

		if err := conn.Ping(); err != nil {
			log.Fatalf("Error connecting to the database: %s", err)
			return
		}

		log.Println("Connected to the database")

		instance = &Types.Database{
			Conn: conn,
		}

		if !dbExists {
			log.Println("Database file does not exist. Creating a new file and seeding database...")
			seedDB()
		}
	})
}

// seedDB seeds the newly created database file with initial data if it does not exist
func seedDB() {
	createTables := []string{
		`CREATE TABLE Users ( 
			id INTEGER PRIMARY KEY,    
			name TEXT NOT NULL,   
			username TEXT UNIQUE NOT NULL,    
			email TEXT UNIQUE NOT NULL,
			password TEXT NOT NULL,
			profilePicture TEXT
		);`,

		`CREATE TABLE Sessions (
			id TEXT PRIMARY KEY,
			userId INTEGER NOT NULL,
			data TEXT,
			createdAt DATETIME NOT NULL,
			expiresAt DATETIME NOT NULL,
			FOREIGN KEY (userId) REFERENCES Users(id)
		);`,

		`CREATE TABLE Categories (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			name TEXT NOT NULL,
			description TEXT,
			icon TEXT
		);`,

		`CREATE TABLE Posts (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			title TEXT NOT NULL,
			content TEXT NOT NULL,
			authorId INTEGER NOT NULL,
			createdAt DATETIME DEFAULT CURRENT_TIMESTAMP,
			updatedAt DATETIME DEFAULT CURRENT_TIMESTAMP,
			attachments TEXT
		);`,

		`CREATE TABLE PostLikes (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			postId INTEGER NOT NULL,
			userId INTEGER NOT NULL,
			isLike INTEGER NOT NULL CHECK (isLike IN (0, 1)),
			createdAt DATETIME DEFAULT CURRENT_TIMESTAMP,
			updatedAt DATETIME DEFAULT CURRENT_TIMESTAMP,
			FOREIGN KEY (postId) REFERENCES posts(id) ON DELETE CASCADE,
			FOREIGN KEY (userId) REFERENCES users(id) ON DELETE CASCADE
		);`,

		`CREATE TABLE PostDislikes (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			postId INTEGER NOT NULL,
			userId INTEGER NOT NULL,
			isDisLike INTEGER NOT NULL CHECK (isDisLike IN (0, 1)),
			createdAt DATETIME DEFAULT CURRENT_TIMESTAMP,
			updatedAt DATETIME DEFAULT CURRENT_TIMESTAMP,
			FOREIGN KEY (postId) REFERENCES posts(id) ON DELETE CASCADE,
			FOREIGN KEY (userId) REFERENCES users(id) ON DELETE CASCADE
		);`,

		`CREATE TABLE Comments (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			postId INTEGER NOT NULL,
			content TEXT NOT NULL,
			authorId INTEGER NOT NULL,
			createdAt DATETIME DEFAULT CURRENT_TIMESTAMP,
			updatedAt DATETIME DEFAULT CURRENT_TIMESTAMP,
			attachments TEXT,
			FOREIGN KEY (postId) REFERENCES posts(id) ON DELETE CASCADE,
			FOREIGN KEY (authorId) REFERENCES users(id) ON DELETE CASCADE
		);`,

		`CREATE TABLE CommentLikes (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			commentId INTEGER NOT NULL,
			userId INTEGER NOT NULL,
			isLike INTEGER NOT NULL CHECK (isLike IN (0, 1)),
			createdAt DATETIME DEFAULT CURRENT_TIMESTAMP,
			updatedAt DATETIME DEFAULT CURRENT_TIMESTAMP,
			FOREIGN KEY (commentId) REFERENCES comments(id) ON DELETE CASCADE,
			FOREIGN KEY (userId) REFERENCES users(id) ON DELETE CASCADE
		);`,

		`CREATE TABLE CommentDislikes (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			commentId INTEGER NOT NULL,
			userId INTEGER NOT NULL,
			isDislike INTEGER NOT NULL CHECK (isDislike IN (0, 1)),
			createdAt DATETIME DEFAULT CURRENT_TIMESTAMP,
			updatedAt DATETIME DEFAULT CURRENT_TIMESTAMP,
			FOREIGN KEY (commentId) REFERENCES comments(id) ON DELETE CASCADE,
			FOREIGN KEY (userId) REFERENCES users(id) ON DELETE CASCADE
		);`,

		`CREATE TABLE PostCategories (
			postId INTEGER NOT NULL,
			categoryId INTEGER NOT NULL,
			PRIMARY KEY (postId, categoryId),
			FOREIGN KEY (postId) REFERENCES posts(id) ON DELETE CASCADE,
			FOREIGN KEY (categoryId) REFERENCES categories(id) ON DELETE CASCADE
		);`,
	}

	// Create tables
	for _, query := range createTables {
		if err := PrepareAndExecute(query); err != nil {
			log.Println(query)
			log.Fatalf("Error creating table: %s", err)
			return
		}
	}

	initialData := []string{
		`INSERT INTO Users (name, username, email, password, profilePicture) VALUES ('John Doe', 'johndoe', 'johndoe@gmail.com', '$2a$10$M9APgO1pJZgsfMdj9SmZEORF94WYnS5RkXrIaVA7ZG6bXgzSB5lEa', 'https://iili.io/dW44kLG.jpg');`,
		`INSERT INTO Users (name, username, email, password, profilePicture) VALUES ('Jane Doe', 'janedoe', 'janedoe@gmail.com', '$2a$10$M9APgO1pJZgsfMdj9SmZEORF94WYnS5RkXrIaVA7ZG6bXgzSB5lEa', 'https://iili.io/dW44kLG.jpg');`,

		`INSERT INTO Categories (name, description, icon) VALUES ('General', 'Talk about any car here', NULL);`,
		`INSERT INTO Categories (name, description, icon) VALUES ('Introductions', 'Drop in and introduce yourself to the community!', NULL);`,
		`INSERT INTO Categories (name, description, icon) VALUES ('Off-Topic', 'Anything but car talk here.', NULL);`,
		`INSERT INTO Categories (name, description, icon) VALUES ('General Question and Answers', 'This can be used for any question you may have. Also ask questions for manufacturers or models that don''t have a forum.', NULL);`,

		`INSERT INTO Posts (title, content, authorId, attachments) VALUES ('Hello World', 'This is the first post on the forum!', 1, 'https://iili.io/dW44kLG.jpg');`,
		`INSERT INTO Posts (title, content, authorId) VALUES ('Hello World 2', 'This is the second post on the forum!', 1);`,
		`INSERT INTO PostCategories (postId, categoryId) VALUES (1, 1);`,
		`INSERT INTO PostCategories (postId, categoryId) VALUES (1, 2);`,
		`INSERT INTO PostCategories (postId, categoryId) VALUES (2, 1);`,
		`INSERT INTO PostCategories (postId, categoryId) VALUES (2, 3);`,

		`INSERT INTO Comments (postId, content, authorId) VALUES (1, 'This is the first comment on the forum!', 1);`,
		`INSERT INTO Comments (postId, content, authorId) VALUES (2, 'This is the second comment on the forum!', 2);`,

		`INSERT INTO PostLikes (postId, userId, isLike) VALUES (1, 2, 1);`,
		`INSERT INTO PostDislikes (postId, userId, isDislike) VALUES (2, 2, 1);`,

		`INSERT INTO CommentLikes (commentId, userId, isLike) VALUES (1, 2, 1);`,
		`INSERT INTO CommentDislikes (commentId, userId, isDislike) VALUES (2, 1, 1);`,
	}

	// Insert initial data
	for _, query := range initialData {
		if err := PrepareAndExecute(query); err != nil {
			log.Println(query)
			log.Fatalf("Error inserting initial data: %s", err)
			return
		}
	}

	log.Println("Database seeded successfully")
}

// GetDB returns the database instance
func GetDB() *sql.DB {
	if instance == nil {
		log.Fatal("Database not initialized. Call Init() first.")
	}
	return instance.Conn
}

// CloseDB closes the database connection
func CloseDB() {
	if instance != nil {
		instance.Mu.Lock()
		defer instance.Mu.Unlock()
		if err := instance.Conn.Close(); err != nil {
			log.Printf("Error closing the database: %s", err)
		}
	}
}

// PrepareAndExecute prepares and executes a query. It returns an error if the query fails.
// May be expanded to return the result of the query in the future.
func PrepareAndExecute(query string, args ...interface{}) error {
	stmt, stmtErr := GetDB().Prepare(query)
	if stmtErr != nil {
		log.Printf("Error preparing query: %s", stmtErr)
		return stmtErr
	}

	defer stmt.Close()

	_, execErr := stmt.Exec(args...)
	if execErr != nil {
		log.Printf("Error executing query: %s", execErr)
		return execErr
	}

	return nil
}
