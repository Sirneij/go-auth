package main

import (
	"bytes"
	"encoding/gob"
	"errors"
	"net/http"

	"goauthbackend.johnowolabiidogun.dev/internal/cookies"
	"goauthbackend.johnowolabiidogun.dev/internal/data"
)

func (app *application) loginUserHandler(w http.ResponseWriter, r *http.Request) {
	// Expected data from the user
	var input struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	// Try reading the user input to JSON
	err := app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	db_user, err := app.models.Users.GetEmail(input.Email, true)
	if err != nil {
		app.logger.PrintError(err, nil)
		app.badRequestResponse(w, r, err)
		return
	}

	match, err := db_user.Password.Matches(input.Password)
	if err != nil {
		app.logger.PrintError(err, nil)
		return
	}

	if !match {
		app.logger.PrintError(errors.New("email and password combination does not match"), nil)
		app.badRequestResponse(w, r, errors.New("email and password combination does not match"))
		return
	}

	var userID = data.UserID{
		Id: db_user.ID,
	}

	var buf bytes.Buffer

	// Gob-encode the user data, storing the encoded output in the buffer.
	err = gob.NewEncoder(&buf).Encode(&userID)
	if err != nil {
		app.logger.PrintFatal(err, nil)
		app.serverErrorResponse(w, r, errors.New("something happened encoding your data"))
		return
	}

	session := buf.String()

	// Store session in redis
	err = app.storeInRedis("sessionid_", session, userID.Id, app.config.secret.sessionExpiration)
	if err != nil {
		app.logger.PrintError(err, nil)
	}

	cookie := http.Cookie{
		Name:     "sessionid",
		Value:    session,
		Path:     "/",
		MaxAge:   int(app.config.secret.sessionExpiration.Seconds()),
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteLaxMode,
	}

	// Write an encrypted cookie containing the gob-encoded data as normal.
	err = cookies.WriteEncrypted(w, cookie, app.config.secret.secretKey)
	if err != nil {
		app.logger.PrintFatal(err, nil)
		app.serverErrorResponse(w, r, errors.New("something happened setting your cookie data"))
		return
	}

	app.writeJSON(w, http.StatusOK, db_user, nil)
}
