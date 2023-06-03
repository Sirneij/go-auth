package main

import (
	"errors"
	"expvar"
	"fmt"
	"net/http"
	"strconv"

	"github.com/felixge/httpsnoop"
)

func (app *application) recoverPanic(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				w.Header().Set("Connection", "close")
				app.serverErrorResponse(w, r, fmt.Errorf("%s", err))
			}
		}()
		next.ServeHTTP(w, r)
	})
}

func (app *application) metrics(next http.Handler) http.Handler {
	totalRequestsReceived := expvar.NewInt("total_requests_received")
	totalResponsesSent := expvar.NewInt("total_responses_sent")
	totalProcessingTimeMicroseconds := expvar.NewInt("total_processing_time_μs")

	totalResponsesSentByStatus := expvar.NewMap("total_responses_sent_by_status")
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		totalRequestsReceived.Add(1)

		metrics := httpsnoop.CaptureMetrics(next, w, r)

		totalResponsesSent.Add(1)

		totalProcessingTimeMicroseconds.Add(metrics.Duration.Microseconds())

		totalResponsesSentByStatus.Add(strconv.Itoa(metrics.Code), 1)
	})
}

func (app *application) authenticateAndAuthorize(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
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
			app.badRequestResponse(w, r, err)
			return
		}
		if !db_user.IsSuperuser {
			app.unauthorizedResponse(w, r, errors.New("you are not authorized to access this resource"))
			return
		}
		next.ServeHTTP(w, r)

	})
}
