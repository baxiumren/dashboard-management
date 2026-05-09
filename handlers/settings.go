package handlers

import (
	"html/template"
	"net/http"
	"strconv"

	"golang.org/x/crypto/bcrypt"
)

var (
	tmplSettings      *template.Template
	settingsGetHash   func() string
	onPasswordChanged func(newHash string)
	settingsGetCaptcha  func() string
	onCaptchaChanged    func(string)
	settingsGetSession  func() int
	onSessionChanged    func(int)
)

func InitSettings(tmpl *template.Template, getHash func() string, onChange func(string)) {
	tmplSettings = tmpl
	settingsGetHash = getHash
	onPasswordChanged = onChange
}

func InitSettingsFitur(getCaptcha func() string, onCaptcha func(string)) {
	settingsGetCaptcha = getCaptcha
	onCaptchaChanged = onCaptcha
}

func InitSettingsSession(getSession func() int, onChange func(int)) {
	settingsGetSession = getSession
	onSessionChanged = onChange
}


func SettingsGet(w http.ResponseWriter, r *http.Request) {
	captcha := "math"
	if settingsGetCaptcha != nil {
		if m := settingsGetCaptcha(); m != "" {
			captcha = m
		}
	}
	sessionHours := 1
	if settingsGetSession != nil {
		if h := settingsGetSession(); h > 0 {
			sessionHours = h
		}
	}
	tmplSettings.ExecuteTemplate(w, "settings.html", map[string]interface{}{
		"Page":           "settings",
		"Success":        r.URL.Query().Get("ok") == "1",
		"FiturOK":        r.URL.Query().Get("fitur") == "1",
		"SessionOK":      r.URL.Query().Get("sesi") == "1",
		"CaptchaMode":    captcha,
		"SessionHours":   sessionHours,
		"SessionOptions": []int{1, 2, 3, 4, 6, 8, 10, 12},
	})
}

func SettingsFiturPost(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	mode := r.FormValue("captcha_mode")
	if mode != "off" && mode != "math" {
		mode = "math"
	}
	if onCaptchaChanged != nil {
		onCaptchaChanged(mode)
	}
	http.Redirect(w, r, "/settings?fitur=1", http.StatusFound)
}

func SettingsSessionPost(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	h, err := strconv.Atoi(r.FormValue("session_hours"))
	if err != nil || h < 1 || h > 12 {
		h = 1
	}
	if onSessionChanged != nil {
		onSessionChanged(h)
	}
	http.Redirect(w, r, "/settings?sesi=1", http.StatusFound)
}

func SettingsPost(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()

	if r.FormValue("_action") == "fitur" {
		SettingsFiturPost(w, r)
		return
	}
	if r.FormValue("_action") == "session" {
		SettingsSessionPost(w, r)
		return
	}

	oldPass := r.FormValue("old_password")
	newPass := r.FormValue("new_password")
	confirm := r.FormValue("confirm_password")

	if bcrypt.CompareHashAndPassword([]byte(settingsGetHash()), []byte(oldPass)) != nil {
		tmplSettings.ExecuteTemplate(w, "settings.html", map[string]interface{}{
			"Page":  "settings",
			"Error": "Password lama salah",
		})
		return
	}
	if newPass != confirm || len(newPass) < 6 {
		tmplSettings.ExecuteTemplate(w, "settings.html", map[string]interface{}{
			"Page":  "settings",
			"Error": "Password baru tidak cocok atau kurang dari 6 karakter",
		})
		return
	}
	hash, _ := bcrypt.GenerateFromPassword([]byte(newPass), bcrypt.DefaultCost)
	onPasswordChanged(string(hash))
	http.Redirect(w, r, "/settings?ok=1", http.StatusFound)
}
