package handlers

import (
	"fmt"
	"html/template"
	"net/http"
	"strconv"

	"dashboard-fb/models"
)

var tmplAkunIklan *template.Template

func InitAkunIklan(tmpl *template.Template) {
	tmplAkunIklan = tmpl
}

func AkunIklanList(w http.ResponseWriter, r *http.Request) {
	search := r.URL.Query().Get("search")
	status := r.URL.Query().Get("status")

	list, err := models.GetAllAkunIklan(search, status)
	if err != nil {
		http.Error(w, "Error: "+err.Error(), 500)
		return
	}

	stats := models.GetAkunIklanStats()
	bmList, _ := models.GetAllBMSimple()
	akunList, _ := models.GetAllAkunFBSimple()

	tmplAkunIklan.ExecuteTemplate(w, "akun_iklan.html", map[string]interface{}{
		"List":         list,
		"Stats":        stats,
		"BMList":       bmList,
		"AkunList":     akunList,
		"Search":       search,
		"StatusFilter": status,
		"Page":         "akun-iklan",
	})
}

func AkunIklanAdd(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	bmID, _ := strconv.Atoi(r.FormValue("bm_id"))
	runnerID, _ := strconv.Atoi(r.FormValue("runner_akun_fb_id"))
	limitHarian, _ := strconv.Atoi(r.FormValue("limit_harian"))
	limitTotal, _ := strconv.Atoi(r.FormValue("limit_total"))

	ai := &models.AkunIklan{
		Nama:           r.FormValue("nama"),
		AdAccountID:    r.FormValue("ad_account_id"),
		BMID:           bmID,
		RunnerAkunFBID: runnerID,
		LimitHarian:    limitHarian,
		LimitTotal:     limitTotal,
		MetodeBayar:    r.FormValue("metode_bayar"),
		MataUang:       r.FormValue("mata_uang"),
		Status:         r.FormValue("status"),
		TglBuat:        r.FormValue("tgl_buat"),
		Catatan:        r.FormValue("catatan"),
	}

	if err := models.CreateAkunIklan(ai); err != nil {
		http.Error(w, "Gagal simpan: "+err.Error(), 500)
		return
	}

	lastID := models.GetLastInsertedAkunIklanID()
	if lastID > 0 && ai.TglBuat != "" {
		models.CreateRiwayat(&models.Riwayat{
			Entitas:   "akun_iklan",
			EntitasID: lastID,
			Tipe:      "buat",
			Tanggal:   ai.TglBuat,
			Catatan:   fmt.Sprintf("Akun iklan dibuat: %s", ai.Nama),
		})
	}

	http.Redirect(w, r, "/akun-iklan?t=tambah", http.StatusFound)
}

func AkunIklanEdit(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Query().Get("id")
	id, _ := strconv.Atoi(idStr)

	existing, err := models.GetAkunIklanByID(id)
	if err != nil || existing == nil {
		http.Redirect(w, r, "/akun-iklan", http.StatusFound)
		return
	}

	r.ParseForm()
	bmID, _ := strconv.Atoi(r.FormValue("bm_id"))
	runnerID, _ := strconv.Atoi(r.FormValue("runner_akun_fb_id"))
	limitHarian, _ := strconv.Atoi(r.FormValue("limit_harian"))
	limitTotal, _ := strconv.Atoi(r.FormValue("limit_total"))

	existing.Nama           = r.FormValue("nama")
	existing.AdAccountID    = r.FormValue("ad_account_id")
	existing.BMID           = bmID
	existing.RunnerAkunFBID = runnerID
	existing.LimitHarian    = limitHarian
	existing.LimitTotal     = limitTotal
	existing.MetodeBayar    = r.FormValue("metode_bayar")
	existing.MataUang       = r.FormValue("mata_uang")
	existing.Status         = r.FormValue("status")
	existing.TglBuat        = r.FormValue("tgl_buat")
	existing.Catatan        = r.FormValue("catatan")

	if err := models.UpdateAkunIklan(existing); err != nil {
		http.Error(w, "Gagal update: "+err.Error(), 500)
		return
	}
	http.Redirect(w, r, "/akun-iklan?t=edit", http.StatusFound)
}

func AkunIklanDelete(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Query().Get("id")
	id, _ := strconv.Atoi(idStr)
	models.DeleteAkunIklan(id)
	http.Redirect(w, r, "/akun-iklan?t=hapus", http.StatusFound)
}

func AkunIklanRiwayat(w http.ResponseWriter, r *http.Request) {
	if r.Header.Get("X-Fragment") != "1" {
		http.Redirect(w, r, "/akun-iklan", http.StatusFound)
		return
	}
	idStr := r.URL.Query().Get("id")
	id, _ := strconv.Atoi(idStr)

	ai, _ := models.GetAkunIklanByID(id)
	riwayat, _ := models.GetRiwayat("akun_iklan", id)

	tmplAkunIklan.ExecuteTemplate(w, "akun_iklan_riwayat.html", map[string]interface{}{
		"Akun":    ai,
		"Riwayat": riwayat,
	})
}

func AkunIklanRiwayatAdd(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	id, _ := strconv.Atoi(r.FormValue("entitas_id"))
	models.CreateRiwayat(&models.Riwayat{
		Entitas:   "akun_iklan",
		EntitasID: id,
		Tipe:      r.FormValue("tipe"),
		Tanggal:   r.FormValue("tanggal"),
		Catatan:   r.FormValue("catatan"),
	})
	http.Redirect(w, r, "/akun-iklan?t=riwayat", http.StatusFound)
}
