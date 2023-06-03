package main

import (
	"errors"
	"net/http"

	"goauthbackend.johnowolabiidogun.dev/internal/data"
	"goauthbackend.johnowolabiidogun.dev/internal/types"
	"goauthbackend.johnowolabiidogun.dev/internal/validator"
)

func (app *application) updateUserHandler(w http.ResponseWriter, r *http.Request) {
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

	db_user, err := app.models.Users.Get(userID.Id)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	var input struct {
		FirstName   *string        `json:"first_name"`
		LastName    *string        `json:"last_name"`
		Thumbnail   *string        `json:"thumbnail"`
		PhoneNumber *string        `json:"phone_number"`
		BirthDate   types.NullTime `json:"birth_date"`
		GithubLink  *string        `json:"github_link"`
	}

	err = app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	if input.FirstName != nil {
		db_user.FirstName = *input.FirstName
	}
	if input.LastName != nil {
		db_user.LastName = *input.LastName
	}
	if input.Thumbnail != nil {
		db_user.Thumbnail = input.Thumbnail
	}
	if input.PhoneNumber != nil {
		db_user.Profile.PhoneNumber = input.PhoneNumber
	}
	if input.BirthDate.Valid {
		db_user.Profile.BirthDate = input.BirthDate
	}
	if input.GithubLink != nil {
		db_user.Profile.GithubLink = input.GithubLink
	}

	v := validator.New()
	if data.ValidateUser(v, db_user); !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	updated_user, err := app.models.Users.Update(db_user)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	err = app.writeJSON(w, http.StatusOK, updated_user, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
	app.logSuccess(r, http.StatusOK, "User updated successfully")
}
