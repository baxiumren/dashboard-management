package main

import (
	"net/http"

	"dashboard-fb/handlers"

	"github.com/gorilla/mux"
)

func registerRoutes(r *mux.Router) {
	// ── Error handlers ──
	r.NotFoundHandler = http.HandlerFunc(handlers.Error404)
	r.MethodNotAllowedHandler = http.HandlerFunc(handlers.Error404)

	// ── Public ──
	r.HandleFunc("/403", handlers.Error403)
	r.HandleFunc("/503", handlers.Error503)

	// ── Public ──
	r.HandleFunc("/login", handlers.LoginPage).Methods("GET")
	r.HandleFunc("/login", func(w http.ResponseWriter, r *http.Request) {
		handlers.LoginPost(w, r, cfg.PasswordHash)
	}).Methods("POST")
	r.HandleFunc("/logout", handlers.Logout).Methods("GET")

	r.PathPrefix("/static/").Handler(
		http.StripPrefix("/static/", http.FileServer(http.Dir("static"))),
	)

	// ── Protected ──
	protected := r.NewRoute().Subrouter()
	protected.Use(handlers.AuthMiddleware)

	protected.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/akun-fb", http.StatusFound)
	})

	// Akun FB
	protected.HandleFunc("/akun-fb", func(w http.ResponseWriter, r *http.Request) {
		handlers.AkunFBList(w, r, decrypt, encrypt, cfg.IsDefaultPassword)
	}).Methods("GET")
	protected.HandleFunc("/akun-fb/tambah", func(w http.ResponseWriter, r *http.Request) {
		handlers.AkunFBAdd(w, r, encrypt)
	}).Methods("POST")
	protected.HandleFunc("/akun-fb/edit", func(w http.ResponseWriter, r *http.Request) {
		handlers.AkunFBEdit(w, r, encrypt, decrypt)
	}).Methods("POST")
	protected.HandleFunc("/akun-fb/hapus", handlers.AkunFBDelete).Methods("GET")
	protected.HandleFunc("/akun-fb/import", func(w http.ResponseWriter, r *http.Request) {
		handlers.AkunFBImportBulk(w, r, encrypt)
	}).Methods("POST")
	protected.HandleFunc("/akun-fb/riwayat", handlers.AkunFBRiwayat).Methods("GET")
	protected.HandleFunc("/akun-fb/riwayat/tambah", handlers.AkunFBRiwayatAdd).Methods("POST")

	// Fanpage
	protected.HandleFunc("/fanpage", handlers.FanpageList).Methods("GET")
	protected.HandleFunc("/fanpage/tambah", handlers.FanpageAdd).Methods("POST")
	protected.HandleFunc("/fanpage/edit", handlers.FanpageEdit).Methods("POST")
	protected.HandleFunc("/fanpage/hapus", handlers.FanpageDelete).Methods("GET")
	protected.HandleFunc("/fanpage/riwayat", handlers.FanpageRiwayat).Methods("GET")
	protected.HandleFunc("/fanpage/riwayat/tambah", handlers.FanpageRiwayatAdd).Methods("POST")

	// Business Manager
	protected.HandleFunc("/bm", handlers.BMList).Methods("GET")
	protected.HandleFunc("/bm/tambah", handlers.BMAdd).Methods("POST")
	protected.HandleFunc("/bm/edit", handlers.BMEdit).Methods("POST")
	protected.HandleFunc("/bm/hapus", handlers.BMDelete).Methods("GET")
	protected.HandleFunc("/bm/riwayat", handlers.BMRiwayat).Methods("GET")
	protected.HandleFunc("/bm/riwayat/tambah", handlers.BMRiwayatAdd).Methods("POST")

	// Akun Iklan
	protected.HandleFunc("/akun-iklan", handlers.AkunIklanList).Methods("GET")
	protected.HandleFunc("/akun-iklan/tambah", handlers.AkunIklanAdd).Methods("POST")
	protected.HandleFunc("/akun-iklan/edit", handlers.AkunIklanEdit).Methods("POST")
	protected.HandleFunc("/akun-iklan/hapus", handlers.AkunIklanDelete).Methods("GET")
	protected.HandleFunc("/akun-iklan/riwayat", handlers.AkunIklanRiwayat).Methods("GET")
	protected.HandleFunc("/akun-iklan/riwayat/tambah", handlers.AkunIklanRiwayatAdd).Methods("POST")

	// Pixel
	protected.HandleFunc("/pixel", handlers.PixelList).Methods("GET")
	protected.HandleFunc("/pixel/tambah", handlers.PixelAdd).Methods("POST")
	protected.HandleFunc("/pixel/edit", handlers.PixelEdit).Methods("POST")
	protected.HandleFunc("/pixel/hapus", handlers.PixelDelete).Methods("GET")
	protected.HandleFunc("/pixel/riwayat", handlers.PixelRiwayat).Methods("GET")
	protected.HandleFunc("/pixel/riwayat/tambah", handlers.PixelRiwayatAdd).Methods("POST")

	// Settings
	protected.HandleFunc("/settings", handlers.SettingsGet).Methods("GET")
	protected.HandleFunc("/settings", handlers.SettingsPost).Methods("POST")

	// Rekap
	protected.HandleFunc("/rekap", handlers.RekapPage).Methods("GET")

}
