package user

import (
	"database/sql"
	"errors"

	"examples.com/auth-service/hashing"
)

func LoginCheck(db *sql.DB, login string) (bool, error) {
	query := "SELECT * FROM user WHERE login = ?"

	rows, err := db.Query(query, login)
	if err != nil {
		return false, err
	}
	defer rows.Close()
	if rows.Next() {
		return true, nil
	} else {
		return false, nil
	}
}

func Registration(db *sql.DB, login string, password string) (bool, error) {
	query := "INSERT INTO user (login, password) VALUES (?, ?)"
	tx, err := db.Begin()
	if err != nil {
		return false, err
	}
	_, err = tx.Exec(query, login, password)
	if err != nil {
		tx.Rollback()
		return false, err
	}
	err = tx.Commit()
	if err != nil {
		return false, nil
	}
	return true, nil

}

func Login(db *sql.DB, login string, password string) (bool, error) {
	errLogin := errors.New("wrong password")
	var tempPass string
	query := "SELECT password FROM user WHERE login = ?"

	tx, err := db.Begin()
	if err != nil {
		return false, err
	}
	rows, err := tx.Query(query, login)
	if err != nil {
		tx.Rollback()
		return false, err
	}
	_ = rows.Next()
	err = rows.Scan(&tempPass)
	if err != nil {
		return false, err
	}
	defer rows.Close()
	tx.Commit()
	checkPass := hashing.CheckPasswordHash(password, tempPass)
	if checkPass {
		return true, nil
	} else {
		return false, errLogin
	}
}
