package main

import (
	"context"
	"crypto/sha256"
	"errors"
	"fmt"
	"net/http"

	"goauthbackend.johnowolabiidogun.dev/internal/tokens"
	"goauthbackend.johnowolabiidogun.dev/internal/validator"
)

func (app *application) activateUserHandler(w http.ResponseWriter, r *http.Request) {
	id, err := app.readIDParam(r)

	if err != nil {
		app.badRequestResponse(w, r, err)
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

	hash, err := app.getFromRedis(fmt.Sprintf("activation_%s", id))
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	tokenHash := fmt.Sprintf("%x\n", sha256.Sum256([]byte(input.Secret)))

	if *hash != tokenHash {
		app.logger.PrintError(errors.New("the supplied token is invalid"), nil, app.config.debug)
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
		}, app.config.debug)

	}

	app.logger.PrintInfo(fmt.Sprintf("Token hash was deleted successfully :activation_%d", deleted), nil, app.config.debug)

	_, err = app.models.Users.ActivateUser(*id)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	app.successResponse(w, r, http.StatusOK, "Account activated successfully.")
}
