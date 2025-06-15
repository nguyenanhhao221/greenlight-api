package main

import (
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/nguyenanhhao221/greenlight-api/internal/data"
	"github.com/nguyenanhhao221/greenlight-api/internal/models"
	"github.com/nguyenanhhao221/greenlight-api/internal/validator"
)

func (app *application) createAuthenticationToken(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	if err := app.readJSON(w, r, &input); err != nil {
		err = fmt.Errorf("error when readJSON in createAuthenticationToken %w", err)
		app.serverErrorResponse(w, r, err)
	}

	// Validate user input
	v := validator.New()
	data.ValidateEmail(v, input.Email)
	data.ValidatePasswordPlainText(v, input.Password)
	if !v.Valid() {
		app.failValidationResponse(w, r, v.Errors)
		return
	}

	user, err := app.models.User.GetByEmail(input.Email)
	if err != nil {
		switch {
		case errors.Is(err, models.ErrRecordNotFound):
			app.notFoundResponse(w, r)
			return
		default:
			err = fmt.Errorf("error GetByEmail in createAuthenticationToken: %w", err)
			app.serverErrorResponse(w, r, err)
		}
	}
	isPasswordMatches, err := user.Password.Matches(input.Password)
	if err != nil {
		err = fmt.Errorf("error calling Password.Matches in createAuthenticationToken: %w", err)
		app.serverErrorResponse(w, r, err)
		return
	}

	if !isPasswordMatches {
		app.invalidCredentialsResponse(w, r)
		return
	}

	// Create new authentication token with 24hour expiry
	authenticationToken, err := app.models.Token.New(user.ID, 24*time.Hour, data.ScopeAuthentication)
	if err != nil {
		err = fmt.Errorf("error calling Token.New in createAuthenticationToken %w", err)
		app.serverErrorResponse(w, r, err)
	}

	if err := app.writeJSON(w, http.StatusCreated, envelop{"authentication_token": authenticationToken}, nil); err != nil {
		app.serverErrorResponse(w, r, err)
	}
}
