package middlewares

import (
	"forum/pkg/db"
	"log"
	"net/http"
	"time"
)

// AuthRequired is a middleware that checks if the user is authenticated.
// If not, it redirects the user to the login page.
func AuthRequired(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		// To avoid ERR_TOO_MANY_REDIRECTS, we need to check if the user is trying to access the login page
		if r.URL.Path == "/login" {
			next.ServeHTTP(w, r)
			return
		}

		cookie, cookieErr := r.Cookie("sessionId")

		if cookieErr != nil {
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}

		sessionsQuery := "SELECT expiresAt FROM sessions WHERE id = ?"
		var expiresAtSession time.Time

		sessionStmt, sessionErr := db.GetDB().Prepare(sessionsQuery)
		if sessionErr != nil {
			log.Println(sessionErr)
			http.Error(w, "Error preparing query", http.StatusInternalServerError)
			return
		}
		defer sessionStmt.Close()

		err := sessionStmt.QueryRow(cookie.Value).Scan(&expiresAtSession)
		if err != nil {
			log.Println(err)
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}

		if time.Now().After(expiresAtSession) {
			deleteSessionQuery := "DELETE FROM sessions WHERE id = ?"

			sessionExecErr := db.PrepareAndExecute(deleteSessionQuery, cookie.Value)
			if sessionExecErr != nil {
				log.Println(sessionExecErr)
				http.Error(w, "Error deleting session", http.StatusInternalServerError)
				return
			}

			http.SetCookie(w, &http.Cookie{
				Name:    "sessionId",
				Value:   "",
				Expires: time.Now(),
			})

			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}

		// check if user has more than one session
		// if so, delete all sessions except the current one

		moreThanOneSessionQuery := "SELECT userId FROM sessions WHERE id = ?"
		var userId int

		moreThanOneSessionStmt, moreThanOneSessionError := db.GetDB().Prepare(moreThanOneSessionQuery)
		if moreThanOneSessionError != nil {
			log.Println(moreThanOneSessionError)
			http.Error(w, "Error preparing query", http.StatusInternalServerError)
			return
		}
		defer moreThanOneSessionStmt.Close()

		moreThanOneSessionError = moreThanOneSessionStmt.QueryRow(cookie.Value).Scan(&userId)

		if moreThanOneSessionError != nil {
			log.Println(moreThanOneSessionError)
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}

		// delete all sessions except the current one

		deleteAllSessionsQuery := "DELETE FROM sessions WHERE userId = ? AND id != ?"

		deleteAllSessionsExecErr := db.PrepareAndExecute(deleteAllSessionsQuery, userId, cookie.Value)
		if deleteAllSessionsExecErr != nil {
			log.Println(deleteAllSessionsExecErr)
			http.Error(w, "Error deleting sessions", http.StatusInternalServerError)
			return
		}

		next.ServeHTTP(w, r)
	})
}
