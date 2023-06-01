package main

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"time"
)

func (app *application) logoutUserHandler(w http.ResponseWriter, r *http.Request) {
	userID, status, err := app.extractParamsFromSession(r)
	if err != nil {
		switch *status {
		case http.StatusUnauthorized:
			app.unauthorizedResponse(w, r, err)

		case http.StatusBadRequest:
			app.badRequestResponse(w, r, errors.New("invalid cookie"))

		case http.StatusInternalServerError:
			app.serverErrorResponse(w, r, err)

		default:
			app.serverErrorResponse(w, r, errors.New("something happened and we could not fullfil your request at the moment"))
		}
		return
	}

	// Get session from redis
	_, err = app.getFromRedis(fmt.Sprintf("sessionid_%s", userID.Id))
	if err != nil {
		app.unauthorizedResponse(w, r, errors.New("you are not authorized to access this resource"))
		return
	}

	// Delete session from redis
	ctx := context.Background()
	_, err = app.redisClient.Del(ctx, fmt.Sprintf("sessionid_%s", userID.Id)).Result()
	if err != nil {
		app.logger.PrintError(err, map[string]string{
			"key": fmt.Sprintf("sessionid_%s", userID.Id),
		})
		app.serverErrorResponse(w, r, errors.New("something happened decosing your cookie data"))
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:    "sessionid",
		Value:   "",
		Expires: time.Now(),
	})

	// Respond with success
	app.successResponse(w, r, http.StatusOK, "You have successfully logged out")
}
