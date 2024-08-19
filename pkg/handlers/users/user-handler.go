package users

import (
	"database/sql"
	"forum/pkg/db"
	templates "forum/pkg/templates"
	types "forum/pkg/types"
	"forum/pkg/utils"
	"log"
	"net/http"
	"regexp"
	"strings"
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
	name := r.FormValue("name")
	username := r.FormValue("username")
	email := r.FormValue("email")
	password := r.FormValue("password")

	data := types.RegValidation{}

	err := RegValidation(w, r, &data)
	if err != nil {
		return
	}

	// log.Println(data)

	if len(data.Errors) != 0 {
		w.WriteHeader(http.StatusConflict)
		templates.RegisterTemplateHandler(w, r, data)
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

func RegValidation(w http.ResponseWriter, r *http.Request, data *types.RegValidation) error {

	name := r.FormValue("name")
	username := strings.TrimSpace(r.FormValue("username"))
	email := strings.TrimSpace(r.FormValue("email"))
	password := r.FormValue("password")

	data.Errors = make(map[string]string)

	// name validation
	fullName := strings.Join(strings.Fields(name), " ")

	re := regexp.MustCompile(`^[a-zA-Z ]{3,50}$`)
	matched := re.MatchString(fullName)
	if !matched {
		data.Errors["Name"] = "Please enter a valid full name"
	} else {
		data.Name = fullName
	}

	// username validation
	re = regexp.MustCompile(`^[a-zA-Z\d_]{3,20}$`)
	matched = re.MatchString(username)
	if !matched || strings.Contains(username, " ") {
		data.Errors["Username"] = "Username must be 3-20 characters long and only caontain alphabets, numbers, and/or _"
	} else {
		data.Username = username
	}

	userQuery := "SELECT EXISTS(SELECT 1 FROM users WHERE username = ?)"
	userStmt, userErr := db.GetDB().Prepare(userQuery)
	if userErr != nil {
		http.Error(w, "Error preparing query", http.StatusInternalServerError)
		return userErr
	}
	defer userStmt.Close()

	var userExists bool
	findUserErr := userStmt.QueryRow(username).Scan(&userExists)
	if findUserErr != nil {
		log.Println(findUserErr)
		http.Error(w, "Error querying database", http.StatusInternalServerError)
		return findUserErr
	}

	if userExists {
		data.Errors["Username"] = "User already exists"
	}

	// Email validation
	re = regexp.MustCompile(`^[\w-\.]+@([\w-]+\.)+[\w-]{2,}$`)
	matched = re.MatchString(email)
	if !matched || strings.Contains(email, " ") {
		data.Errors["Email"] = "Invalid email"
	} else {
		data.Email = email
	}

	emailQuery := "SELECT EXISTS(SELECT 1 FROM users WHERE email = ?)"
	emailStmt, emailErr := db.GetDB().Prepare(emailQuery)
	if emailErr != nil {
		http.Error(w, "Error preparing query", http.StatusInternalServerError)
		return emailErr
	}
	defer emailStmt.Close()

	var emailExists bool
	findEmailErr := emailStmt.QueryRow(email).Scan(&emailExists)
	if findEmailErr != nil {
		log.Println(findEmailErr)
		http.Error(w, "Error querying database", http.StatusInternalServerError)
		return findEmailErr
	}

	if emailExists {
		data.Errors["Email"] = "Account already exists! Please log-in"
	}

	// password validation
	reLower := regexp.MustCompile(`[a-z]`)
	reUpper := regexp.MustCompile(`[A-Z]`)
	reDigit := regexp.MustCompile(`\d`)
	
	if len(password) < 8 || len(password) > 128 {
		data.Errors["Password"] = "Password must be between 8 and 128 characters long"
	} else if !reLower.MatchString(password) {
		data.Errors["Password"] = "Password must contain at least one lowercase letter"
	} else if !reUpper.MatchString(password) {
		data.Errors["Password"] = "Password must contain at least one uppercase letter"
	} else if !reDigit.MatchString(password) {
		data.Errors["Password"] = "Password must contain at least one digit"
	} else {
		data.Password = password
	}	

	return nil
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
