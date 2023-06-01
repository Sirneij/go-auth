package main

import (
	"errors"
	"net/http"
	"time"

	"goauthbackend.johnowolabiidogun.dev/internal/data"
	"goauthbackend.johnowolabiidogun.dev/internal/tokens"
	"goauthbackend.johnowolabiidogun.dev/internal/validator"
)

func (app *application) registerUserHandler(w http.ResponseWriter, r *http.Request) {
	// Expected data from the user
	var input struct {
		Email     string `json:"email"`
		FirstName string `json:"first_name"`
		LastName  string `json:"last_name"`
		Password  string `json:"password"`
	}
	// Try reading the user input to JSON
	err := app.readJSON(w, r, &input)
	if err != nil {
		app.logger.PrintError(err, nil)
		app.badRequestResponse(w, r, err)
		return
	}

	user := &data.User{
		Email:     input.Email,
		FirstName: input.FirstName,
		LastName:  input.LastName,
	}

	// Hash user password
	err = user.Password.Set(input.Password)
	if err != nil {
		app.logger.PrintError(err, nil)
		app.serverErrorResponse(w, r, err)
		return
	}

	// Validate the user input
	v := validator.New()
	if data.ValidateUser(v, user); !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	// Save the user in the database
	userID, err := app.models.Users.Insert(user)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrDuplicateEmail):
			v.AddError("email", "A user with this email address already exists")
			app.logger.PrintError(err, nil)
			app.failedValidationResponse(w, r, v.Errors)
		default:
			app.logger.PrintError(err, nil)
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	// Generate 6-digit token
	otp, err := tokens.GenerateOTP()
	if err != nil {
		app.logger.PrintError(err, nil)
	}

	err = app.storeInRedis("activation_", otp.Hash, userID.Id, app.config.tokenExpiration.duration)
	if err != nil {
		app.logger.PrintError(err, nil)
	}

	now := time.Now()
	expiration := now.Add(app.config.tokenExpiration.duration)
	exact := expiration.Format(time.RFC1123)

	// Send email to user, using separate goroutine, for account activation
	app.background(func() {
		data := map[string]interface{}{
			"token":       tokens.FormatOTP(otp.Secret),
			"userID":      userID.Id,
			"frontendURL": app.config.frontendURL,
			"expiration":  app.config.tokenExpiration.durationString,
			"exact":       exact,
		}
		err = app.mailer.Send(user.Email, "user_welcome.tmpl", data)
		if err != nil {
			app.logger.PrintError(err, nil)
		}
		app.logger.PrintInfo("Email successfully sent.", nil)
	})

	// Respond with success
	app.successResponse(
		w,
		r,
		http.StatusAccepted,
		"Your account creation was accepted successfully. Check your email address and follow the instruction to activate your account. Ensure you activate your account before the token expires",
	)
}
