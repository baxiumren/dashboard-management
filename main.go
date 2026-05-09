package main

import (
	"html/template"
	"log"
	"net/http"

	"dashboard-fb/db"
	"dashboard-fb/handlers"

	"github.com/gorilla/mux"
)

func main() {
	loadConfig()
	db.Init("data.db")

	funcMap := buildFuncMap()

	// Base template: hanya layout (sidebar, loading, popup)
	baseTmpl := template.Must(
		template.New("").Funcs(funcMap).ParseFiles("templates/layout.html"),
	)
	// Login berdiri sendiri (tidak pakai layout)
	loginTmpl := template.Must(
		template.New("").Funcs(funcMap).ParseFiles("templates/login.html"),
	)

	// Per-halaman: clone base lalu parse file halaman + partial yang dibutuhkan
	akunFBTmpl      := cloneParse(baseTmpl, "templates/pages/akun_fb.html", "templates/partials/riwayat.html")
	fanpageTmpl     := cloneParse(baseTmpl, "templates/pages/fanpage.html", "templates/partials/fanpage_riwayat.html")
	bmTmpl          := cloneParse(baseTmpl, "templates/pages/bm.html", "templates/partials/bm_riwayat.html")
	akunIklanTmpl   := cloneParse(baseTmpl, "templates/pages/akun_iklan.html", "templates/partials/akun_iklan_riwayat.html")
	pixelTmpl       := cloneParse(baseTmpl, "templates/pages/pixel.html", "templates/partials/pixel_riwayat.html")
	settingsTmpl    := cloneParse(baseTmpl, "templates/pages/settings.html")
	rekapTmpl       := cloneParse(baseTmpl, "templates/pages/rekap.html")
	errorTmpl       := template.Must(template.New("").ParseFiles("templates/error.html"))

	// Inisialisasi semua handler dengan template masing-masing
	handlers.InitAuth(loginTmpl)
	handlers.SetAuthFuncs(
		func() string { return cfg.Username },
		func() string { return cfg.PasswordHash },
	)
	handlers.SetCaptchaModeFn(func() string { return cfg.CaptchaMode })
	handlers.SetSessionHoursFn(func() int { return cfg.SessionHours })
	handlers.InitAkunFB(akunFBTmpl)
	handlers.InitFanpage(fanpageTmpl)
	handlers.InitBM(bmTmpl)
	handlers.InitAkunIklan(akunIklanTmpl)
	handlers.InitPixel(pixelTmpl)
	handlers.InitSettings(settingsTmpl,
		func() string { return cfg.PasswordHash },
		func(newHash string) {
			cfg.PasswordHash = newHash
			cfg.IsDefaultPassword = false
			saveConfig()
		},
	)
	handlers.InitSettingsFitur(
		func() string { return cfg.CaptchaMode },
		func(mode string) {
			cfg.CaptchaMode = mode
			saveConfig()
		},
	)
	handlers.InitSettingsSession(
		func() int { return cfg.SessionHours },
		func(h int) {
			cfg.SessionHours = h
			saveConfig()
		},
	)
	handlers.InitRekap(rekapTmpl)
	handlers.InitError(errorTmpl)

	r := mux.NewRouter()
	registerRoutes(r)

	addr := ":" + cfg.Port
	log.Printf("FB Manager jalan di http://localhost%s", addr)
	log.Fatal(http.ListenAndServe(addr, r))
}
