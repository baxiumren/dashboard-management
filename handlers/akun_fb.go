package handlers

import (
	"fmt"
	"html/template"
	"net/http"
	"strconv"
	"strings"

	"dashboard-fb/models"
)

var tmplAkunFB *template.Template

func InitAkunFB(tmpl *template.Template) {
	tmplAkunFB = tmpl
}

func AkunFBList(w http.ResponseWriter, r *http.Request,
	decryptFn func(string) (string, error),
	encryptFn func(string) (string, error),
	isDefaultPassword bool) {

	search := r.URL.Query().Get("search")
	status := r.URL.Query().Get("status")

	list, err := models.GetAllAkunFB(search, status)
	if err != nil {
		http.Error(w, "Error: "+err.Error(), 500)
		return
	}

	// Decrypt sensitive fields for display
	for i := range list {
		list[i].Password, _ = decryptFn(list[i].Password)
		list[i].PasswordMail, _ = decryptFn(list[i].PasswordMail)
		list[i].Cookie, _ = decryptFn(list[i].Cookie)
		list[i].TwoFASecret, _ = decryptFn(list[i].TwoFASecret)
	}

	stats := models.GetAkunFBStats()

	tmplAkunFB.ExecuteTemplate(w, "akun_fb.html", map[string]interface{}{
		"List":         list,
		"Stats":        stats,
		"Search":       search,
		"StatusFilter": status,
		"Page":         "akun-fb",
		"Welcome":      isDefaultPassword,
	})
}

func AkunFBAdd(w http.ResponseWriter, r *http.Request, encryptFn func(string) (string, error)) {
	r.ParseForm()

	harga, _ := strconv.Atoi(r.FormValue("harga_beli"))
	pass, _ := encryptFn(r.FormValue("password"))
	passMail, _ := encryptFn(r.FormValue("password_mail"))
	cookie, _ := encryptFn(r.FormValue("cookie"))
	twofa, _ := encryptFn(r.FormValue("twofa_secret"))

	a := &models.AkunFB{
		Nama:         r.FormValue("nama"),
		FBID:         r.FormValue("fb_id"),
		Email:        r.FormValue("email"),
		Password:     pass,
		PasswordMail: passMail,
		RecoveryMail: r.FormValue("recovery_mail"),
		Cookie:       cookie,
		TwoFASecret:  twofa,
		Status:       r.FormValue("status"),
		TglBeli:      r.FormValue("tgl_beli"),
		HargaBeli:    harga,
		Seller:       r.FormValue("seller"),
		Catatan:      r.FormValue("catatan"),
	}

	if err := models.CreateAkunFB(a); err != nil {
		http.Error(w, "Gagal simpan: "+err.Error(), 500)
		return
	}

	// Auto-create riwayat "beli"
	if a.TglBeli != "" {
		var lastID int
		models.GetAkunFBStats() // just to ensure db is ready
		row := models.GetLastInsertedAkunFBID()
		if row > 0 {
			models.CreateRiwayat(&models.Riwayat{
				Entitas:   "akun_fb",
				EntitasID: row,
				Tipe:      "beli",
				Tanggal:   a.TglBeli,
				Catatan:   fmt.Sprintf("Beli dari %s, harga Rp %d", a.Seller, harga),
			})
			_ = lastID
		}
	}

	http.Redirect(w, r, "/akun-fb?t=tambah", http.StatusFound)
}

func AkunFBEdit(w http.ResponseWriter, r *http.Request, encryptFn func(string) (string, error), decryptFn func(string) (string, error)) {
	idStr := r.URL.Query().Get("id")
	id, _ := strconv.Atoi(idStr)

	existing, err := models.GetAkunFBByID(id)
	if err != nil || existing == nil {
		http.Redirect(w, r, "/akun-fb", http.StatusFound)
		return
	}

	r.ParseForm()
	harga, _ := strconv.Atoi(r.FormValue("harga_beli"))

	// Only re-encrypt if new value provided, else keep existing
	pass := existing.Password
	if v := r.FormValue("password"); v != "" {
		pass, _ = encryptFn(v)
	}
	passMail := existing.PasswordMail
	if v := r.FormValue("password_mail"); v != "" {
		passMail, _ = encryptFn(v)
	}
	cookie := existing.Cookie
	if v := r.FormValue("cookie"); v != "" {
		cookie, _ = encryptFn(v)
	}
	twofa := existing.TwoFASecret
	if v := r.FormValue("twofa_secret"); v != "" {
		twofa, _ = encryptFn(v)
	}

	existing.Nama = r.FormValue("nama")
	existing.FBID = r.FormValue("fb_id")
	existing.Email = r.FormValue("email")
	existing.Password = pass
	existing.PasswordMail = passMail
	existing.RecoveryMail = r.FormValue("recovery_mail")
	existing.Cookie = cookie
	existing.TwoFASecret = twofa
	existing.Status = r.FormValue("status")
	existing.TglBeli = r.FormValue("tgl_beli")
	existing.HargaBeli = harga
	existing.Seller = r.FormValue("seller")
	existing.Catatan = r.FormValue("catatan")

	if err := models.UpdateAkunFB(existing); err != nil {
		http.Error(w, "Gagal update: "+err.Error(), 500)
		return
	}
	http.Redirect(w, r, "/akun-fb?t=edit", http.StatusFound)
}

func AkunFBDelete(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Query().Get("id")
	id, _ := strconv.Atoi(idStr)
	models.DeleteAkunFB(id)
	http.Redirect(w, r, "/akun-fb?t=hapus", http.StatusFound)
}

func AkunFBImportBulk(w http.ResponseWriter, r *http.Request, encryptFn func(string) (string, error)) {
	r.ParseForm()
	raw := r.FormValue("raw_text")

	parsed := parseBulkText(raw)

	for _, a := range parsed {
		pass, _ := encryptFn(a.Password)
		passMail, _ := encryptFn(a.PasswordMail)
		cookie, _ := encryptFn(a.Cookie)
		twofa, _ := encryptFn(a.TwoFASecret)
		a.Password = pass
		a.PasswordMail = passMail
		a.Cookie = cookie
		a.TwoFASecret = twofa
		a.Status = "aktif"
		models.CreateAkunFB(&a)
	}

	http.Redirect(w, r, "/akun-fb?t=import", http.StatusFound)
}

// parseBulkText detects format (TSV dari Sheets/Excel atau teks seller) lalu parse.
func parseBulkText(raw string) []models.AkunFB {
	for _, line := range strings.Split(raw, "\n") {
		if strings.Contains(line, "\t") {
			return parseTSV(raw)
		}
	}
	return parseSellerText(raw)
}

// parseTSV: copy-paste dari Google Sheets / Excel.
// Kolom: Nama | FB ID | Password | 2FA | Email | Password Email | Recovery Email
func parseTSV(raw string) []models.AkunFB {
	var results []models.AkunFB
	headerKeywords := map[string]bool{
		"nama fb": true, "nama": true, "name": true,
		"id": true, "fb id": true, "fb sedang jalan": true,
	}

	get := func(cols []string, i int) string {
		if i < len(cols) {
			return strings.TrimSpace(cols[i])
		}
		return ""
	}

	for _, line := range strings.Split(raw, "\n") {
		line = strings.TrimRight(line, "\r")
		if strings.TrimSpace(line) == "" {
			continue
		}
		cols := strings.Split(line, "\t")

		// Lewati baris header
		first := strings.ToLower(strings.TrimSpace(cols[0]))
		if headerKeywords[first] {
			continue
		}

		fbID := get(cols, 1)
		if fbID == "" {
			continue
		}

		results = append(results, models.AkunFB{
			Nama:         get(cols, 0),
			FBID:         fbID,
			Password:     get(cols, 2),
			TwoFASecret:  get(cols, 3),
			Email:        get(cols, 4),
			PasswordMail: get(cols, 5),
			RecoveryMail: get(cols, 6),
		})
	}
	return results
}

// parseSellerText: format teks dari seller (key : value per baris).
func parseSellerText(raw string) []models.AkunFB {
	var results []models.AkunFB
	var current models.AkunFB
	inBlock := false

	for _, line := range strings.Split(raw, "\n") {
		line = strings.TrimSpace(line)
		if line == "" {
			if inBlock && current.FBID != "" {
				results = append(results, current)
				current = models.AkunFB{}
				inBlock = false
			}
			continue
		}

		// Lewati nomor urut "1.", "2.", dst.
		if len(line) <= 3 && strings.HasSuffix(line, ".") {
			continue
		}

		parts := strings.SplitN(line, ":", 2)
		if len(parts) != 2 {
			continue
		}
		key := strings.ToLower(strings.TrimSpace(parts[0]))
		val := strings.TrimSpace(parts[1])
		inBlock = true

		switch key {
		case "id":
			current.FBID = val
		case "password":
			current.Password = val
		case "email":
			current.Email = val
		case "passwordmail":
			current.PasswordMail = val
		case "recoverymail":
			current.RecoveryMail = val
		case "cookie":
			current.Cookie = val
		case "2fa", "twofa", "2fa secret":
			current.TwoFASecret = val
		}
	}
	if inBlock && current.FBID != "" {
		results = append(results, current)
	}
	return results
}

// AkunFBRiwayat returns riwayat for a given akun_fb id
func AkunFBRiwayat(w http.ResponseWriter, r *http.Request) {
	if r.Header.Get("X-Fragment") != "1" {
		http.Redirect(w, r, "/akun-fb", http.StatusFound)
		return
	}
	idStr := r.URL.Query().Get("id")
	id, _ := strconv.Atoi(idStr)

	akun, _ := models.GetAkunFBByID(id)
	riwayat, _ := models.GetRiwayat("akun_fb", id)

	tmplAkunFB.ExecuteTemplate(w, "riwayat.html", map[string]interface{}{
		"Akun":    akun,
		"Riwayat": riwayat,
	})
}

func AkunFBRiwayatAdd(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	id, _ := strconv.Atoi(r.FormValue("entitas_id"))
	models.CreateRiwayat(&models.Riwayat{
		Entitas:   "akun_fb",
		EntitasID: id,
		Tipe:      r.FormValue("tipe"),
		Tanggal:   r.FormValue("tanggal"),
		Catatan:   r.FormValue("catatan"),
	})
	http.Redirect(w, r, "/akun-fb?t=riwayat", http.StatusFound)
}
