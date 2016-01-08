package util

import (
	"database/sql"
	"fmt"
	"log"
	"regexp"
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

func GetLastProxies(tableName string, db *sql.DB) ([]string, int) {
	proxies := []string{}
	level := GetLastLevel(tableName, db)

	rows, err := db.Query(fmt.Sprintf(
		"SELECT ip, port FROM %s WHERE level=%d",
		tableName,
		level,
	))
	if err == nil {
		defer rows.Close()
		for rows.Next() {
			var ip string
			var port string
			err := rows.Scan(&ip, &port)
			if err == nil {
				proxies = append(proxies, fmt.Sprintf("%s:%s", ip, port))
			}
		}
	} else {
		log.Fatal(err)
	}
	return proxies, level
}

func GetLastProxiesByDomain(tableName string, domain string, db *sql.DB) ([]string, string, int) {
	proxies := []string{}
	level := GetLastLevel(tableName, db)

	var rows *sql.Rows
	var err error
	if domain == "" {
		rows, err = db.Query(fmt.Sprintf(
			"SELECT ip, port, domain FROM %s WHERE level=%d",
			tableName,
			level,
		))
	} else {
		rows, err = db.Query(fmt.Sprintf(
			"SELECT ip, port, domain FROM %s WHERE level=%d AND domain='%s'",
			tableName,
			level,
			domain,
		))
	}
	if err == nil {
		defer rows.Close()
		for rows.Next() {
			var ip string
			var port string
			err := rows.Scan(&ip, &port, &domain)
			if err == nil {
				proxies = append(proxies, fmt.Sprintf("%s:%s", ip, port))
			}
		}
	} else {
		log.Fatal(err)
	}
	return proxies, domain, level
}

func ExtractDomain(url string) string {
	return regexp.MustCompile(`http[s]?://([\w\-\.]+)`).FindStringSubmatch(url)[1]
}
