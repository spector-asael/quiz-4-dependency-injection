package handler

import (
	"net/http"
	"github.com/spector-asael/lab4-crud/internal/data"
	"github.com/spector-asael/lab4-crud/internal/validator"
)

func (a *ApplicationDependencies) transferHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		FromAccountID int64   `json:"from_account_id"`
		ToAccountID   int64   `json:"to_account_id"`
		Amount        float64 `json:"amount"`
	}

	// Decode JSON
	if err := a.readJSON(w, r, &input); err != nil {
		a.badRequestResponse(w, r, err)
		return
	}

	// Validate
	v := validator.New()
	data.ValidateTransfer(v, data.TransferInput(input))
	if !v.IsEmpty() {
		a.failedValidationResponse(w, r, v.Errors)
		return
	}

	// Perform transfer using model
	if err := a.Models.Transfers.Create(data.TransferInput(input)); err != nil {
		a.serverErrorResponse(w, r, err)
		return
	}

	// Respond
	a.writeJSON(w, http.StatusOK, envelope{
		"message":         "transfer completed successfully",
		"from_account_id": input.FromAccountID,
		"to_account_id":   input.ToAccountID,
		"amount":          input.Amount,
	}, nil)
}