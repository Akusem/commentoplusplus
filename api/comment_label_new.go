package main

import (
	"net/http"
)

func commentNewLabel(commentHex string, labelHex string) error {
	if commentHex == "" || labelHex == "" {
		return errorMissingField
	}

	statement := `
		INSERT INTO comments_labels(commentHex, labelHex)
		VALUES 					   (   $1, 		   $2   );
	`
	_, err := db.Exec(statement, commentHex, labelHex)

	if err != nil {
		return errorCannotLinkLabel
	}

	return nil
}

func commentNewLabelsList(commentHex string, labelsHex []string) error {
	if commentHex == "" || len(labelsHex) == 0 {
		return errorMissingField
	}

	for _, labelHex := range labelsHex {
		err := commentNewLabel(commentHex, labelHex)

		if err != nil {
			return errorCannotLinkLabel
		}
	}

	return nil
}

func commentLabelNewHandler(w http.ResponseWriter, r *http.Request) {
	type request struct {
		CommenterToken *string `json:"commenterToken"`
		CommentHex     *string `json:"commentHex"`
		LabelHex       *string `json:"labelHex"`
	}

	// Access request
	var x request
	if err := bodyUnmarshal(r, &x); err != nil {
		bodyMarshal(w, response{"success": false, "message": err.Error()})
		return
	}

	// Verify the request come from the commenter or from a moderator/domain owner
	c, err := commenterGetByCommenterToken(*x.CommenterToken)
	if err != nil {
		bodyMarshal(w, response{"success": false, "message": err.Error()})
		return
	}

	cm, err := commentGetByCommentHex(*x.CommentHex)
	if err != nil {
		bodyMarshal(w, response{"success": false, "message": err.Error()})
		return
	}

	domain, _, err := commentDomainPathGet(*x.CommentHex)
	if err != nil {
		bodyMarshal(w, response{"success": false, "message": err.Error()})
		return
	}

	isModerator, err := isDomainModerator(domain, c.Email)
	if err != nil {
		bodyMarshal(w, response{"success": false, "message": err.Error()})
		return
	}

	if !isModerator && cm.CommenterHex != c.CommenterHex {
		bodyMarshal(w, response{"success": false, "message": errorNotModerator.Error()})
		return
	}

	if err = commentNewLabel(*x.CommentHex, *x.LabelHex); err != nil {
		bodyMarshal(w, response{"success": false, "message": err.Error()})
		return
	}

	bodyMarshal(w, response{"success": true})
}
