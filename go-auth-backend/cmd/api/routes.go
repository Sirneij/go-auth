package main

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func (app *application) routes() http.Handler {
	router := httprouter.New()

	router.NotFound = http.HandlerFunc(app.notFoundResponse)
	router.MethodNotAllowed = http.HandlerFunc(app.methodNotAllowedResponse)

	router.HandlerFunc(http.MethodGet, "/healthcheck/", app.healthcheckHandler)
	router.HandlerFunc(http.MethodPost, "/users/register/", app.registerUserHandler)
	router.HandlerFunc(http.MethodPost, "/users/login/", app.loginUserHandler)
	router.HandlerFunc(http.MethodGet, "/users/current-user/", app.currentUserHandler)
	router.HandlerFunc(http.MethodPut, "/users/activate/:id/", app.activateUserHandler)

	return app.recoverPanic(router)
}
