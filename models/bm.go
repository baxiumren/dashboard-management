package models

import "dashboard-fb/db"

type BM struct {
	ID            int
	Nama          string
	BMID          string
	OwnerAkunFBID int
	Status        string
	TglBuat       string
	Catatan       string
	CreatedAt     string
	UpdatedAt     string

	// Join fields untuk display
	OwnerNama string
}

type BMStats struct {
	Total    int
	Aktif    int
	Disabled int
	Banned   int
}

func GetAllBM(search, status string) ([]BM, error) {
	query := `
		SELECT b.id, b.nama, b.bm_id, b.owner_akun_fb_id,
		       b.status, b.tgl_buat, b.catatan, b.created_at, b.updated_at,
		       COALESCE(a.nama,'')
		FROM bm b
		LEFT JOIN akun_fb a ON a.id = b.owner_akun_fb_id
		WHERE 1=1`
	args := []interface{}{}

	if search != "" {
		query += ` AND (b.nama LIKE ? OR b.bm_id LIKE ? OR a.nama LIKE ?)`
		s := "%" + search + "%"
		args = append(args, s, s, s)
	}
	if status != "" {
		query += ` AND b.status = ?`
		args = append(args, status)
	}
	query += ` ORDER BY b.created_at DESC`

	rows, err := db.DB.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []BM
	for rows.Next() {
		var b BM
		err := rows.Scan(&b.ID, &b.Nama, &b.BMID, &b.OwnerAkunFBID,
			&b.Status, &b.TglBuat, &b.Catatan, &b.CreatedAt, &b.UpdatedAt,
			&b.OwnerNama)
		if err != nil {
			return nil, err
		}
		list = append(list, b)
	}
	return list, nil
}

func GetBMByID(id int) (*BM, error) {
	row := db.DB.QueryRow(`
		SELECT b.id, b.nama, b.bm_id, b.owner_akun_fb_id,
		       b.status, b.tgl_buat, b.catatan, b.created_at, b.updated_at,
		       COALESCE(a.nama,'')
		FROM bm b
		LEFT JOIN akun_fb a ON a.id = b.owner_akun_fb_id
		WHERE b.id = ?`, id)
	var b BM
	err := row.Scan(&b.ID, &b.Nama, &b.BMID, &b.OwnerAkunFBID,
		&b.Status, &b.TglBuat, &b.Catatan, &b.CreatedAt, &b.UpdatedAt,
		&b.OwnerNama)
	if err != nil {
		return nil, nil
	}
	return &b, nil
}

func GetAllBMSimple() ([]BM, error) {
	rows, err := db.DB.Query(`SELECT id, nama, bm_id FROM bm ORDER BY nama`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var list []BM
	for rows.Next() {
		var b BM
		rows.Scan(&b.ID, &b.Nama, &b.BMID)
		list = append(list, b)
	}
	return list, nil
}

func CreateBM(b *BM) error {
	_, err := db.DB.Exec(`
		INSERT INTO bm (nama, bm_id, owner_akun_fb_id, status, tgl_buat, catatan)
		VALUES (?,?,?,?,?,?)`,
		b.Nama, b.BMID, b.OwnerAkunFBID, b.Status, b.TglBuat, b.Catatan)
	return err
}

func UpdateBM(b *BM) error {
	_, err := db.DB.Exec(`
		UPDATE bm SET nama=?, bm_id=?, owner_akun_fb_id=?,
		status=?, tgl_buat=?, catatan=?, updated_at=CURRENT_TIMESTAMP
		WHERE id=?`,
		b.Nama, b.BMID, b.OwnerAkunFBID, b.Status, b.TglBuat, b.Catatan, b.ID)
	return err
}

func DeleteBM(id int) error {
	DeleteRiwayatByEntitas("bm", id)
	_, err := db.DB.Exec(`DELETE FROM bm WHERE id=?`, id)
	return err
}

func GetBMStats() BMStats {
	var s BMStats
	db.DB.QueryRow(`SELECT COUNT(*) FROM bm`).Scan(&s.Total)
	db.DB.QueryRow(`SELECT COUNT(*) FROM bm WHERE status='aktif'`).Scan(&s.Aktif)
	db.DB.QueryRow(`SELECT COUNT(*) FROM bm WHERE status='disabled'`).Scan(&s.Disabled)
	db.DB.QueryRow(`SELECT COUNT(*) FROM bm WHERE status='banned'`).Scan(&s.Banned)
	return s
}

func GetLastInsertedBMID() int {
	var id int
	db.DB.QueryRow(`SELECT id FROM bm ORDER BY id DESC LIMIT 1`).Scan(&id)
	return id
}
