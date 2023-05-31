package main

import "net/http"

func (app *application) successResponse(w http.ResponseWriter, r *http.Request, status int, message interface{}) {
	env := envelope{"message": message}

	err := app.writeJSON(w, status, env, nil)
	if err != nil {
		app.logError(r, err)
		w.WriteHeader(500)
	}
}
