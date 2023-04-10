package repository

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/jackc/pgx/v5"
	"gitlab.com/learn-micorservices/profile-service/exception"
	"gitlab.com/learn-micorservices/profile-service/model/domain"
)

var dbName = os.Getenv("DB_NAME")

type ProfileRepository interface {
	GetProfileByID(c context.Context, userID string) (domain.User, error)
	GetProfilesByQuery(c context.Context, params string, value string) (domain.User, error)
	UpdateProfile(c context.Context, user domain.User) error
	UpdatePassword(c context.Context, user domain.User) error

	// kafka
	CreateUser(c context.Context, user domain.User) error
	DeleteUser(c context.Context, user_id string) error
}

type profileRepository struct {
	Database func(dbName string) *pgx.Conn
}

func NewProfileRepository(database func(dbName string) *pgx.Conn) ProfileRepository {
	return &profileRepository{
		Database: database,
	}
}

func (repository *profileRepository) GetProfileByID(c context.Context, userID string) (domain.User, error) {
	ctx, cancel := context.WithTimeout(c, 10*time.Second)
	defer cancel()

	db := repository.Database(dbName)
	defer db.Close(ctx)

	query := `SELECT users.*, roles.name
		FROM users
		LEFT JOIN roles ON roles.id = users.role_id 
		WHERE users.id = $1`

	user, err := db.Query(ctx, query, userID)
	if err != nil {
		fmt.Printf("CollectRows error: %v", err)
		return domain.User{}, exception.ErrInternalServer(err.Error())
	}

	defer user.Close()

	data, err := pgx.CollectOneRow(user, pgx.RowToStructByPos[domain.User])

	if data.ID == "" {
		return domain.User{}, exception.ErrNotFound("user not found")
	}

	if err != nil {
		fmt.Printf("CollectRows error: %v", err)
		return domain.User{}, exception.ErrNotFound("user not found")
	}

	return data, nil
}

func (repository *profileRepository) GetProfilesByQuery(c context.Context, params string, value string) (domain.User, error) {
	ctx, cancel := context.WithTimeout(c, 10*time.Second)
	defer cancel()

	db := repository.Database(dbName)
	defer db.Close(ctx)

	query := fmt.Sprintf(`SELECT SELECT users.*, roles.name
		FROM users
		LEFT JOIN roles ON roles.id = users.role_id 
		WHERE %s = $1`, params)

	user := db.QueryRow(ctx, query, value)

	var data domain.User
	user.Scan(&data.ID, &data.Name, &data.Username, &data.Email, &data.Password, &data.Phone, &data.RoleID, &data.CreatedAt, &data.UpdatedAt, &data.RoleName)

	return data, nil
}

func (repository *profileRepository) UpdateProfile(c context.Context, user domain.User) error {
	ctx, cancel := context.WithTimeout(c, 10*time.Second)
	defer cancel()

	db := repository.Database(dbName)
	defer db.Close(ctx)

	query := `UPDATE users SET 
		name = $1, 
		username = $2, 
		email = $3, 
		password = $4, 
		phone = $5, 
		role_id = $6, 
		updated_at = $7
		WHERE id = $8`

	if _, err := db.Prepare(c, "data", query); err != nil {
		return exception.ErrInternalServer(err.Error())
	}

	if _, err := db.Exec(ctx, "data",
		user.Name,
		user.Username,
		user.Email,
		user.Password,
		user.Phone,
		user.RoleID,
		user.UpdatedAt,
		user.ID); err != nil {
		return exception.ErrUnprocessableEntity(err.Error())
	}

	return nil
}

func (repository *profileRepository) UpdatePassword(c context.Context, user domain.User) error {
	ctx, cancel := context.WithTimeout(c, 10*time.Second)
	defer cancel()

	db := repository.Database(dbName)
	defer db.Close(ctx)

	query := "UPDATE users SET password = $1 WHERE id = $2"

	if _, err := db.Prepare(c, "data", query); err != nil {
		return exception.ErrInternalServer(err.Error())
	}

	if _, err := db.Exec(c, "data", user.Password, user.ID); err != nil {
		return exception.ErrInternalServer(err.Error())
	}

	return nil
}

// function for kafka
func (repository *profileRepository) CreateUser(c context.Context, user domain.User) error {
	ctx, cancel := context.WithTimeout(c, 10*time.Second)
	defer cancel()

	db := repository.Database(dbName)
	defer db.Close(ctx)

	query := `INSERT INTO users (
		id,
		name,
		username,
		email,
		password,
		phone,
		role_id,
		created_at,
		updated_at
	)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)`

	if _, err := db.Prepare(ctx, "data", query); err != nil {
		return exception.ErrInternalServer(err.Error())
	}

	if _, err := db.Exec(ctx, "data",
		user.ID,
		user.Name,
		user.Username,
		user.Email,
		user.Password,
		user.Phone,
		user.RoleID,
		user.CreatedAt,
		user.UpdatedAt); err != nil {
		return exception.ErrUnprocessableEntity(err.Error())
	}

	return nil
}

func (repository *profileRepository) DeleteUser(c context.Context, user_id string) error {
	ctx, cancel := context.WithTimeout(c, 10*time.Second)
	defer cancel()

	db := repository.Database(dbName)
	defer db.Close(ctx)

	query := `DELETE FROM users WHERE id = $1`

	if _, err := db.Prepare(ctx, "data", query); err != nil {
		return exception.ErrInternalServer(err.Error())
	}

	if _, err := db.Exec(ctx, "data", user_id); err != nil {
		return exception.ErrUnprocessableEntity(err.Error())
	}

	return nil
}
