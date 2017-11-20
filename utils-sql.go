package main

import (
	"database/sql"
	"log"
)

func sqlExecOrFail(db *sql.DB, query string, parameters ...string) int64 {
	result, err := db.Exec(query, arrayStringToInterface(parameters)...)
	if err != nil {
		log.Fatal(err)
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		log.Fatal(err)
	}
	return rowsAffected
}

func sqlStringArrayOrFail(db *sql.DB, query string, parameters ...string) []string {
	rows, err := db.Query(query, arrayStringToInterface(parameters)...)
	defer rows.Close()
	var arr []string
	for rows.Next() {
		var item string
		err = rows.Scan(&item)
		if err != nil {
			log.Fatal(err)
		}
		arr = append(arr, item)
	}
	return arr
}
