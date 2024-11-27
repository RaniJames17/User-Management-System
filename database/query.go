package database

import (
	"database/sql"
	"log"
)

func ExecuteQuery(query string, args ...interface{}) (sql.Result, error) {
	stmt, err := DB.Prepare(query)
	if err != nil {
		log.Printf("Error preparing query: %v", err)
		return nil, err
	}
	defer stmt.Close()

	result, err := stmt.Exec(args...)
	if err != nil {
		log.Printf("Error executing query: %v", err)
		return nil, err
	}

	return result, nil
}
