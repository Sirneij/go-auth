package main

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"goauthbackend.johnowolabiidogun.dev/internal/data"
	"goauthbackend.johnowolabiidogun.dev/internal/tokens"
	"goauthbackend.johnowolabiidogun.dev/internal/validator"
)

func (app *application) requestChangePasswordHandler(w http.ResponseWriter, r *http.Request) {

	expirationInt, err := strconv.Atoi(strings.Split(app.config.tokenExpiration.durationString, "m")[0])
	if err != nil {
		app.logger.PrintError(err, nil)
		app.serverErrorResponse(w, r,
			errors.New("something happened and we could not fullfil your request at the moment"))

		return
	}
	expirationStr := fmt.Sprintf("%dm", expirationInt)

	// Expected data from the user
	var input struct {
		Email string `json:"email"`
	}
	// Try reading the user input to JSON
	err = app.readJSON(w, r, &input)
	if err != nil {
		app.logger.PrintError(err, nil)
		app.badRequestResponse(w, r, err)
		return
	}

	// Validate the user input
	v := validator.New()
	if data.ValidateEmail(v, input.Email); !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	db_user, err := app.models.Users.GetEmail(input.Email, true)
	if err != nil {
		app.logger.PrintError(err, nil)
		app.badRequestResponse(w, r, err)
		return
	}

	// Generate 6-digit token
	otp, err := tokens.GenerateOTP()
	if err != nil {
		app.logger.PrintError(err, nil)
		app.serverErrorResponse(w, r, errors.New("something happened and we could not fullfil your request at the moment"))

		return
	}

	err = app.storeInRedis("password_reset_", otp.Hash, db_user.ID, (app.config.tokenExpiration.duration * 2))
	if err != nil {
		app.logger.PrintError(err, nil)
		app.serverErrorResponse(w, r,
			errors.New("something happened and we could not fullfil your request at the moment"),
		)
		return
	}

	now := time.Now()
	expiration := now.Add(app.config.tokenExpiration.duration * 2)
	exact := expiration.Format(time.RFC1123)

	// Send email to user, using separate goroutine, for account activation
	app.background(func() {
		data := map[string]interface{}{
			"name":        fmt.Sprintf("%s %s", db_user.FirstName, db_user.LastName),
			"token":       tokens.FormatOTP(otp.Secret),
			"userID":      db_user.ID,
			"frontendURL": app.config.frontendURL,
			"expiration":  expirationStr,
			"exact":       exact,
		}
		err = app.mailer.Send(db_user.Email, "password_reset.tmpl", data)
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
		"You requested a password change. Check your email address and follow the instruction to change your password. Ensure your password is changed before the token expires",
	)
}
