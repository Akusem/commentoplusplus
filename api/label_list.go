package main

import (
	"database/sql"
	"net/http"
)

func labelListAll(domain string) ([]label, error) {
	if domain == "" {
		return nil, errorMissingField
	}

	statement := `
		SELECT
			labelHex,
			color,
			domain
		FROM labels
		WHERE
			canon(domain) LIKE canon($1);
	`

	var rows *sql.Rows
	var err error

	rows, err = db.Query(statement, domain)

	if err != nil {
		logger.Errorf("cannot get labels: %v", err)
		return nil, errorInternal
	}
	defer rows.Close()

	labels := []label{}
	for rows.Next() {
		l := label{}
		if err = rows.Scan(
			&l.LabelHex,
			&l.Color,
			&l.Domain,
		); err != nil {
			return nil, errorInternal
		}

		labels = append(labels, l)
	}

	return labels, nil
}

func labelListAllHandler(w http.ResponseWriter, r *http.Request) {
	type request struct {
		OwnerToken *string `json:"ownerToken"`
		Domain     *string `json:"domain"`
	}

	var x request
	if err := bodyUnmarshal(r, &x); err != nil {
		bodyMarshal(w, response{"success": false, "message": err.Error()})
		return
	}

	o, err := ownerGetByOwnerToken(*x.OwnerToken)
	if err != nil {
		bodyMarshal(w, response{"success": false, "message": err.Error()})
		return
	}

	domain := domainStrip(*x.Domain)
	isOwner, err := domainOwnershipVerify(o.OwnerHex, domain)
	if err != nil {
		bodyMarshal(w, response{"success": false, "message": err.Error()})
		return
	}

	if !isOwner {
		bodyMarshal(w, response{"success": false, "message": errorNotAuthorised.Error()})
		return
	}

	labels, err := labelListAll(domain)
	if err != nil {
		bodyMarshal(w, response{"success": false, "message": err.Error()})
		return
	}

	bodyMarshal(w, response{
		"success": true,
		"domain":  domain,
		"labels":  labels,
	})
}
