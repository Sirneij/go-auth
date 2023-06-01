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
	router.HandlerFunc(http.MethodPost, "/users/logout/", app.logoutUserHandler)
	router.HandlerFunc(http.MethodGet, "/users/current-user/", app.currentUserHandler)
	router.HandlerFunc(http.MethodPut, "/users/activate/:id/", app.activateUserHandler)
	router.HandlerFunc(http.MethodPost, "/users/regenerate-token/", app.regenerateTokenHandler)
	router.HandlerFunc(http.MethodPost, "/users/password/request-password-change/", app.requestChangePasswordHandler)
	router.HandlerFunc(http.MethodPut, "/users/password/change-user-password/:id/", app.changePasswordHandler)

	return app.recoverPanic(router)
}
