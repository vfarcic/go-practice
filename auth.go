package main

import (
	"net/http"
	"strings"
	"fmt"
	"github.com/parnurzeal/gorequest"
	"encoding/json"
	"github.com/stretchr/objx"
	"log"
)

type authHandler struct {
	next http.Handler
}

func (h *authHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if cookie, err := r.Cookie("authID"); err == http.ErrNoCookie || len(cookie.Value) == 0 {
		w.Header().Set("Location", "/login")
		w.WriteHeader(http.StatusTemporaryRedirect)
	} else if err != nil {
		panic(err.Error())
	} else {
		h.next.ServeHTTP(w, r)
	}
}

func MustAuth(handler http.Handler) http.Handler {
	return &authHandler{next: handler}
}

// Format: /auth/{provider}/{action}
func loginHandler(w http.ResponseWriter, r *http.Request) {
	segs := strings.Split(r.URL.Path, "/")
	action := segs[3]
	switch action {
	case "callback":
		authID := r.URL.Query().Get("authID")
		log.Println(r.URL.Query())
		// TODO: Change to params
		_, body, err := gorequest.New().Get("http://localhost:8080/auth/api/v1/user/" + authID).End()
		if len(err) > 0 {
			http.Error(w, err[0].Error(), http.StatusInternalServerError)
			return
		}
		var user User
		json.Unmarshal([]byte(body), &user)
		authCookieValue := objx.New(map[string]interface {}{
			"name": user.Name,
			"avatarURL": user.AvatarURL,
		}).MustBase64()
		http.SetCookie(w, &http.Cookie{
			Name: "auth",
			Value: authCookieValue,
			Path: "/"})
		w.Header()["Location"] = []string{"/chat?authID=" + user.AuthID}
		w.WriteHeader(http.StatusTemporaryRedirect)
	default:
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprintf(w, "Auth action %s not supported", action)
	}
}

func logoutHandler(w http.ResponseWriter, r *http.Request) {
	http.SetCookie(w, &http.Cookie{
		Name: "auth",
		Value: "",
		Path: "/",
		MaxAge: -1,
	})
	w.Header()["Location"] = []string{"/chat"}
	w.WriteHeader(http.StatusTemporaryRedirect)
}

type User struct {
	Email string
	Name string
	Nickname string
	AvatarURL string `json:"avatar_url"`
	AuthID string `json:"auth_id"`
}