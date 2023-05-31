package main

import (
	"bytes"
	"context"
	"crypto/sha256"
	"encoding/gob"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"goauthbackend.johnowolabiidogun.dev/internal/cookies"
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
			app.failedValidationResponse(w, r, v.Errors)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	// Generate 6-digit token
	otp, err := tokens.GenerateOTP()
	if err != nil {
		app.logger.PrintError(err, nil)
	}

	err = app.storeHashInRedis(otp.Hash, userID.Id)
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
		"Your account creation was accepted successfully. Check your email address and follow the instruction to actuivate your account. Ensure you activate your account before the token expires",
	)
}

func (app *application) activateUserHandler(w http.ResponseWriter, r *http.Request) {
	id, err := app.readIDParam(r)

	if err != nil {
		app.notFoundResponse(w, r)
		return
	}

	var input struct {
		Secret string `json:"token"`
	}

	err = app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	v := validator.New()
	if tokens.ValidateSecret(v, input.Secret); !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	hash, err := app.getHashFromRedis(fmt.Sprintf("activation_%s", id))
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	tokenHash := fmt.Sprintf("%x\n", sha256.Sum256([]byte(input.Secret)))

	if *hash != tokenHash {
		app.failedValidationResponse(w, r, map[string]string{
			"token": "The supplied token is invalid",
		})
		return
	}

	ctx := context.Background()
	deleted, err := app.redisClient.Del(ctx, fmt.Sprintf("activation_%s", id)).Result()

	if err != nil {
		app.logger.PrintError(err, map[string]string{
			"key": fmt.Sprintf("activation_%s", id),
		})

	}

	app.logger.PrintInfo(fmt.Sprintf("Token hash was deleted successfully :activation_%d", deleted), nil)

	result, err := app.models.Users.ActivateUser(*id)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	app.logger.PrintInfo(fmt.Sprintf("%x", result), nil)
}

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

	db_user, err := app.models.Users.GetEmail(input.Email)

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

	app.logger.PrintInfo("Creating a cookie", nil)
	cookie := http.Cookie{
		Name:     "sessionid",
		Value:    buf.String(),
		Path:     "/",
		MaxAge:   3600,
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

func (app *application) currentUserHandler(w http.ResponseWriter, r *http.Request) {

	gobEncodedValue, err := cookies.ReadEncrypted(r, "sessionid", app.config.secret.secretKey)

	if err != nil {
		switch {
		case errors.Is(err, http.ErrNoCookie):
			app.logger.PrintError(err, nil)
			app.badRequestResponse(w, r, errors.New("cookie not found"))
		case errors.Is(err, cookies.ErrInvalidValue):
			app.logger.PrintError(err, nil)
			app.badRequestResponse(w, r, errors.New("invalid cookie"))
		default:
			app.logger.PrintError(err, nil)
			app.serverErrorResponse(w, r, errors.New("something happened getting your cookie data"))
		}
		return
	}

	var userID data.UserID

	reader := strings.NewReader(gobEncodedValue)
	if err := gob.NewDecoder(reader).Decode(&userID); err != nil {
		app.logger.PrintError(err, nil)
		app.serverErrorResponse(w, r, errors.New("something happened decosing your cookie data"))
		return
	}

	db_user, err := app.models.Users.Get(userID.Id)

	if err != nil {
		app.logger.PrintFatal(err, nil)
		app.badRequestResponse(w, r, err)
		return
	}

	app.writeJSON(w, http.StatusOK, db_user, nil)
}
