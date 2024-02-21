package middlewares

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/alexedwards/scs/v2"
)

func Authmiddleware(sessionManager *scs.SessionManager) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			userID := sessionManager.GetInt32(r.Context(), "userID")

			isAuthenticated := userID > 0
            if isAuthenticated {
                if strings.Contains(r.URL.Path, "auth") {
                    fmt.Println("here")
                    http.Redirect(w, r, "/", http.StatusSeeOther)
                }
            } else {
                if !strings.Contains(r.URL.Path, "auth") {
                    http.Redirect(w, r, "/auth", http.StatusSeeOther)
                }
            }

			next.ServeHTTP(w, r)
		})
	}
}
