package main

import "net/http"

func labelEdit(labelHex string, name string, color string) error {
	// Test for missing field
	if labelHex == "" || name == "" || color == "" {
		return errorMissingField
	}

	statement := `
		UPDATE labels
		SET name = $2, color = $3
		WHERE labelHex = $1;
	`

	_, err := db.Exec(statement, labelHex, name, color)

	if err != nil {
		return errorNoSuchLabel
	}

	return nil
}

func labelEditHandler(w http.ResponseWriter, r *http.Request) {
	type request struct {
		OwnerToken *string `json:"ownerToken"`
		LabelHex   *string `json:"labelHex"`
		Domain     *string `json:"domain"`
		Color      *string `json:"color"`
		Name       *string `json:"name"`
	}

	// Extract request
	var x request
	if err := bodyUnmarshal(r, &x); err != nil {
		bodyMarshal(w, response{"success": false, "message": err.Error()})
		return
	}

	// Get Owner
	o, err := ownerGetByOwnerToken(*x.OwnerToken)
	if err != nil {
		bodyMarshal(w, response{"success": false, "message": err.Error()})
		return
	}

	// Verify Domain ownership
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

	if err := labelEdit(*x.LabelHex, *x.Name, *x.Color); err != nil {
		bodyMarshal(w, response{"success": false, "message": err.Error()})
		return
	}

	bodyMarshal(w, response{"success": true})
}
