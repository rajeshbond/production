package auth

import (
	"net/http"
	"strconv"

	"github.com/alexedwards/scs/v2"
)

// Session is injected from main
var Session *scs.SessionManager

// AuthRequired middleware
func AuthRequired(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if Session == nil {
			http.Error(w, "Session not initialized", http.StatusInternalServerError)
			return
		}
		userID := Session.GetInt(r.Context(), "userID")
		if userID == 0 {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}
		next(w, r)
	}
}

// DashboardHandler example
func DashboardHandler(w http.ResponseWriter, r *http.Request) {
	userID := Session.GetInt(r.Context(), "userID")
	w.Write([]byte("Welcome User ID: " + strconv.Itoa(userID)))
}
