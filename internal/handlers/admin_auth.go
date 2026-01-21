package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/alexedwards/scs/v2"
	"github.com/henryarin/portfolio-backend-go/internal/auth"
	"github.com/henryarin/portfolio-backend-go/internal/config"
)

type AdminAuthHandler struct {
	Sessions *scs.SessionManager
	Config   config.Config
}

type loginReq struct {
	Password string `json:"password"`
}

func (h *AdminAuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusNoContent)
		return
	}

	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req loginReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid json", http.StatusBadRequest)
		return
	}

	if h.Config.AdminPassword == "" {
		http.Error(w, "server misconfigured", http.StatusInternalServerError)
		return
	}

	if req.Password != h.Config.AdminPassword {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	h.Sessions.Put(r.Context(), auth.AdminSessionKey, true)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]any{"ok": true})
}

func (h *AdminAuthHandler) Logout(w http.ResponseWriter, r *http.Request) {
	if err := h.Sessions.Destroy(r.Context()); err != nil {
		http.Error(w, "failed to logout", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
