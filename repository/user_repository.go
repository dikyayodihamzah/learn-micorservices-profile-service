package repository

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/jackc/pgx/v5"
	"gitlab.com/learn-micorservices/user-service/exception"
	"gitlab.com/learn-micorservices/user-service/model/domain"
)

var dbName = os.Getenv("DB_NAME")

type UserRepository interface {
	CreateUser(c context.Context, user domain.User) error
	GetAllUser(c context.Context) ([]domain.User, error)
	GetUserByID(c context.Context, userID string) (domain.User, error)
	GetUsersByQuery(c context.Context, params string, value string) (domain.User, error)
	UpdateUser(c context.Context, user domain.User) error
	UpdateUserPassword(c context.Context, user domain.User) error
	DeleteUser(c context.Context, user_id string) error
}

type userRepository struct {
	Database func(dbName string) *pgx.Conn
}

func NewUserRepository(database func(dbName string) *pgx.Conn) UserRepository {
	return &userRepository{
		Database: database,
	}
}

func (repository *userRepository) CreateUser(c context.Context, user domain.User) error {
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
		role,
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
		user.Role,
		user.CreatedAt,
		user.UpdatedAt); err != nil {
		return exception.ErrUnprocessableEntity(err.Error())
	}

	log.Printf("success insert user %s into DB", user.Username)
	return nil
}

func (repository *userRepository) GetAllUser(c context.Context) ([]domain.User, error) {
	ctx, cancel := context.WithTimeout(c, 10*time.Second)
	defer cancel()

	db := repository.Database(dbName)
	defer db.Close(ctx)

	query := `SELECT 
		id,
		name,
		username,
		email,
		password,
		phone,
		role,
		created_at,
		updated_at
		FROM users`

	user, err := db.Query(ctx, query)
	if err != nil {
		return []domain.User{}, exception.ErrInternalServer(err.Error())
	}

	defer user.Close()

	users, err := pgx.CollectRows(user, pgx.RowToStructByPos[domain.User])
	if err != nil {
		fmt.Printf("CollectRows error: %v", err)
		return []domain.User{}, exception.ErrInternalServer(err.Error())
	}

	return users, nil
}

func (repository *userRepository) GetUserByID(c context.Context, userID string) (domain.User, error) {
	ctx, cancel := context.WithTimeout(c, 10*time.Second)
	defer cancel()

	db := repository.Database(dbName)
	defer db.Close(ctx)

	query := `SELECT 
		id,
		name,
		username,
		email,
		password,
		phone,
		role,
		created_at,
		updated_at
		FROM users
		WHERE id = $1`

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

func (repository *userRepository) GetUsersByQuery(c context.Context, params string, value string) (domain.User, error) {
	ctx, cancel := context.WithTimeout(c, 10*time.Second)
	defer cancel()

	db := repository.Database(dbName)
	defer db.Close(ctx)

	query := fmt.Sprintf("SELECT * FROM users WHERE %s = $1", params)

	user := db.QueryRow(ctx, query, value)

	var data domain.User
	user.Scan(&data.ID, &data.Name, &data.Username, &data.Email, &data.Password, &data.Phone, &data.Role, &data.CreatedAt, &data.UpdatedAt)

	return data, nil
}

func (repository *userRepository) UpdateUser(c context.Context, user domain.User) error {
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

func (repository *userRepository) UpdateUserPassword(c context.Context, user domain.User) error {
	ctx, cancel := context.WithTimeout(c, 10*time.Second)
	defer cancel()

	db := repository.Database(dbName)
	defer db.Close(ctx)

	query := "UPDATE users SET password = $1 WHERE nik = $2"

	if _, err := db.Prepare(c, "data", query); err != nil {
		return exception.ErrInternalServer(err.Error())
	}

	if _, err := db.Exec(c, "data", user.Password, user.ID); err != nil {
		return exception.ErrInternalServer(err.Error())
	}

	return nil
}

func (repository *userRepository) DeleteUser(c context.Context, user_id string) error {
	ctx, cancel := context.WithTimeout(c, 10*time.Second)
	defer cancel()

	db := repository.Database(dbName)
	defer db.Close(ctx)

	query := `DELETE FROM users WHERE id = $1`

	if _, err := db.Prepare(c, "data", query); err != nil {
		return exception.ErrInternalServer(err.Error())
	}

	if _, err := db.Exec(c, "data", user_id); err != nil {
		return exception.ErrUnprocessableEntity(err.Error())
	}

	return nil
}
