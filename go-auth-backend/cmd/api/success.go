package main

import (
	"fmt"
	"net/http"
	"strconv"
)

func (app *application) logSuccess(r *http.Request, status int, message interface{}) {
	app.logger.PrintInfo(fmt.Sprintf("%v", message), map[string]string{
		"request_method": r.Method,
		"request_url":    r.URL.String(),
		"status":         strconv.Itoa(status),
	}, app.config.debug)
}

func (app *application) successResponse(w http.ResponseWriter, r *http.Request, status int, message interface{}) {
	env := envelope{"message": message}

	err := app.writeJSON(w, status, env, nil)
	if err != nil {
		app.logError(r, err)
		w.WriteHeader(500)
	}
	app.logSuccess(r, status, message)
}
