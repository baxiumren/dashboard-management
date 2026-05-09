package handlers

import (
	"fmt"
	"html/template"
	"net/http"
	"strconv"

	"dashboard-fb/models"
)

var tmplFanpage *template.Template

func InitFanpage(tmpl *template.Template) {
	tmplFanpage = tmpl
}

func FanpageList(w http.ResponseWriter, r *http.Request) {
	search := r.URL.Query().Get("search")
	status := r.URL.Query().Get("status")

	list, err := models.GetAllFanpage(search, status)
	if err != nil {
		http.Error(w, "Error: "+err.Error(), 500)
		return
	}

	stats := models.GetFanpageStats()
	akunList, _ := models.GetAllAkunFBSimple()

	tmplFanpage.ExecuteTemplate(w, "fanpage.html", map[string]interface{}{
		"List":         list,
		"Stats":        stats,
		"AkunList":     akunList,
		"Search":       search,
		"StatusFilter": status,
		"Page":         "fanpage",
	})
}

func FanpageAdd(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	akunFBID, _ := strconv.Atoi(r.FormValue("akun_fb_id"))
	bmID, _ := strconv.Atoi(r.FormValue("bm_id"))

	fp := &models.Fanpage{
		Nama:    r.FormValue("nama"),
		PageID:  r.FormValue("page_id"),
		AkunFBID: akunFBID,
		BMID:    bmID,
		Status:  r.FormValue("status"),
		TglBuat: r.FormValue("tgl_buat"),
		Catatan: r.FormValue("catatan"),
	}

	if err := models.CreateFanpage(fp); err != nil {
		http.Error(w, "Gagal simpan: "+err.Error(), 500)
		return
	}

	// Auto-create riwayat "buat"
	lastID := models.GetLastInsertedFanpageID()
	if lastID > 0 && fp.TglBuat != "" {
		models.CreateRiwayat(&models.Riwayat{
			Entitas:   "fanpage",
			EntitasID: lastID,
			Tipe:      "buat",
			Tanggal:   fp.TglBuat,
			Catatan:   fmt.Sprintf("Fanpage dibuat: %s", fp.Nama),
		})
	}

	http.Redirect(w, r, "/fanpage?t=tambah", http.StatusFound)
}

func FanpageEdit(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Query().Get("id")
	id, _ := strconv.Atoi(idStr)

	existing, err := models.GetFanpageByID(id)
	if err != nil || existing == nil {
		http.Redirect(w, r, "/fanpage", http.StatusFound)
		return
	}

	r.ParseForm()
	akunFBID, _ := strconv.Atoi(r.FormValue("akun_fb_id"))
	bmID, _ := strconv.Atoi(r.FormValue("bm_id"))

	existing.Nama     = r.FormValue("nama")
	existing.PageID   = r.FormValue("page_id")
	existing.AkunFBID = akunFBID
	existing.BMID     = bmID
	existing.Status   = r.FormValue("status")
	existing.TglBuat  = r.FormValue("tgl_buat")
	existing.Catatan  = r.FormValue("catatan")

	if err := models.UpdateFanpage(existing); err != nil {
		http.Error(w, "Gagal update: "+err.Error(), 500)
		return
	}
	http.Redirect(w, r, "/fanpage?t=edit", http.StatusFound)
}

func FanpageDelete(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Query().Get("id")
	id, _ := strconv.Atoi(idStr)
	models.DeleteFanpage(id)
	http.Redirect(w, r, "/fanpage?t=hapus", http.StatusFound)
}

func FanpageRiwayat(w http.ResponseWriter, r *http.Request) {
	if r.Header.Get("X-Fragment") != "1" {
		http.Redirect(w, r, "/fanpage", http.StatusFound)
		return
	}
	idStr := r.URL.Query().Get("id")
	id, _ := strconv.Atoi(idStr)

	fp, _ := models.GetFanpageByID(id)
	riwayat, _ := models.GetRiwayat("fanpage", id)

	tmplFanpage.ExecuteTemplate(w, "fanpage_riwayat.html", map[string]interface{}{
		"Akun":    fp,
		"Riwayat": riwayat,
	})
}

func FanpageRiwayatAdd(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	id, _ := strconv.Atoi(r.FormValue("entitas_id"))
	models.CreateRiwayat(&models.Riwayat{
		Entitas:   "fanpage",
		EntitasID: id,
		Tipe:      r.FormValue("tipe"),
		Tanggal:   r.FormValue("tanggal"),
		Catatan:   r.FormValue("catatan"),
	})
	http.Redirect(w, r, "/fanpage?t=riwayat", http.StatusFound)
}
