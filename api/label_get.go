package main

import (
	"database/sql"
)

func commentGetLabelsHex(commentHex string) ([]string, error) {
	if commentHex == "" {
		return nil, errorMissingField
	}

	statement := `
		SELECT labelHex
		FROM comments_labels
		WHERE commentHex=$1
	`

	var rows *sql.Rows
	var err error

	rows, err = db.Query(statement, commentHex)
	if err != nil {
		logger.Errorf("cannot get comment's label: %v", err)
		return nil, errorInternal
	}

	labels := []string{}
	for rows.Next() {
		var label string
		if err = rows.Scan(&label); err != nil {
			return nil, errorInternal
		}
		labels = append(labels, label)
	}

	return labels, nil
}
