package main

import (
	"context"
	"encoding/gob"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/julienschmidt/httprouter"
	"goauthbackend.johnowolabiidogun.dev/internal/cookies"
	"goauthbackend.johnowolabiidogun.dev/internal/data"
)

func (app *application) writeJSON(w http.ResponseWriter, status int, data interface{}, headers http.Header) error {
	js, err := json.Marshal(data)

	if err != nil {
		return err
	}

	js = append(js, '\n')

	for key, value := range headers {
		w.Header()[key] = value
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	w.Write(js)
	return nil
}

func (app *application) readJSON(w http.ResponseWriter, r *http.Request, dst interface{}) error {

	maxBytes := 1_048_576
	r.Body = http.MaxBytesReader(w, r.Body, int64(maxBytes))
	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()
	err := dec.Decode(dst)
	if err != nil {

		var syntaxError *json.SyntaxError
		var unmarshalTypeError *json.UnmarshalTypeError
		var invalidUnmarshalError *json.InvalidUnmarshalError
		switch {

		case errors.As(err, &syntaxError):
			return fmt.Errorf("body contains badly-formed JSON (at character %d)", syntaxError.Offset)

		case errors.Is(err, io.ErrUnexpectedEOF):
			return errors.New("body contains badly-formed JSON")

		case errors.As(err, &unmarshalTypeError):
			if unmarshalTypeError.Field != "" {
				return fmt.Errorf("body contains incorrect JSON type for field %q", unmarshalTypeError.Field)
			}
			return fmt.Errorf("body contains incorrect JSON type (at character %d)", unmarshalTypeError.Offset)

		case errors.Is(err, io.EOF):
			return errors.New("body must not be empty")

		case strings.HasPrefix(err.Error(), "json: unknown field "):
			fieldName := strings.TrimPrefix(err.Error(), "json: unknown field ")
			return fmt.Errorf("body contains unknown key %s", fieldName)

		case err.Error() == "http: request body too large":
			return fmt.Errorf("body must not be larger than %d bytes", maxBytes)

		case errors.As(err, &invalidUnmarshalError):
			panic(err)

		default:
			return err
		}
	}
	err = dec.Decode(&struct{}{})
	if err != io.EOF {
		return errors.New("body must only contain a single JSON value")
	}
	return nil

}

func (app *application) background(fn func()) {
	app.wg.Add(1)

	go func() {

		defer app.wg.Done()
		// Recover any panic.
		defer func() {
			if err := recover(); err != nil {
				app.logger.PrintError(fmt.Errorf("%s", err), nil)
			}
		}()
		// Execute the arbitrary function that we passed as the parameter.
		fn()
	}()
}

func (app *application) storeInRedis(prefix string, hash string, userID uuid.UUID, expiration time.Duration) error {
	ctx := context.Background()
	err := app.redisClient.Set(
		ctx,
		fmt.Sprintf("%s%s", prefix, userID),
		hash,
		expiration,
	).Err()
	if err != nil {
		return err
	}

	return nil
}

func (app *application) getFromRedis(key string) (*string, error) {
	ctx := context.Background()

	hash, err := app.redisClient.Get(ctx, key).Result()
	if err != nil {
		return nil, err
	}

	return &hash, nil
}

func (app *application) readIDParam(r *http.Request) (*uuid.UUID, error) {
	params := httprouter.ParamsFromContext(r.Context())
	id, err := uuid.Parse(params.ByName("id"))
	if err != nil {
		return nil, errors.New("invalid id parameter")
	}
	return &id, nil
}

func (app *application) extractParamsFromSession(r *http.Request) (*data.UserID, *int, error) {
	gobEncodedValue, err := cookies.ReadEncrypted(r, "sessionid", app.config.secret.secretKey)

	if err != nil {
		var errorData error
		var status int
		switch {
		case errors.Is(err, http.ErrNoCookie):
			app.logger.PrintError(err, nil)
			status = http.StatusUnauthorized
			errorData = errors.New("you are not authorized to access this resource")

		case errors.Is(err, cookies.ErrInvalidValue):
			app.logger.PrintError(err, nil)
			status = http.StatusBadRequest
			errorData = errors.New("invalid cookie")

		default:
			app.logger.PrintError(err, nil)
			status = http.StatusInternalServerError
			errorData = errors.New("something happened getting your cookie data")

		}
		return nil, &status, errorData
	}

	var userID data.UserID

	reader := strings.NewReader(gobEncodedValue)
	if err := gob.NewDecoder(reader).Decode(&userID); err != nil {
		app.logger.PrintError(err, nil)
		status := http.StatusInternalServerError
		return nil, &status, errors.New("something happened decosing your cookie data")
	}

	return &userID, nil, nil
}
