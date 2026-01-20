package handlers

import (
	"encoding/json"
	"net/http"
	"os"

	"github.com/henryarin/portfolio-backend-go/internal/auth"

	"github.com/alexedwards/scs/v2"
)

type AdminAuthHandler struct {
	Sessions *scs.SessionManager
}

type loginReq struct {
	Token string `json:"token"`
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

	adminToken := os.Getenv("ADMIN_TOKEN")
	if adminToken == "" {
		http.Error(w, "server misconfigured", http.StatusInternalServerError)
		return
	}

	if req.Token != adminToken {
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
