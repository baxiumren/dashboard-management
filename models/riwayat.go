package models

import (
	"dashboard-fb/db"
)

type Riwayat struct {
	ID        int
	Entitas   string
	EntitasID int
	Tipe      string
	Tanggal   string
	Catatan   string
	CreatedAt string
}

func GetRiwayat(entitas string, entitasID int) ([]Riwayat, error) {
	rows, err := db.DB.Query(`SELECT id, entitas, entitas_id, tipe, tanggal, catatan, created_at
		FROM riwayat WHERE entitas=? AND entitas_id=? ORDER BY tanggal DESC, created_at DESC`,
		entitas, entitasID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var list []Riwayat
	for rows.Next() {
		var r Riwayat
		rows.Scan(&r.ID, &r.Entitas, &r.EntitasID, &r.Tipe, &r.Tanggal, &r.Catatan, &r.CreatedAt)
		list = append(list, r)
	}
	return list, nil
}

func CreateRiwayat(r *Riwayat) error {
	_, err := db.DB.Exec(`INSERT INTO riwayat (entitas, entitas_id, tipe, tanggal, catatan)
		VALUES (?,?,?,?,?)`, r.Entitas, r.EntitasID, r.Tipe, r.Tanggal, r.Catatan)
	return err
}

func DeleteRiwayat(id int) error {
	_, err := db.DB.Exec(`DELETE FROM riwayat WHERE id=?`, id)
	return err
}

func DeleteRiwayatByEntitas(entitas string, entitasID int) {
	db.DB.Exec(`DELETE FROM riwayat WHERE entitas=? AND entitas_id=?`, entitas, entitasID)
}
