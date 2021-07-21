package main

import (
	"net/http"
)

func labelDelete(labelHex string, domain string) error {
	if labelHex == "" {
		return errorMissingField
	}

	statement := `
		DELETE FROM labels
		WHERE
			labelHex = $1;
	`
	_, err := db.Exec(statement, labelHex)

	if err != nil {
		// TODO: make sure this is the error is actually non-existant labelHex
		return errorNoSuchLabel
	}

	return nil
}

func labelDeleteHandler(w http.ResponseWriter, r *http.Request) {
	type request struct {
		OwnerToken *string `json:"ownerToken"`
		LabelHex   *string `json:"labelHex"`
		Domain     *string `json:"domain"`
	}

	// Extract request
	var x request
	if err := bodyUnmarshal(r, &x); err != nil {
		bodyMarshal(w, response{"success": false, "message": err.Error()})
		return
	}

	domain := domainStrip(*x.Domain)

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

	if err = labelDelete(*x.LabelHex, domain); err != nil {
		bodyMarshal(w, response{"success": false, "message": err.Error()})
		return
	}

	bodyMarshal(w, response{"success": true})
}
