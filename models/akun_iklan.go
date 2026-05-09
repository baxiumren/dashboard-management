package models

import "dashboard-fb/db"

type AkunIklan struct {
	ID              int
	Nama            string
	AdAccountID     string
	BMID            int
	RunnerAkunFBID  int
	LimitHarian     int
	LimitTotal      int
	MetodeBayar     string
	MataUang        string
	Status          string
	TglBuat         string
	Catatan         string
	CreatedAt       string
	UpdatedAt       string

	// Join fields untuk display
	BMNama     string
	RunnerNama string
}

type AkunIklanStats struct {
	Total    int
	Aktif    int
	Disabled int
	Banned   int
}

func GetAllAkunIklan(search, status string) ([]AkunIklan, error) {
	query := `
		SELECT ai.id, ai.nama, ai.ad_account_id, ai.bm_id,
		       COALESCE(ai.runner_akun_fb_id,0),
		       ai.limit_harian, ai.limit_total,
		       ai.metode_bayar, ai.mata_uang,
		       ai.status, ai.tgl_buat, ai.catatan,
		       ai.created_at, ai.updated_at,
		       COALESCE(b.nama,''), COALESCE(a.nama,'')
		FROM akun_iklan ai
		LEFT JOIN bm b ON b.id = ai.bm_id
		LEFT JOIN akun_fb a ON a.id = ai.runner_akun_fb_id
		WHERE 1=1`
	args := []interface{}{}

	if search != "" {
		query += ` AND (ai.nama LIKE ? OR ai.ad_account_id LIKE ? OR b.nama LIKE ? OR a.nama LIKE ?)`
		s := "%" + search + "%"
		args = append(args, s, s, s, s)
	}
	if status != "" {
		query += ` AND ai.status = ?`
		args = append(args, status)
	}
	query += ` ORDER BY ai.created_at DESC`

	rows, err := db.DB.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []AkunIklan
	for rows.Next() {
		var ai AkunIklan
		err := rows.Scan(
			&ai.ID, &ai.Nama, &ai.AdAccountID, &ai.BMID,
			&ai.RunnerAkunFBID, &ai.LimitHarian, &ai.LimitTotal,
			&ai.MetodeBayar, &ai.MataUang,
			&ai.Status, &ai.TglBuat, &ai.Catatan,
			&ai.CreatedAt, &ai.UpdatedAt,
			&ai.BMNama, &ai.RunnerNama,
		)
		if err != nil {
			return nil, err
		}
		list = append(list, ai)
	}
	return list, nil
}

func GetAkunIklanByID(id int) (*AkunIklan, error) {
	row := db.DB.QueryRow(`
		SELECT ai.id, ai.nama, ai.ad_account_id, ai.bm_id,
		       COALESCE(ai.runner_akun_fb_id,0),
		       ai.limit_harian, ai.limit_total,
		       ai.metode_bayar, ai.mata_uang,
		       ai.status, ai.tgl_buat, ai.catatan,
		       ai.created_at, ai.updated_at,
		       COALESCE(b.nama,''), COALESCE(a.nama,'')
		FROM akun_iklan ai
		LEFT JOIN bm b ON b.id = ai.bm_id
		LEFT JOIN akun_fb a ON a.id = ai.runner_akun_fb_id
		WHERE ai.id = ?`, id)
	var ai AkunIklan
	err := row.Scan(
		&ai.ID, &ai.Nama, &ai.AdAccountID, &ai.BMID,
		&ai.RunnerAkunFBID, &ai.LimitHarian, &ai.LimitTotal,
		&ai.MetodeBayar, &ai.MataUang,
		&ai.Status, &ai.TglBuat, &ai.Catatan,
		&ai.CreatedAt, &ai.UpdatedAt,
		&ai.BMNama, &ai.RunnerNama,
	)
	if err != nil {
		return nil, nil
	}
	return &ai, nil
}

func CreateAkunIklan(ai *AkunIklan) error {
	var runnerID interface{}
	if ai.RunnerAkunFBID > 0 {
		runnerID = ai.RunnerAkunFBID
	}
	_, err := db.DB.Exec(`
		INSERT INTO akun_iklan
		(nama, ad_account_id, bm_id, runner_akun_fb_id,
		 limit_harian, limit_total, metode_bayar, mata_uang,
		 status, tgl_buat, catatan)
		VALUES (?,?,?,?,?,?,?,?,?,?,?)`,
		ai.Nama, ai.AdAccountID, ai.BMID, runnerID,
		ai.LimitHarian, ai.LimitTotal, ai.MetodeBayar, ai.MataUang,
		ai.Status, ai.TglBuat, ai.Catatan)
	return err
}

func UpdateAkunIklan(ai *AkunIklan) error {
	var runnerID interface{}
	if ai.RunnerAkunFBID > 0 {
		runnerID = ai.RunnerAkunFBID
	}
	_, err := db.DB.Exec(`
		UPDATE akun_iklan SET
		nama=?, ad_account_id=?, bm_id=?, runner_akun_fb_id=?,
		limit_harian=?, limit_total=?, metode_bayar=?, mata_uang=?,
		status=?, tgl_buat=?, catatan=?, updated_at=CURRENT_TIMESTAMP
		WHERE id=?`,
		ai.Nama, ai.AdAccountID, ai.BMID, runnerID,
		ai.LimitHarian, ai.LimitTotal, ai.MetodeBayar, ai.MataUang,
		ai.Status, ai.TglBuat, ai.Catatan, ai.ID)
	return err
}

func DeleteAkunIklan(id int) error {
	DeleteRiwayatByEntitas("akun_iklan", id)
	_, err := db.DB.Exec(`DELETE FROM akun_iklan WHERE id=?`, id)
	return err
}

func GetAkunIklanStats() AkunIklanStats {
	var s AkunIklanStats
	db.DB.QueryRow(`SELECT COUNT(*) FROM akun_iklan`).Scan(&s.Total)
	db.DB.QueryRow(`SELECT COUNT(*) FROM akun_iklan WHERE status='aktif'`).Scan(&s.Aktif)
	db.DB.QueryRow(`SELECT COUNT(*) FROM akun_iklan WHERE status='disabled'`).Scan(&s.Disabled)
	db.DB.QueryRow(`SELECT COUNT(*) FROM akun_iklan WHERE status='banned'`).Scan(&s.Banned)
	return s
}

func GetLastInsertedAkunIklanID() int {
	var id int
	db.DB.QueryRow(`SELECT id FROM akun_iklan ORDER BY id DESC LIMIT 1`).Scan(&id)
	return id
}
