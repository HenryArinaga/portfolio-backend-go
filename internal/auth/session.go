package auth

import (
	"database/sql"
	"net/http"
	"time"

	"github.com/alexedwards/scs/sqlite3store"
	"github.com/alexedwards/scs/v2"
)

const AdminSessionKey = "is_admin"

func NewSessionManager(db *sql.DB, isProd bool) *scs.SessionManager {
	sm := scs.New()

	sm.Store = sqlite3store.New(db)

	sm.Lifetime = 7 * 24 * time.Hour
	sm.IdleTimeout = 24 * time.Hour

	sm.Cookie.Name = "admin_session"
	sm.Cookie.HttpOnly = true
	sm.Cookie.Path = "/"
	sm.Cookie.SameSite = http.SameSiteLaxMode

	sm.Cookie.Secure = isProd

	return sm
}
