package main

import (
	"errors"
	"net/http"
)

func (app *application) deleteFileOnS3Handler(w http.ResponseWriter, r *http.Request) {
	_, status, err := app.extractParamsFromSession(r)
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

	_, err = app.deleteFileFromS3(r)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	app.successResponse(w, r, http.StatusNoContent, "Image deleted successfully.")
}
