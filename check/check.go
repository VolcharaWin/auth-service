package check

import (
	"database/sql"
)

func LoginCheck(db *sql.DB, login string, password string) (bool, error) {
	query := "SELECT * FROM user WHERE login = ?"

	rows, err := db.Query(query, login)
	if err != nil {
		return false, err
	}
	if rows.Next() {
		return true, nil
	} else {
		return false, nil
	}
}
