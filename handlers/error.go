package handlers

import (
	"html/template"
	"net/http"
)

var tmplError *template.Template

func InitError(tmpl *template.Template) {
	tmplError = tmpl
}

type errorData struct {
	Code    int
	Title   string
	Message string
}

func renderError(w http.ResponseWriter, code int, title, message string) {
	w.WriteHeader(code)
	tmplError.ExecuteTemplate(w, "error.html", errorData{
		Code:    code,
		Title:   title,
		Message: message,
	})
}

func Error404(w http.ResponseWriter, r *http.Request) {
	renderError(w, 404,
		"Halaman Tidak Ditemukan",
		"Halaman yang kamu cari tidak ada atau sudah dipindahkan. Pastikan URL sudah benar.",
	)
}

func Error403(w http.ResponseWriter, r *http.Request) {
	renderError(w, 403,
		"Akses Ditolak",
		"Kamu tidak punya izin untuk mengakses halaman ini. Silakan login terlebih dahulu.",
	)
}

func Error503(w http.ResponseWriter, r *http.Request) {
	renderError(w, 503,
		"Layanan Tidak Tersedia",
		"Server sedang dalam maintenance atau mengalami gangguan. Coba lagi dalam beberapa saat.",
	)
}
