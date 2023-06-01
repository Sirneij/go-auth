package main

import (
	"errors"
	"fmt"
	"net/http"
)

func (app *application) currentUserHandler(w http.ResponseWriter, r *http.Request) {

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
			app.serverErrorResponse(
				w,
				r,
				errors.New("something happened and we could not fullfil your request at the moment"),
			)
		}
		return
	}

	// Get session from redis
	_, err = app.getFromRedis(fmt.Sprintf("sessionid_%s", userID.Id))
	if err != nil {
		app.unauthorizedResponse(w, r, errors.New("you are not authorized to access this resource"))
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
