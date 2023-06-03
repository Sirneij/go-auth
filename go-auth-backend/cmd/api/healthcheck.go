package main

import (
	"net/http"
	"strconv"
)

func (app *application) healthcheckHandler(w http.ResponseWriter, r *http.Request) {
	data := map[string]string{
		"status":  "available",
		"debug":   strconv.FormatBool(app.config.debug),
		"version": version,
	}

	err := app.writeJSON(w, http.StatusOK, data, nil)

	if err != nil {
		app.serverErrorResponse(w, r, err)
	}

}
