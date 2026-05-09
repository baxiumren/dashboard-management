package handlers

import (
	"fmt"
	"html/template"
	"net/http"
	"strconv"

	"dashboard-fb/models"
)

var tmplPixel *template.Template

func InitPixel(tmpl *template.Template) {
	tmplPixel = tmpl
}

func PixelList(w http.ResponseWriter, r *http.Request) {
	search := r.URL.Query().Get("search")
	status := r.URL.Query().Get("status")

	list, err := models.GetAllPixel(search, status)
	if err != nil {
		http.Error(w, "Error: "+err.Error(), 500)
		return
	}

	stats := models.GetPixelStats()
	bmList, _ := models.GetAllBMSimple()
	akunList, _ := models.GetAllAkunFBSimple()

	tmplPixel.ExecuteTemplate(w, "pixel.html", map[string]interface{}{
		"List":         list,
		"Stats":        stats,
		"BMList":       bmList,
		"AkunList":     akunList,
		"Search":       search,
		"StatusFilter": status,
		"Page":         "pixel",
	})
}

func PixelAdd(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	bmID, _ := strconv.Atoi(r.FormValue("bm_id"))
	akunFBID, _ := strconv.Atoi(r.FormValue("akun_fb_id"))

	p := &models.Pixel{
		Nama:     r.FormValue("nama"),
		PixelID:  r.FormValue("pixel_id"),
		BMID:     bmID,
		AkunFBID: akunFBID,
		Status:   r.FormValue("status"),
		TglBuat:  r.FormValue("tgl_buat"),
		Catatan:  r.FormValue("catatan"),
	}

	if err := models.CreatePixel(p); err != nil {
		http.Error(w, "Gagal simpan: "+err.Error(), 500)
		return
	}

	lastID := models.GetLastInsertedPixelID()
	if lastID > 0 && p.TglBuat != "" {
		models.CreateRiwayat(&models.Riwayat{
			Entitas:   "pixel",
			EntitasID: lastID,
			Tipe:      "buat",
			Tanggal:   p.TglBuat,
			Catatan:   fmt.Sprintf("Pixel dibuat: %s", p.Nama),
		})
	}

	http.Redirect(w, r, "/pixel?t=tambah", http.StatusFound)
}

func PixelEdit(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Query().Get("id")
	id, _ := strconv.Atoi(idStr)

	existing, err := models.GetPixelByID(id)
	if err != nil || existing == nil {
		http.Redirect(w, r, "/pixel", http.StatusFound)
		return
	}

	r.ParseForm()
	bmID, _ := strconv.Atoi(r.FormValue("bm_id"))
	akunFBID, _ := strconv.Atoi(r.FormValue("akun_fb_id"))

	existing.Nama     = r.FormValue("nama")
	existing.PixelID  = r.FormValue("pixel_id")
	existing.BMID     = bmID
	existing.AkunFBID = akunFBID
	existing.Status   = r.FormValue("status")
	existing.TglBuat  = r.FormValue("tgl_buat")
	existing.Catatan  = r.FormValue("catatan")

	if err := models.UpdatePixel(existing); err != nil {
		http.Error(w, "Gagal update: "+err.Error(), 500)
		return
	}
	http.Redirect(w, r, "/pixel?t=edit", http.StatusFound)
}

func PixelDelete(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Query().Get("id")
	id, _ := strconv.Atoi(idStr)
	models.DeletePixel(id)
	http.Redirect(w, r, "/pixel?t=hapus", http.StatusFound)
}

func PixelRiwayat(w http.ResponseWriter, r *http.Request) {
	if r.Header.Get("X-Fragment") != "1" {
		http.Redirect(w, r, "/pixel", http.StatusFound)
		return
	}
	idStr := r.URL.Query().Get("id")
	id, _ := strconv.Atoi(idStr)

	p, _ := models.GetPixelByID(id)
	riwayat, _ := models.GetRiwayat("pixel", id)

	tmplPixel.ExecuteTemplate(w, "pixel_riwayat.html", map[string]interface{}{
		"Akun":    p,
		"Riwayat": riwayat,
	})
}

func PixelRiwayatAdd(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	id, _ := strconv.Atoi(r.FormValue("entitas_id"))
	models.CreateRiwayat(&models.Riwayat{
		Entitas:   "pixel",
		EntitasID: id,
		Tipe:      r.FormValue("tipe"),
		Tanggal:   r.FormValue("tanggal"),
		Catatan:   r.FormValue("catatan"),
	})
	http.Redirect(w, r, "/pixel?t=riwayat", http.StatusFound)
}
