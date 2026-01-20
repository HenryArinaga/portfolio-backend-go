package middleware

import (
	"net/http"

	"github.com/alexedwards/scs/v2"
	"github.com/henryarin/portfolio-backend-go/internal/auth"
)

func RequireAdminSession(sm *scs.SessionManager) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if !sm.GetBool(r.Context(), auth.AdminSessionKey) {
				http.Error(w, "unauthorized", http.StatusUnauthorized)
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}
