package utils

import (
	"fmt"
	"forum/pkg/db"
	Types "forum/pkg/types"
	"github.com/gofrs/uuid/v5"
	"golang.org/x/crypto/bcrypt"
	"log"
	"net/http"
	"time"
	// "strings"
)

// GenerateSessionID generates a new UUID for session ID
func GenerateSessionID() string {
	uuid, _ := uuid.NewV4()
	return uuid.String()
}

// HashPassword generates a bcrypt hash of the password
func HashPassword(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)

	if err != nil {
		return "", err
	}

	return string(hash), nil
}

// CompareHashAndPassword compares a bcrypt hashed password with its possible plaintext equivalent.
// Returns nil on success, or an error on failure.
func CompareHashAndPassword(hashedPassword, password string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
}

// IsAuthenticated checks if the user is logged in by checking if the session ID cookie exists in the database
// Returns true if the user is logged in, false otherwise. If an error occurs, it returns false. Also checks if the session has expired.
// Used for rendering buttons based on the user's authentication status.
func IsAuthenticated(r *http.Request) bool {
	cookie, err := r.Cookie("sessionId")
	if err != nil {
		return false
	}

	sessionsQuery := "SELECT expiresAt FROM sessions WHERE id = ?"
	var expiresAtSession time.Time

	sessionStmt, sessionErr := db.GetDB().Prepare(sessionsQuery)
	if sessionErr != nil {
		log.Println(sessionErr)
		return false
	}
	defer sessionStmt.Close()

	sessionRowErr := sessionStmt.QueryRow(cookie.Value).Scan(&expiresAtSession)

	if sessionRowErr != nil {
		log.Println(err)
		return false
	}

	return !time.Now().After(expiresAtSession)
}

// Returns the user object based on the session ID cookie
func GetUserInfoBySession(w http.ResponseWriter, r *http.Request) Types.User {
	cookie, cookieErr := r.Cookie("sessionId")

	if cookieErr != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return Types.User{}
	}

	sessionId := cookie.Value

	sessionsQuery := "SELECT userId FROM sessions WHERE id = ?"

	sessionStmt, sessionErr := db.GetDB().Prepare(sessionsQuery)
	if sessionErr != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		log.Println(sessionErr)
		return Types.User{}
	}

	var userId int

	sessionRowErr := sessionStmt.QueryRow(sessionId).Scan(&userId)

	if sessionRowErr != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return Types.User{}
	}

	userQuery := "SELECT * FROM users WHERE id = ?"

	userStmt, userErr := db.GetDB().Prepare(userQuery)

	if userErr != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		log.Println(userErr)
		return Types.User{}
	}

	var user Types.User

	userRowErr := userStmt.QueryRow(userId).Scan(&user.ID, &user.Name, &user.Username, &user.Email, &user.Password, &user.ProfilePicture)

	if userRowErr != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		log.Println(userRowErr)
		return Types.User{}
	}

	return user
}

func GetFilteredPosts(w http.ResponseWriter, r *http.Request) string {
	//protect? against user-input (manually) queries

	params := r.URL.Query()
	categories := params["category"]
	userPosts := params.Get("user-posts")
	likedPosts := params.Get("liked-posts")

	if len(categories) < 1 && userPosts != "true" && likedPosts != "true" {
		return ""
	}

	// fmt.Println("Categories:", categories)
	// fmt.Println("userPosts:", userPosts)
	// fmt.Println("likedPosts:", likedPosts)

	filters := "\nWHERE "

	if len(categories) > 0 {
		for i, category := range categories {
			if i > 0 {
				filters += " OR "
			}

			filters += fmt.Sprintf(`categories LIKE '%%"categoryID":%s%%'`, category)
		}
		
		filters += "\n"
	}


	var userID int
	if IsAuthenticated(r) {
		userID = GetUserInfoBySession(w, r).ID

		if userPosts == "true" {
			if len(filters) > 10 {
				filters += "AND "
			}
			filters += fmt.Sprintf("userID = %d\n", userID)
		}
	
		if likedPosts == "true" {
			if len(filters) > 10 {
				filters += "AND "
			}
			filters += fmt.Sprintf("postID IN (SELECT postId FROM PostLikes pl WHERE userId = %d AND isLike = 1)\n", userID)
		}
	}
	// if no user cookie, ignore any query other than categories

	return filters
}
