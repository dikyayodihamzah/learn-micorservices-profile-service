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
	GetProfileByID(c context.Context, userID string) (domain.Profile, error)
	GetProfilesByQuery(c context.Context, params string, value string) (domain.Profile, error)
	UpdateProfile(c context.Context, user domain.Profile) error
}

type profileRepository struct {
	Database func(dbName string) *pgx.Conn
}

func NewProfileRepository(database func(dbName string) *pgx.Conn) ProfileRepository {
	return &profileRepository{
		Database: database,
	}
}

func (repository *profileRepository) GetProfileByID(c context.Context, userID string) (domain.Profile, error) {
	ctx, cancel := context.WithTimeout(c, 10*time.Second)
	defer cancel()

	db := repository.Database(dbName)
	defer db.Close(ctx)

	query := `SELECT users.*, roles.name,
		FROM users
		LEFT JOIN roles ON roles.id = users.role_id 
		WHERE users.id = $1`

	user, err := db.Query(ctx, query, userID)
	if err != nil {
		fmt.Printf("CollectRows error: %v", err)
		return domain.Profile{}, exception.ErrInternalServer(err.Error())
	}

	defer user.Close()

	data, err := pgx.CollectOneRow(user, pgx.RowToStructByPos[domain.Profile])

	if data.ID == "" {
		return domain.Profile{}, exception.ErrNotFound("user not found")
	}

	if err != nil {
		fmt.Printf("CollectRows error: %v", err)
		return domain.Profile{}, exception.ErrNotFound("user not found")
	}

	return data, nil
}

func (repository *profileRepository) GetProfilesByQuery(c context.Context, params string, value string) (domain.Profile, error) {
	ctx, cancel := context.WithTimeout(c, 10*time.Second)
	defer cancel()

	db := repository.Database(dbName)
	defer db.Close(ctx)

	query := fmt.Sprintf(`SELECT SELECT users.*,
		FROM users
		WHERE %s = $1`, params)

	user := db.QueryRow(ctx, query, value)

	var data domain.Profile
	user.Scan(&data.ID, &data.Name, &data.Username, &data.Email, &data.Password, &data.Phone, &data.Role, &data.CreatedAt, &data.UpdatedAt)

	return data, nil
}

func (repository *profileRepository) UpdateProfile(c context.Context, user domain.Profile) error {
	ctx, cancel := context.WithTimeout(c, 10*time.Second)
	defer cancel()

	db := repository.Database(dbName)
	defer db.Close(ctx)

	query := `UPDATE profiles SET 
		name = $1, 
		username = $2, 
		email = $3, 
		password = $4, 
		phone = $5, 
		role = $6, 
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
		user.Role,
		user.UpdatedAt); err != nil {
		return exception.ErrUnprocessableEntity(err.Error())
	}

	return nil
}
