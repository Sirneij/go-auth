package data

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/google/uuid"
)

func (um UserModel) Insert(user *User) (*UserID, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	tx, err := um.DB.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}

	var userID uuid.UUID

	query_user := `
	INSERT INTO users (email, password, first_name, last_name) VALUES ($1, $2, $3, $4) RETURNING id`
	args_user := []interface{}{user.Email, user.Password.hash, user.FirstName, user.LastName}

	if err := tx.QueryRowContext(ctx, query_user, args_user...).Scan(&userID); err != nil {
		switch {
		case err.Error() == `pq: duplicate key value violates unique constraint "users_email_key"`:
			return nil, ErrDuplicateEmail
		default:
			return nil, err
		}

	}

	query_user_profile := `
	INSERT INTO user_profile (user_id) VALUES ($1) ON CONFLICT (user_id) DO NOTHING RETURNING user_id`

	_, err = tx.ExecContext(ctx, query_user_profile, userID)

	if err != nil {
		return nil, err
	}

	if err = tx.Commit(); err != nil {
		return nil, err
	}
	id := UserID{
		Id: userID,
	}

	return &id, nil
}

func (um UserModel) Get(id uuid.UUID) (*User, error) {
	query := `
	SELECT 
		u.*, p.* 
	FROM 
		users u 
		LEFT JOIN user_profile p ON p.user_id = u.id 
	WHERE 
		u.is_active = true AND u.id = $1
	`
	var user User
	var userP UserProfile
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	err := um.DB.QueryRowContext(ctx, query, id).Scan(&user.ID,
		&user.Email, &user.Password.hash, &user.FirstName, &user.LastName, &user.IsActive, &user.IsStaff, &user.IsSuperuser, &user.Thumbnail, &user.DateJoined, &userP.ID, &userP.UserID, &userP.PhoneNumber, &userP.BirthDate, &userP.GithubLink,
	)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrRecordNotFound
		default:
			return nil, err
		}
	}
	user.Profile = userP
	return &user, nil
}

func (um UserModel) ActivateUser(userID uuid.UUID) (*sql.Result, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := `UPDATE users SET is_active = true WHERE id = $1`

	result, err := um.DB.ExecContext(ctx, query, userID)
	if err != nil {
		return nil, err
	}

	return &result, nil
}

func (um UserModel) UpdateUserPassword(user *User) (*sql.Result, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := `UPDATE users SET password = $1 WHERE id = $2`

	result, err := um.DB.ExecContext(ctx, query, user.Password.hash, user.ID)
	if err != nil {
		return nil, err
	}

	return &result, nil
}

func (um UserModel) Update(user *User, userP *UserProfile) (*User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var userOut User
	var userPOut UserProfile

	tx, err := um.DB.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}

	query_user := `
	UPDATE 
		users 
	SET 
		first_name = COALESCE($1, first_name), 
		last_name = COALESCE($2, last_name), 
		thumbnail = COALESCE($3, thumbnail)
	WHERE 
		id = $4 
		AND is_active = true 
		AND (
			$1 IS NOT NULL 
			AND $1 IS DISTINCT 
			FROM 
				first_name 
				OR $2 IS NOT NULL 
				AND $2 IS DISTINCT 
			FROM 
				last_name 
				OR $3 IS DISTINCT 
			FROM 
				thumbnail
		)
	RETURNING id, email, password, first_name, last_name, is_active, is_staff, is_superuser, thumbnail, date_joined,
	`
	args_user := []interface{}{user.FirstName, user.LastName, user.Thumbnail, user.ID}

	err = tx.QueryRowContext(ctx, query_user, args_user...).Scan(&userOut.ID,
		&userOut.Email, &userOut.Password.hash, &userOut.FirstName, &userOut.LastName, &userOut.IsActive, &userOut.IsStaff, &userOut.IsSuperuser, &userOut.Thumbnail, &userOut.DateJoined)

	if err != nil {
		return nil, err
	}

	query_user_profile := `
	UPDATE 
		user_profile 
	SET 
		phone_number = NULLIF($1, ''), 
		birth_date = $2, 
		github_link = NULLIF($3, '')
	WHERE 
		user_id = $4 
		AND (
			$1 IS DISTINCT 
			FROM 
				phone_number 
				OR $2 IS DISTINCT 
			FROM 
				birth_date 
				OR $3 IS DISTINCT 
			FROM 
				github_link
		)
	RETURNING id, user_id, phone_number, birth_date, github_link
	`

	args_profile_user := []interface{}{
		userP.PhoneNumber,
		userP.BirthDate,
		userP.GithubLink,
		user.ID,
	}

	err = tx.QueryRowContext(ctx, query_user_profile, args_profile_user...).Scan(&userPOut.ID, &userPOut.UserID, &userPOut.PhoneNumber, &userPOut.BirthDate, &userPOut.GithubLink)

	if err != nil {
		return nil, err
	}

	if err = tx.Commit(); err != nil {
		return nil, err
	}

	userOut.Profile = userPOut

	return &userOut, nil
}

func (um UserModel) GetEmail(email string, active bool) (*User, error) {
	query := `
	SELECT 
		u.*, p.*
	FROM 
		users u 
		JOIN user_profile p ON p.user_id = u.id 
	WHERE 
		u.is_active = $2 AND u.email = $1`

	var user User
	var userP UserProfile

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := um.DB.QueryRowContext(ctx, query, email, active).Scan(
		&user.ID,
		&user.Email,
		&user.Password.hash,
		&user.FirstName,
		&user.LastName,
		&user.IsActive,
		&user.IsStaff,
		&user.IsSuperuser,
		&user.Thumbnail,
		&user.DateJoined,
		&userP.ID,
		&userP.UserID,
		&userP.PhoneNumber,
		&userP.BirthDate,
		&userP.GithubLink,
	)

	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			if active {
				return nil, ErrRecordNotFound
			} else {
				return nil, errors.New("an inactive user with the provided email address was not found")
			}
		default:
			return nil, err
		}
	}

	user.Profile = userP

	return &user, nil
}
