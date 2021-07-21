package main

import (
	"net/http"
)

func labelNew(name string, color string, domain string) (string, error) {
	// Test for missing field
	if name == "" || color == "" || domain == "" {
		return "", errorMissingField
	}

	// Define primary key
	labelHex, err := randomHex(32)
	if err != nil {
		return "", err
	}

	// Create label in db
	statement := `
		INSERT INTO
		labels (labelHex, name, color, domain)
		VALUES   ($1,      $2,   $3,    $4);
	`
	_, err = db.Exec(statement, labelHex, name, color, domain)
	if err != nil {
		logger.Errorf("cannot insert comment: %v", err)
		return "", errorInternal
	}

	return labelHex, nil
}

func labelNewHandler(w http.ResponseWriter, r *http.Request) {
	type request struct {
		OwnerToken *string `json:"ownerToken"`
		Domain     *string `json:"domain"`
		Name       *string `json:"name"`
		Color      *string `json:"color"`
	}

	// Extract request
	var x request
	if err := bodyUnmarshal(r, &x); err != nil {
		bodyMarshal(w, response{"success": false, "message": err.Error()})
		return
	}

	domain := domainStrip(*x.Domain)

	// Verify domain existence
	_, err := domainGet(domain)
	if err != nil {
		bodyMarshal(w, response{"success": false, "message": err.Error()})
		return
	}

	// Verify user is owner of the Domain
	o, err := ownerGetByOwnerToken(*x.OwnerToken)
	if err != nil {
		bodyMarshal(w, response{"success": false, "message": err.Error()})
		return
	}
	isOwner, err := domainOwnershipVerify(o.OwnerHex, domain)
	if err != nil {
		bodyMarshal(w, response{"success": false, "message": err.Error()})
		return
	}
	if !isOwner {
		bodyMarshal(w, response{"success": false, "message": errorNotAuthorised.Error()})
		return
	}

	// Create label and get primary key
	labelHex, err := labelNew(*x.Name, *x.Color, domain)
	if err != nil {
		bodyMarshal(w, response{"success": false, "message": err.Error()})
		return
	}

	bodyMarshal(w, response{"success": true, "labelHex": labelHex})
}
