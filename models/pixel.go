package models

import "dashboard-fb/db"

type Pixel struct {
	ID        int
	Nama      string
	PixelID   string
	BMID      int
	AkunFBID  int
	Status    string
	TglBuat   string
	Catatan   string
	CreatedAt string
	UpdatedAt string

	BMNama     string
	AkunFBNama string
}

type PixelStats struct {
	Total    int
	Aktif    int
	Nonaktif int
}

func GetAllPixel(search, status string) ([]Pixel, error) {
	query := `
		SELECT p.id, p.nama, p.pixel_id, COALESCE(p.bm_id,0), COALESCE(p.akun_fb_id,0),
		       p.status, p.tgl_buat, p.catatan, p.created_at, p.updated_at,
		       COALESCE(b.nama,''), COALESCE(a.nama,'')
		FROM pixel p
		LEFT JOIN bm b ON b.id = p.bm_id
		LEFT JOIN akun_fb a ON a.id = p.akun_fb_id
		WHERE 1=1`
	args := []interface{}{}

	if search != "" {
		query += ` AND (p.nama LIKE ? OR p.pixel_id LIKE ? OR b.nama LIKE ? OR a.nama LIKE ?)`
		s := "%" + search + "%"
		args = append(args, s, s, s, s)
	}
	if status != "" {
		query += ` AND p.status = ?`
		args = append(args, status)
	}
	query += ` ORDER BY p.created_at DESC`

	rows, err := db.DB.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []Pixel
	for rows.Next() {
		var p Pixel
		err := rows.Scan(&p.ID, &p.Nama, &p.PixelID, &p.BMID, &p.AkunFBID,
			&p.Status, &p.TglBuat, &p.Catatan, &p.CreatedAt, &p.UpdatedAt,
			&p.BMNama, &p.AkunFBNama)
		if err != nil {
			return nil, err
		}
		list = append(list, p)
	}
	return list, nil
}

func GetPixelByID(id int) (*Pixel, error) {
	row := db.DB.QueryRow(`
		SELECT p.id, p.nama, p.pixel_id, COALESCE(p.bm_id,0), COALESCE(p.akun_fb_id,0),
		       p.status, p.tgl_buat, p.catatan, p.created_at, p.updated_at,
		       COALESCE(b.nama,''), COALESCE(a.nama,'')
		FROM pixel p
		LEFT JOIN bm b ON b.id = p.bm_id
		LEFT JOIN akun_fb a ON a.id = p.akun_fb_id
		WHERE p.id = ?`, id)
	var p Pixel
	err := row.Scan(&p.ID, &p.Nama, &p.PixelID, &p.BMID, &p.AkunFBID,
		&p.Status, &p.TglBuat, &p.Catatan, &p.CreatedAt, &p.UpdatedAt,
		&p.BMNama, &p.AkunFBNama)
	if err != nil {
		return nil, nil
	}
	return &p, nil
}

func CreatePixel(p *Pixel) error {
	var bmID, akunFBID interface{}
	if p.BMID > 0 {
		bmID = p.BMID
	}
	if p.AkunFBID > 0 {
		akunFBID = p.AkunFBID
	}
	_, err := db.DB.Exec(`
		INSERT INTO pixel (nama, pixel_id, bm_id, akun_fb_id, status, tgl_buat, catatan)
		VALUES (?,?,?,?,?,?,?)`,
		p.Nama, p.PixelID, bmID, akunFBID, p.Status, p.TglBuat, p.Catatan)
	return err
}

func UpdatePixel(p *Pixel) error {
	var bmID, akunFBID interface{}
	if p.BMID > 0 {
		bmID = p.BMID
	}
	if p.AkunFBID > 0 {
		akunFBID = p.AkunFBID
	}
	_, err := db.DB.Exec(`
		UPDATE pixel SET nama=?, pixel_id=?, bm_id=?, akun_fb_id=?,
		status=?, tgl_buat=?, catatan=?, updated_at=CURRENT_TIMESTAMP
		WHERE id=?`,
		p.Nama, p.PixelID, bmID, akunFBID, p.Status, p.TglBuat, p.Catatan, p.ID)
	return err
}

func DeletePixel(id int) error {
	DeleteRiwayatByEntitas("pixel", id)
	_, err := db.DB.Exec(`DELETE FROM pixel WHERE id=?`, id)
	return err
}

func GetPixelStats() PixelStats {
	var s PixelStats
	db.DB.QueryRow(`SELECT COUNT(*) FROM pixel`).Scan(&s.Total)
	db.DB.QueryRow(`SELECT COUNT(*) FROM pixel WHERE status='aktif'`).Scan(&s.Aktif)
	db.DB.QueryRow(`SELECT COUNT(*) FROM pixel WHERE status='nonaktif'`).Scan(&s.Nonaktif)
	return s
}

func GetLastInsertedPixelID() int {
	var id int
	db.DB.QueryRow(`SELECT id FROM pixel ORDER BY id DESC LIMIT 1`).Scan(&id)
	return id
}

