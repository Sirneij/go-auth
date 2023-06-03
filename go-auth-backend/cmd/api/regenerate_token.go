package main

import (
	"errors"
	"net/http"
	"time"

	"goauthbackend.johnowolabiidogun.dev/internal/data"
	"goauthbackend.johnowolabiidogun.dev/internal/tokens"
	"goauthbackend.johnowolabiidogun.dev/internal/validator"
)

func (app *application) regenerateTokenHandler(w http.ResponseWriter, r *http.Request) {
	// Expected data from the user
	var input struct {
		Email string `json:"email"`
	}
	// Try reading the user input to JSON
	err := app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	// Validate the user input
	v := validator.New()
	if data.ValidateEmail(v, input.Email); !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	db_user, err := app.models.Users.GetEmail(input.Email, false)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	// Generate 6-digit token
	otp, err := tokens.GenerateOTP()
	if err != nil {
		app.serverErrorResponse(w, r, errors.New("something happened and we could not fullfil your request at the moment"))

		return
	}

	err = app.storeInRedis("activation_", otp.Hash, db_user.ID, app.config.tokenExpiration.duration)
	if err != nil {
		app.logError(r, err)
	}

	now := time.Now()
	expiration := now.Add(app.config.tokenExpiration.duration)
	exact := expiration.Format(time.RFC1123)

	// Send email to user, using separate goroutine, for account activation
	app.background(func() {
		data := map[string]interface{}{
			"token":       tokens.FormatOTP(otp.Secret),
			"userID":      db_user.ID,
			"frontendURL": app.config.frontendURL,
			"expiration":  app.config.tokenExpiration.durationString,
			"exact":       exact,
		}
		err = app.mailer.Send(db_user.Email, "user_welcome.tmpl", data)
		if err != nil {
			app.logError(r, err)
		}
		app.logger.PrintInfo("Email successfully sent.", nil, app.config.debug)
	})

	// Respond with success
	app.successResponse(
		w,
		r,
		http.StatusAccepted,
		"Account activation link has been sent to your email address. Kindly take action before its expiration",
	)
}
