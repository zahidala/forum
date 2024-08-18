package users

import (
	"database/sql"
	"forum/pkg/db"
	templates "forum/pkg/templates"
	types "forum/pkg/types"
	"forum/pkg/utils"
	"log"
	"net/http"
	"time"
)

func GetUserByIDHandler(w http.ResponseWriter, r *http.Request) {
	// Implementation specific to fetching users

	// cookie, cookieErr := r.Cookie("session_id")

	// if cookieErr == nil {
	// 	http.Error(w, "Unauthorized", http.StatusUnauthorized)
	// 	return
	// }

	// sessionId := cookie.Value
}

func CreateUserHandler(w http.ResponseWriter, r *http.Request) {
	username := r.FormValue("username")
	name := r.FormValue("name")
	email := r.FormValue("email")
	password := r.FormValue("password")

	isValid, err := RegValidation(w, r)
	if !isValid || err != nil {
		return
	}

	hashedPassword, hashErr := utils.HashPassword(password)
	if hashErr != nil {
		log.Println(hashErr)
		http.Error(w, "Error generating password hash", http.StatusInternalServerError)
		return
	}

	query := "INSERT INTO users (name, username, email, password) VALUES (?, ?, ?, ?)"

	userAddExecErr := db.PrepareAndExecute(query,
		name,
		username,
		email,
		hashedPassword,
	)
	if userAddExecErr != nil {
		http.Error(w, "Error creating user", http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/login", http.StatusSeeOther)
}

func RegValidation(w http.ResponseWriter, r *http.Request) (bool, error) {
	isValid := true

	username := r.FormValue("username")
	email := r.FormValue("email")

	userQuery := "SELECT EXISTS(SELECT 1 FROM users WHERE username = ?)"

	userStmt, userErr := db.GetDB().Prepare(userQuery)
	if userErr != nil {
		http.Error(w, "Error preparing query", http.StatusInternalServerError)
		return isValid, userErr
	}
	defer userStmt.Close()

	var userExists bool
	findUserErr := userStmt.QueryRow(username).Scan(&userExists)
	if findUserErr != nil {
		log.Println(findUserErr)
		http.Error(w, "Error querying database", http.StatusInternalServerError)
		return isValid, findUserErr
	}

	if userExists {
		// return user already exists message
		data := types.ErrorPageProps{
			Error: types.Error{Message: "User already exists"},
			Title: "username",
		}
		w.WriteHeader(http.StatusConflict)
		templates.RegisterTemplateHandler(w, r, data)
		return false, nil
	}

	emailQuery := "SELECT EXISTS(SELECT 1 FROM users WHERE email = ?)"
	emailStmt, emailErr := db.GetDB().Prepare(emailQuery)
	if emailErr != nil {
		http.Error(w, "Error preparing query", http.StatusInternalServerError)
		return isValid, emailErr
	}
	defer emailStmt.Close()

	var emailExists bool
	findEmailErr := emailStmt.QueryRow(email).Scan(&emailExists)
	if findEmailErr != nil {
		log.Println(findEmailErr)
		http.Error(w, "Error querying database", http.StatusInternalServerError)
		return isValid, findEmailErr
	}

	if emailExists {
		// return account already exists message
		data := types.ErrorPageProps{
			Error: types.Error{Message: "Account already exists! Please log-in"},
			Title: "email",
		}
		w.WriteHeader(http.StatusConflict)
		templates.RegisterTemplateHandler(w, r, data)
		return false, nil
	}
	return isValid, nil
}

func UserLoginHandler(w http.ResponseWriter, r *http.Request) {
	username := r.FormValue("username")
	password := r.FormValue("password")

	userQuery := "SELECT id, password FROM users WHERE username = ?"

	userStmt, userErr := db.GetDB().Prepare(userQuery)
	if userErr != nil {
		http.Error(w, "Error preparing query", http.StatusInternalServerError)
		return
	}
	defer userStmt.Close()

	var userId int
	var hashedPassword string

	findUserErr := userStmt.QueryRow(username).Scan(&userId, &hashedPassword)
	switch {
	case findUserErr == sql.ErrNoRows:
		data := types.Error{Message: "User not found"}
		w.WriteHeader(http.StatusNotFound)
		templates.LoginTemplateHandler(w, r, data)
		// http.Error(w, "User not found", http.StatusNotFound)
		return
	case findUserErr != nil:
		http.Error(w, "Error querying database", http.StatusInternalServerError)
		return
	}

	compareErr := utils.CompareHashAndPassword(hashedPassword, password)
	if compareErr != nil {
		data := types.Error{Message: "Incorrect username/password"}
		w.WriteHeader(http.StatusUnauthorized)
		templates.LoginTemplateHandler(w, r, data)
		// http.Error(w, "Invalid password", http.StatusUnauthorized)
		return
	}

	sessionId := utils.GenerateSessionID()

	createdAt := time.Now()
	expiresAt := createdAt.Add(24 * time.Hour)

	http.SetCookie(w, &http.Cookie{
		Name:    "sessionId",
		Value:   sessionId,
		Expires: expiresAt,
	})

	sessionsAddQuery := "INSERT INTO sessions (id, userId, createdAt, expiresAt) VALUES (?, ?, ?, ?)"

	sessionAddExecErr := db.PrepareAndExecute(sessionsAddQuery, sessionId, userId, createdAt, expiresAt)
	if sessionAddExecErr != nil {
		http.Error(w, "Error creating session", http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func UserLogoutHandler(w http.ResponseWriter, r *http.Request) {
	cookie, cookieErr := r.Cookie("sessionId")
	if cookieErr != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	sessionId := cookie.Value

	sessionDeleteQuery := "DELETE FROM sessions WHERE id = ?"

	sessionDeleteExecErr := db.PrepareAndExecute(sessionDeleteQuery, sessionId)
	if sessionDeleteExecErr != nil {
		http.Error(w, "Error deleting session", http.StatusInternalServerError)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:    "sessionId",
		Value:   "",
		Expires: time.Now(),
	})

	http.Redirect(w, r, "/", http.StatusSeeOther)
}
