package handlers

import (
	"crypto/rand"
	"encoding/hex"
	"html/template"
	"net/http"
	"sync"
	"time"

	"golang.org/x/crypto/bcrypt"
)

var (
	sessions   = sync.Map{}
	cookieName = "fb_session"
	tmplLogin  *template.Template
)

func InitAuth(tmpl *template.Template) {
	tmplLogin = tmpl
}

func captchaMode() string {
	if getCaptchaMode != nil {
		if m := getCaptchaMode(); m != "" {
			return m
		}
	}
	return "math"
}

func LoginPage(w http.ResponseWriter, r *http.Request) {
	if isLoggedIn(r) {
		http.Redirect(w, r, "/akun-fb", http.StatusFound)
		return
	}
	tmplLogin.ExecuteTemplate(w, "login.html", map[string]interface{}{
		"Error":       "",
		"CaptchaMode": captchaMode(),
	})
}

func LoginPost(w http.ResponseWriter, r *http.Request, passwordHash string) {
	r.ParseForm()
	username := r.FormValue("username")
	password := r.FormValue("password")

	if username != GetUsername() || bcrypt.CompareHashAndPassword([]byte(passwordHash), []byte(password)) != nil {
		tmplLogin.ExecuteTemplate(w, "login.html", map[string]interface{}{
			"Error":       "Username atau password salah",
			"CaptchaMode": captchaMode(),
		})
		return
	}

	token := generateToken()
	sessions.Store(token, time.Now())

	dur := sessionDuration()
	http.SetCookie(w, &http.Cookie{
		Name:     cookieName,
		Value:    token,
		Path:     "/",
		HttpOnly: true,
		MaxAge:   int(dur.Seconds()),
	})
	http.Redirect(w, r, "/akun-fb", http.StatusFound)
}

func Logout(w http.ResponseWriter, r *http.Request) {
	if cookie, err := r.Cookie(cookieName); err == nil {
		sessions.Delete(cookie.Value)
	}
	http.SetCookie(w, &http.Cookie{
		Name:   cookieName,
		Value:  "",
		Path:   "/",
		MaxAge: -1,
	})
	http.Redirect(w, r, "/login", http.StatusFound)
}

func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !isLoggedIn(r) {
			http.Redirect(w, r, "/login", http.StatusFound)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func isLoggedIn(r *http.Request) bool {
	cookie, err := r.Cookie(cookieName)
	if err != nil {
		return false
	}
	val, ok := sessions.Load(cookie.Value)
	if !ok {
		return false
	}
	loginTime := val.(time.Time)
	if time.Since(loginTime) >= sessionDuration() {
		sessions.Delete(cookie.Value)
		return false
	}
	return true
}

func generateToken() string {
	b := make([]byte, 32)
	rand.Read(b)
	return hex.EncodeToString(b)
}

// These vars will be set from main
var (
	getUsername     func() string
	getPasswordHash func() string
	getCaptchaMode  func() string
	getSessionHours func() int
)

func SetAuthFuncs(userFn func() string, passFn func() string) {
	getUsername = userFn
	getPasswordHash = passFn
}

func SetCaptchaModeFn(fn func() string) {
	getCaptchaMode = fn
}

func SetSessionHoursFn(fn func() int) {
	getSessionHours = fn
}

func sessionDuration() time.Duration {
	h := 1
	if getSessionHours != nil {
		if v := getSessionHours(); v > 0 {
			h = v
		}
	}
	return time.Duration(h) * time.Hour
}

func GetUsername() string {
	if getUsername != nil {
		return getUsername()
	}
	return "admin"
}
