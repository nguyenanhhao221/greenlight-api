package main

import (
	"errors"
	"net/http"

	"github.com/nguyenanhhao221/greenlight-api/internal/data"
	"github.com/nguyenanhhao221/greenlight-api/internal/models"
	"github.com/nguyenanhhao221/greenlight-api/internal/validator"
)

func (app *application) createUserHandler(w http.ResponseWriter, r *http.Request) {
	var userInput struct {
		Name     string `json:"name"`
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	err := app.readJSON(w, r, &userInput)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	user := &data.User{Name: userInput.Name, Email: userInput.Email, Activated: false}
	err = user.Password.Set(userInput.Password)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	// Validate user input
	validator := validator.New()
	if data.ValidateUser(validator, user); !validator.Valid() {
		app.failValidationResponse(w, r, validator.Errors)
		return
	}

	// Add user to db
	err = app.models.User.Create(user)
	if err != nil {
		if errors.Is(err, models.ErrDuplicateEmail) {
			validator.AddError("email", "a user with this email address already existed")
			app.failValidationResponse(w, r, validator.Errors)
			return
		}
		app.serverErrorResponse(w, r, err)
		return
	}

	// Send welcome email to user
	app.logger.Info("Create user successfully, sending welcome email!")
	err = app.mailer.Send(*user, "user_welcome.tmpl")
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}

	err = app.writeJSON(w, http.StatusCreated, envelop{"user": user}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}
}
