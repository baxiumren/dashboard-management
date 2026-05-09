package handlers

import (
	"html/template"
	"net/http"

	"dashboard-fb/models"
)

var tmplRekap *template.Template

func InitRekap(tmpl *template.Template) {
	tmplRekap = tmpl
}

func RekapPage(w http.ResponseWriter, r *http.Request) {
	fb := models.GetAkunFBStats()
	fp := models.GetFanpageStats()
	bm := models.GetBMStats()
	ai := models.GetAkunIklanStats()
	px := models.GetPixelStats()

	totalAset := fb.Total + fp.Total + bm.Total + ai.Total + px.Total
	totalAktif := fb.Aktif + fp.Aktif + bm.Aktif + ai.Aktif + px.Aktif

	tmplRekap.ExecuteTemplate(w, "rekap.html", map[string]interface{}{
		"Page":           "rekap",
		"AkunFBStats":    fb,
		"FanpageStats":   fp,
		"BMStats":        bm,
		"AkunIklanStats": ai,
		"PixelStats":     px,
		"TotalAset":      totalAset,
		"TotalAktif":     totalAktif,
	})
}
