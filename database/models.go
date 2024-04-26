package database

import (
	"database/sql"

	_ "github.com/mattn/go-sqlite3"
)

type PWItem struct {
	ID       int
	Name     string
	Email    string
	Password []byte
}

var db *sql.DB

func InitDB(dsn string) error {
	var err error
	db, err = sql.Open("sqlite3", dsn)

	if err != nil {
		return err
	}

	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS pwitems(
		id INTEGER PRIMARY KEY,
		name TEXT NOT NULL,
		email TEXT NOT NULL,
		password BLOB NOT NULL)`)

	if err != nil {
		return err
	}

	return nil
}

func AddNewPassword(name string, email string, password []byte) (int, error) {
	res, err := db.Exec("INSERT INTO pwitems(name, email, password) VALUES (?1, ?2 , ?3)", name, email, password)
	if err != nil {
		return 0, err
	}
	id, err := res.LastInsertId()
	if err != nil {
		return 0, err
	}
	return int(id), nil
}

func ReadAllPasswords() ([]PWItem, error) {
	rows, err := db.Query("SELECT * FROM pwitems ORDER BY name")
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	pws := []PWItem{}
	for rows.Next() {
		pw := PWItem{}
		err := rows.Scan(&pw.ID, &pw.Name, &pw.Email, &pw.Password)
		if err != nil {
			return nil, err
		}
		pws = append(pws, pw)
	}

	err = rows.Err()
	if err != nil {
		return nil, err
	}

	if len(pws) == 0 {
		return nil, sql.ErrNoRows
	}

	return pws, nil
}

func ReadOnePassword(id int) (PWItem, error) {
	var p PWItem
	err := db.QueryRow("SELECT * FROM pwitems WHERE id = ?1", id).Scan(&p.ID, &p.Name, &p.Email, &p.Password)
	return p, err
}

func UpdateOnePassword(id int, name string, email string, password []byte) error {
	_, err := db.Exec("UPDATE pwitems SET name=?1, email=?2, password=?3 WHERE id=?4", name, email, password, id)
	return err
}

func DeletePassword(id int) error {
	_, err := db.Exec("DELETE FROM pwitems WHERE id=?1", id)
	return err
}
