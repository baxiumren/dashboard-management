package handlers

import (
	"fmt"
	"html/template"
	"net/http"
	"strconv"

	"dashboard-fb/models"
)

var tmplBM *template.Template

func InitBM(tmpl *template.Template) {
	tmplBM = tmpl
}

func BMList(w http.ResponseWriter, r *http.Request) {
	search := r.URL.Query().Get("search")
	status := r.URL.Query().Get("status")

	list, err := models.GetAllBM(search, status)
	if err != nil {
		http.Error(w, "Error: "+err.Error(), 500)
		return
	}

	stats := models.GetBMStats()
	akunList, _ := models.GetAllAkunFBSimple()

	tmplBM.ExecuteTemplate(w, "bm.html", map[string]interface{}{
		"List":         list,
		"Stats":        stats,
		"AkunList":     akunList,
		"Search":       search,
		"StatusFilter": status,
		"Page":         "bm",
	})
}

func BMAdd(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	ownerID, _ := strconv.Atoi(r.FormValue("owner_akun_fb_id"))

	b := &models.BM{
		Nama:          r.FormValue("nama"),
		BMID:          r.FormValue("bm_id"),
		OwnerAkunFBID: ownerID,
		Status:        r.FormValue("status"),
		TglBuat:       r.FormValue("tgl_buat"),
		Catatan:       r.FormValue("catatan"),
	}

	if err := models.CreateBM(b); err != nil {
		http.Error(w, "Gagal simpan: "+err.Error(), 500)
		return
	}

	lastID := models.GetLastInsertedBMID()
	if lastID > 0 && b.TglBuat != "" {
		models.CreateRiwayat(&models.Riwayat{
			Entitas:   "bm",
			EntitasID: lastID,
			Tipe:      "buat",
			Tanggal:   b.TglBuat,
			Catatan:   fmt.Sprintf("BM dibuat: %s", b.Nama),
		})
	}

	http.Redirect(w, r, "/bm?t=tambah", http.StatusFound)
}

func BMEdit(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Query().Get("id")
	id, _ := strconv.Atoi(idStr)

	existing, err := models.GetBMByID(id)
	if err != nil || existing == nil {
		http.Redirect(w, r, "/bm", http.StatusFound)
		return
	}

	r.ParseForm()
	ownerID, _ := strconv.Atoi(r.FormValue("owner_akun_fb_id"))

	existing.Nama          = r.FormValue("nama")
	existing.BMID          = r.FormValue("bm_id")
	existing.OwnerAkunFBID = ownerID
	existing.Status        = r.FormValue("status")
	existing.TglBuat       = r.FormValue("tgl_buat")
	existing.Catatan       = r.FormValue("catatan")

	if err := models.UpdateBM(existing); err != nil {
		http.Error(w, "Gagal update: "+err.Error(), 500)
		return
	}
	http.Redirect(w, r, "/bm?t=edit", http.StatusFound)
}

func BMDelete(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Query().Get("id")
	id, _ := strconv.Atoi(idStr)
	models.DeleteBM(id)
	http.Redirect(w, r, "/bm?t=hapus", http.StatusFound)
}

func BMRiwayat(w http.ResponseWriter, r *http.Request) {
	if r.Header.Get("X-Fragment") != "1" {
		http.Redirect(w, r, "/bm", http.StatusFound)
		return
	}
	idStr := r.URL.Query().Get("id")
	id, _ := strconv.Atoi(idStr)

	b, _ := models.GetBMByID(id)
	riwayat, _ := models.GetRiwayat("bm", id)

	tmplBM.ExecuteTemplate(w, "bm_riwayat.html", map[string]interface{}{
		"Akun":    b,
		"Riwayat": riwayat,
	})
}

func BMRiwayatAdd(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	id, _ := strconv.Atoi(r.FormValue("entitas_id"))
	models.CreateRiwayat(&models.Riwayat{
		Entitas:   "bm",
		EntitasID: id,
		Tipe:      r.FormValue("tipe"),
		Tanggal:   r.FormValue("tanggal"),
		Catatan:   r.FormValue("catatan"),
	})
	http.Redirect(w, r, "/bm?t=riwayat", http.StatusFound)
}
