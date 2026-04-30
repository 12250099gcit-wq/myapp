// myapp/controller/admin.go
package controller

import (
	"encoding/json"
	"net/http"
	"time"

	"myapp/model"
	"myapp/utils/httpResp"
)

// Signup handles POST /signup
func Signup(w http.ResponseWriter, r *http.Request) {
	var admin model.Admin
	if err := json.NewDecoder(r.Body).Decode(&admin); err != nil {
		httpResp.RespondWithError(w, http.StatusBadRequest, "invalid json body")
		return
	}
	defer r.Body.Close()

	if saveErr := admin.Create(); saveErr != nil {
		httpResp.RespondWithError(w, http.StatusBadRequest, saveErr.Error())
		return
	}
	httpResp.RespondWithJSON(w, http.StatusCreated, map[string]string{"status": "admin added"})
}

// Login handles POST /login
func Login(w http.ResponseWriter, r *http.Request) {
	var admin model.Admin
	if err := json.NewDecoder(r.Body).Decode(&admin); err != nil {
		httpResp.RespondWithError(w, http.StatusBadRequest, "invalid json body")
		return
	}
	defer r.Body.Close()

	if getErr := admin.Get(); getErr != nil {
		httpResp.RespondWithError(w, http.StatusUnauthorized, getErr.Error())
		return
	}

	// Set session cookie on successful login.
	// Secure false is used here so the test server and local development
	// environment can use the cookie over HTTP.
	cookie := http.Cookie{
		Name:     "my-cookie",
		Value:    "my-value",
		Expires:  time.Now().Add(30 * time.Minute),
		Secure:   false,
		HttpOnly: true,
	}
	http.SetCookie(w, &cookie)

	httpResp.RespondWithJSON(w, http.StatusOK, map[string]string{"message": "login success"})
}

// Logout handles GET /logout
// Destroys the session by expiring the cookie immediately.
func Logout(w http.ResponseWriter, r *http.Request) {
	http.SetCookie(w, &http.Cookie{
		Name:    "my-cookie",
		Expires: time.Now(), // setting Expires to now causes the browser to delete it
	})
	httpResp.RespondWithJSON(w, http.StatusOK, map[string]string{"message": "cookie deleted"})
}

// VerifyCookie checks that the incoming request has a valid session cookie.
// Returns true if the cookie is valid, false otherwise (response already written).
//
// Call this at the top of every protected handler:
//
//	if !VerifyCookie(w, r) { return }
func VerifyCookie(w http.ResponseWriter, r *http.Request) bool {
	cookie, err := r.Cookie("my-cookie")
	if err != nil {
		if err == http.ErrNoCookie {
			httpResp.RespondWithError(w, http.StatusSeeOther, "cookie not found")
			return false
		}
		httpResp.RespondWithError(w, http.StatusInternalServerError, "internal server error")
		return false
	}

	if cookie.Value != "my-value" {
		httpResp.RespondWithError(w, http.StatusUnauthorized, "cookie does not match")
		return false
	}
	return true
}
