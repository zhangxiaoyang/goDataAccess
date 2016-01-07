package util

import (
	"database/sql"
	"fmt"
)

func InitTable(initSql string, dbFilePath string) (*sql.DB, error) {
	db, err := sql.Open("sqlite3", dbFilePath)
	if err != nil {
		return nil, err
	}
	_, err = db.Exec(initSql)
	if err != nil {
		return nil, err
	}
	return db, nil
}

func GetLastLevel(tableName string, db *sql.DB) int {
	level := 0
	rows, err := db.Query(fmt.Sprintf(
		"SELECT level FROM %s ORDER BY level DESC LIMIT 1", tableName,
	))
	if err == nil {
		defer rows.Close()
		for rows.Next() {
			rows.Scan(&level)
			break
		}
	}
	return level
}
