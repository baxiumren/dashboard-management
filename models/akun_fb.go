package models

import (
	"database/sql"
	"time"

	"dashboard-fb/db"
)

type AkunFB struct {
	ID           int
	Nama         string
	FBID         string
	Email        string
	Password     string
	PasswordMail string
	RecoveryMail string
	Cookie       string
	TwoFASecret  string
	Status       string
	TglBeli      string
	HargaBeli    int
	Seller       string
	Catatan      string
	CreatedAt    string
	UpdatedAt    string
}

type AkunFBStats struct {
	Total      int
	Aktif      int
	Checkpoint int
	Suspend    int
	Kunci      int
	Selfie     int
	Nonaktif   int
	Banned     int
}

func GetAllAkunFB(search, status string) ([]AkunFB, error) {
	query := `SELECT id, nama, fb_id, email, password, password_mail, recovery_mail,
		cookie, twofa_secret, status, tgl_beli, harga_beli, seller, catatan,
		created_at, updated_at FROM akun_fb WHERE 1=1`
	args := []interface{}{}

	if search != "" {
		query += ` AND (nama LIKE ? OR email LIKE ? OR fb_id LIKE ? OR seller LIKE ?)`
		s := "%" + search + "%"
		args = append(args, s, s, s, s)
	}
	if status != "" {
		query += ` AND status = ?`
		args = append(args, status)
	}
	query += ` ORDER BY created_at DESC`

	rows, err := db.DB.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []AkunFB
	for rows.Next() {
		var a AkunFB
		err := rows.Scan(&a.ID, &a.Nama, &a.FBID, &a.Email, &a.Password,
			&a.PasswordMail, &a.RecoveryMail, &a.Cookie, &a.TwoFASecret,
			&a.Status, &a.TglBeli, &a.HargaBeli, &a.Seller, &a.Catatan,
			&a.CreatedAt, &a.UpdatedAt)
		if err != nil {
			return nil, err
		}
		list = append(list, a)
	}
	return list, nil
}

func GetAkunFBByID(id int) (*AkunFB, error) {
	row := db.DB.QueryRow(`SELECT id, nama, fb_id, email, password, password_mail,
		recovery_mail, cookie, twofa_secret, status, tgl_beli, harga_beli, seller,
		catatan, created_at, updated_at FROM akun_fb WHERE id = ?`, id)
	var a AkunFB
	err := row.Scan(&a.ID, &a.Nama, &a.FBID, &a.Email, &a.Password,
		&a.PasswordMail, &a.RecoveryMail, &a.Cookie, &a.TwoFASecret,
		&a.Status, &a.TglBeli, &a.HargaBeli, &a.Seller, &a.Catatan,
		&a.CreatedAt, &a.UpdatedAt)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	return &a, err
}

func GetAllAkunFBSimple() ([]AkunFB, error) {
	rows, err := db.DB.Query(`SELECT id, nama, fb_id FROM akun_fb ORDER BY nama`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var list []AkunFB
	for rows.Next() {
		var a AkunFB
		rows.Scan(&a.ID, &a.Nama, &a.FBID)
		list = append(list, a)
	}
	return list, nil
}

func CreateAkunFB(a *AkunFB) error {
	_, err := db.DB.Exec(`INSERT INTO akun_fb
		(nama, fb_id, email, password, password_mail, recovery_mail, cookie,
		twofa_secret, status, tgl_beli, harga_beli, seller, catatan)
		VALUES (?,?,?,?,?,?,?,?,?,?,?,?,?)`,
		a.Nama, a.FBID, a.Email, a.Password, a.PasswordMail, a.RecoveryMail,
		a.Cookie, a.TwoFASecret, a.Status, a.TglBeli, a.HargaBeli, a.Seller, a.Catatan)
	return err
}

func UpdateAkunFB(a *AkunFB) error {
	_, err := db.DB.Exec(`UPDATE akun_fb SET
		nama=?, fb_id=?, email=?, password=?, password_mail=?, recovery_mail=?,
		cookie=?, twofa_secret=?, status=?, tgl_beli=?, harga_beli=?, seller=?,
		catatan=?, updated_at=? WHERE id=?`,
		a.Nama, a.FBID, a.Email, a.Password, a.PasswordMail, a.RecoveryMail,
		a.Cookie, a.TwoFASecret, a.Status, a.TglBeli, a.HargaBeli, a.Seller,
		a.Catatan, time.Now().Format("2006-01-02 15:04:05"), a.ID)
	return err
}

func DeleteAkunFB(id int) error {
	DeleteRiwayatByEntitas("akun_fb", id)
	_, err := db.DB.Exec(`DELETE FROM akun_fb WHERE id = ?`, id)
	return err
}

func GetLastInsertedAkunFBID() int {
	var id int
	db.DB.QueryRow(`SELECT id FROM akun_fb ORDER BY id DESC LIMIT 1`).Scan(&id)
	return id
}

func GetAkunFBStats() AkunFBStats {
	var s AkunFBStats
	db.DB.QueryRow(`SELECT COUNT(*) FROM akun_fb`).Scan(&s.Total)
	db.DB.QueryRow(`SELECT COUNT(*) FROM akun_fb WHERE status='aktif'`).Scan(&s.Aktif)
	db.DB.QueryRow(`SELECT COUNT(*) FROM akun_fb WHERE status='checkpoint'`).Scan(&s.Checkpoint)
	db.DB.QueryRow(`SELECT COUNT(*) FROM akun_fb WHERE status='suspend'`).Scan(&s.Suspend)
	db.DB.QueryRow(`SELECT COUNT(*) FROM akun_fb WHERE status='kunci'`).Scan(&s.Kunci)
	db.DB.QueryRow(`SELECT COUNT(*) FROM akun_fb WHERE status='selfie'`).Scan(&s.Selfie)
	db.DB.QueryRow(`SELECT COUNT(*) FROM akun_fb WHERE status='nonaktif'`).Scan(&s.Nonaktif)
	db.DB.QueryRow(`SELECT COUNT(*) FROM akun_fb WHERE status='banned'`).Scan(&s.Banned)
	return s
}
