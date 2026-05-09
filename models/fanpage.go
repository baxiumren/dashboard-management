package models

import (
	"dashboard-fb/db"
)

type Fanpage struct {
	ID        int
	Nama      string
	PageID    string
	AkunFBID  int
	BMID      int
	Status    string
	TglBuat   string
	Catatan   string
	CreatedAt string
	UpdatedAt string

	// Join fields untuk display
	AkunFBNama string
	BMNama     string
}

type FanpageStats struct {
	Total    int
	Aktif    int
	Nonaktif int
	Banned   int
}

func GetAllFanpage(search, status string) ([]Fanpage, error) {
	query := `
		SELECT f.id, f.nama, f.page_id, f.akun_fb_id, COALESCE(f.bm_id,0),
		       f.status, f.tgl_buat, f.catatan, f.created_at, f.updated_at,
		       COALESCE(a.nama,''), COALESCE(b.nama,'')
		FROM fanpage f
		LEFT JOIN akun_fb a ON a.id = f.akun_fb_id
		LEFT JOIN bm b ON b.id = f.bm_id
		WHERE 1=1`
	args := []interface{}{}

	if search != "" {
		query += ` AND (f.nama LIKE ? OR f.page_id LIKE ? OR a.nama LIKE ? OR a.fb_id LIKE ?)`
		s := "%" + search + "%"
		args = append(args, s, s, s, s)
	}
	if status != "" {
		query += ` AND f.status = ?`
		args = append(args, status)
	}
	query += ` ORDER BY f.created_at DESC`

	rows, err := db.DB.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []Fanpage
	for rows.Next() {
		var fp Fanpage
		err := rows.Scan(&fp.ID, &fp.Nama, &fp.PageID, &fp.AkunFBID, &fp.BMID,
			&fp.Status, &fp.TglBuat, &fp.Catatan, &fp.CreatedAt, &fp.UpdatedAt,
			&fp.AkunFBNama, &fp.BMNama)
		if err != nil {
			return nil, err
		}
		list = append(list, fp)
	}
	return list, nil
}

func GetFanpageByID(id int) (*Fanpage, error) {
	row := db.DB.QueryRow(`
		SELECT f.id, f.nama, f.page_id, f.akun_fb_id, COALESCE(f.bm_id,0),
		       f.status, f.tgl_buat, f.catatan, f.created_at, f.updated_at,
		       COALESCE(a.nama,''), COALESCE(b.nama,'')
		FROM fanpage f
		LEFT JOIN akun_fb a ON a.id = f.akun_fb_id
		LEFT JOIN bm b ON b.id = f.bm_id
		WHERE f.id = ?`, id)
	var fp Fanpage
	err := row.Scan(&fp.ID, &fp.Nama, &fp.PageID, &fp.AkunFBID, &fp.BMID,
		&fp.Status, &fp.TglBuat, &fp.Catatan, &fp.CreatedAt, &fp.UpdatedAt,
		&fp.AkunFBNama, &fp.BMNama)
	if err != nil {
		return nil, nil
	}
	return &fp, nil
}

func CreateFanpage(fp *Fanpage) error {
	var bmID interface{}
	if fp.BMID > 0 {
		bmID = fp.BMID
	}
	_, err := db.DB.Exec(`
		INSERT INTO fanpage (nama, page_id, akun_fb_id, bm_id, status, tgl_buat, catatan)
		VALUES (?,?,?,?,?,?,?)`,
		fp.Nama, fp.PageID, fp.AkunFBID, bmID, fp.Status, fp.TglBuat, fp.Catatan)
	return err
}

func UpdateFanpage(fp *Fanpage) error {
	var bmID interface{}
	if fp.BMID > 0 {
		bmID = fp.BMID
	}
	_, err := db.DB.Exec(`
		UPDATE fanpage SET nama=?, page_id=?, akun_fb_id=?, bm_id=?,
		status=?, tgl_buat=?, catatan=?, updated_at=CURRENT_TIMESTAMP
		WHERE id=?`,
		fp.Nama, fp.PageID, fp.AkunFBID, bmID, fp.Status, fp.TglBuat, fp.Catatan, fp.ID)
	return err
}

func DeleteFanpage(id int) error {
	DeleteRiwayatByEntitas("fanpage", id)
	_, err := db.DB.Exec(`DELETE FROM fanpage WHERE id=?`, id)
	return err
}

func GetFanpageStats() FanpageStats {
	var s FanpageStats
	db.DB.QueryRow(`SELECT COUNT(*) FROM fanpage`).Scan(&s.Total)
	db.DB.QueryRow(`SELECT COUNT(*) FROM fanpage WHERE status='aktif'`).Scan(&s.Aktif)
	db.DB.QueryRow(`SELECT COUNT(*) FROM fanpage WHERE status='nonaktif'`).Scan(&s.Nonaktif)
	db.DB.QueryRow(`SELECT COUNT(*) FROM fanpage WHERE status='banned'`).Scan(&s.Banned)
	return s
}

func GetLastInsertedFanpageID() int {
	var id int
	db.DB.QueryRow(`SELECT id FROM fanpage ORDER BY id DESC LIMIT 1`).Scan(&id)
	return id
}
