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

	// Create a activation token to be sent to user welcome email
	token, err := app.models.Token.New(user.ID, (24 * 3 * time.Hour), data.ScopeActivation)
	// TODO: if there is an error when creating token, should rollback the user creation as well, or has another endpoint to re-trigger the token creation
	if err != nil {
		err = fmt.Errorf("error create token %w", err)
		app.serverErrorResponse(w, r, err)
		return
	}
	// Send welcome email to user in the background with activation token
	app.logger.Info("Create user and activation token successfully, sending welcome email!")
	app.background(func() {
		data := map[string]any{
			"activationToken": token.Plain,
			"userId":          user.ID,
		}
		err = app.mailer.Send(user.Email, "user_welcome.tmpl", data)
		if err != nil {
			app.logger.Error("Error sending user to", "user: email", user.Email)
		}
		app.logger.Info("Email sent successfully for:", "user email", user.Email)
	})

	err = app.writeJSON(w, http.StatusCreated, envelop{"user": user}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}
}

