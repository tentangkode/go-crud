package models

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/jeypc/go-crud/config"
	"github.com/jeypc/go-crud/entities"
)

type PasienModel struct {
	conn *sql.DB
}

func NewPasienModel() *PasienModel {
	conn, err := config.DBConnection()
	if err != nil {
		panic(err)
	}

	return &PasienModel{
		conn: conn,
	}
}

func (p *PasienModel) FindAll() ([]entities.Pasien, error) {

	rows, err := p.conn.Query("select * from pasien")
	if err != nil {
		return []entities.Pasien{}, err
	}
	defer rows.Close()

	var dataPasien []entities.Pasien
	for rows.Next() {
		var pasien entities.Pasien
		rows.Scan(&pasien.Id,
			&pasien.NamaLengkap,
			&pasien.NIK,
			&pasien.JenisKelamin,
			&pasien.TempatLahir,
			&pasien.TanggalLahir,
			&pasien.Alamat,
			&pasien.NoHp)

		if pasien.JenisKelamin == "1" {
			pasien.JenisKelamin = "Laki-laki"
		} else {
			pasien.JenisKelamin = "Perempuan"
		}
		// 2006-01-02 => yyyy-mm-dd
		tgl_lahir, _ := time.Parse("2006-01-02", pasien.TanggalLahir)
		// 02-01-2006 => dd-mm-yyyy
		pasien.TanggalLahir = tgl_lahir.Format("02-01-2006")

		dataPasien = append(dataPasien, pasien)
	}

	return dataPasien, nil

}

func (p *PasienModel) Create(pasien entities.Pasien) bool {

	result, err := p.conn.Exec("insert into pasien (nama_lengkap, nik, jenis_kelamin, tempat_lahir, tanggal_lahir, alamat, no_hp) values(?,?,?,?,?,?,?)",
		pasien.NamaLengkap, pasien.NIK, pasien.JenisKelamin, pasien.TempatLahir, pasien.TanggalLahir, pasien.Alamat, pasien.NoHp)

	if err != nil {
		fmt.Println(err)
		return false
	}

	lastInsertId, _ := result.LastInsertId()

	return lastInsertId > 0
}

func (p *PasienModel) Find(id int64, pasien *entities.Pasien) error {

	return p.conn.QueryRow("select * from pasien where id = ?", id).Scan(
		&pasien.Id,
		&pasien.NamaLengkap,
		&pasien.NIK,
		&pasien.JenisKelamin,
		&pasien.TempatLahir,
		&pasien.TanggalLahir,
		&pasien.Alamat,
		&pasien.NoHp)
}

func (p *PasienModel) Update(pasien entities.Pasien) error {

	_, err := p.conn.Exec(
		"update pasien set nama_lengkap = ?, nik = ?, jenis_kelamin = ?, tempat_lahir = ?, tanggal_lahir = ?, alamat = ?, no_hp = ? where id = ?",
		pasien.NamaLengkap, pasien.NIK, pasien.JenisKelamin, pasien.TempatLahir, pasien.TanggalLahir, pasien.Alamat, pasien.NoHp, pasien.Id)

	if err != nil {
		return err
	}

	return nil
}

func (p *PasienModel) Delete(id int64) {
	p.conn.Exec("delete from pasien where id = ?", id)
}
