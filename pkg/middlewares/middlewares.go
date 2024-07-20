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
		var expiresAt time.Time

		sessionStmt, sessionErr := db.GetDB().Prepare(sessionsQuery)
		if sessionErr != nil {
			log.Println(sessionErr)
			http.Error(w, "Error preparing query", http.StatusInternalServerError)
			return
		}
		defer sessionStmt.Close()

		err := sessionStmt.QueryRow(cookie.Value).Scan(&expiresAt)
		if err != nil {
			log.Println(err)
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}

		if time.Now().After(expiresAt) {
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

		next.ServeHTTP(w, r)
	})
}
