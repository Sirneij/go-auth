package main

import (
	"errors"
	"net/http"
)

func (app *application) uploadFileToS3Handler(w http.ResponseWriter, r *http.Request) {
	_, status, err := app.extractParamsFromSession(r)
	if err != nil {
		switch *status {
		case http.StatusUnauthorized:
			app.logger.PrintError(err, nil)
			app.unauthorizedResponse(w, r, err)

		case http.StatusBadRequest:
			app.logger.PrintError(err, nil)
			app.badRequestResponse(w, r, errors.New("invalid cookie"))

		case http.StatusInternalServerError:
			app.logger.PrintError(err, nil)
			app.serverErrorResponse(w, r, err)

		default:
			app.logger.PrintError(err, nil)
			app.serverErrorResponse(w, r, errors.New("something happened and we could not fullfil your request at the moment"))
		}
		return
	}

	s3URL, err := app.uploadFileToS3(r)
	if err != nil {
		app.logger.PrintError(err, nil)
		app.badRequestResponse(w, r, err)
		return
	}

	env := envelope{"s3_url": s3URL}

	err = app.writeJSON(w, http.StatusOK, env, nil)
	if err != nil {
		app.logger.PrintError(err, nil)
		app.serverErrorResponse(w, r, err)
	}
}
