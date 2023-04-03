package service

import (
	"context"
	"time"

	"github.com/go-playground/validator"
	"gitlab.com/learn-micorservices/user-service/exception"
	"gitlab.com/learn-micorservices/user-service/helper"
	"gitlab.com/learn-micorservices/user-service/model/domain"
	"gitlab.com/learn-micorservices/user-service/model/web"
	"gitlab.com/learn-micorservices/user-service/repository"
)

type UserService interface {
	CreateUser(c context.Context, claims helper.JWTClaims, request web.CreateUserRequest) (web.UserResponse, error)
	GetAllUser(c context.Context) ([]web.UserResponse, error)
	GetUserByID(c context.Context, claims helper.JWTClaims) (web.UserResponse, error)
	UpdateUser(c context.Context, claims helper.JWTClaims, request web.UpdateUserRequest) (web.UserResponse, error)
	UpdateUserPassword(c context.Context, claims helper.JWTClaims, request web.UpdatePasswordRequest) (web.UserResponse, error)
	DeleteUser(c context.Context, claims helper.JWTClaims) error
}

type userService struct {
	UserRepository repository.UserRepository
	RoleRepository repository.RoleRepository
	Validate       *validator.Validate
}

func NewUserService(userRepository repository.UserRepository, roleRepository repository.RoleRepository, validate *validator.Validate) UserService {
	return &userService{
		UserRepository: userRepository,
		RoleRepository: roleRepository,
		Validate:       validate,
	}
}

func (service *userService) CreateUser(c context.Context, claims helper.JWTClaims, request web.CreateUserRequest) (web.UserResponse, error) {
	if err := service.Validate.Struct(request); err != nil {
		return web.UserResponse{}, exception.ErrBadRequest(err.Error())
	}

	user := domain.User{
		Name:     request.Name,
		Username: request.Username,
		Email:    request.Email,
		Phone:    request.Phone,
		Role: domain.Role{
			ID: request.RoleID,
		},
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	user.GenerateID()
	user.SetPassword(request.Password)

	err := service.UserRepository.CreateUser(c, user)
	if err != nil {
		return web.UserResponse{}, err
	}

	user, err = service.UserRepository.GetUserByID(c, user.ID)
	if err != nil {
		return web.UserResponse{}, err
	}

	// KAFKA

	return helper.ToUserResponse(user), nil
}

func (service *userService) GetAllUser(c context.Context) ([]web.UserResponse, error) {
	users, err := service.UserRepository.GetAllUser(c)
	if err != nil {
		return []web.UserResponse{}, err
	}
	return helper.ToAllUserResponses(users), nil
}

func (service *userService) GetUserByID(c context.Context, claims helper.JWTClaims) (web.UserResponse, error) {
	user, err := service.UserRepository.GetUserByID(c, claims.User.Id)
	if err != nil {
		return web.UserResponse{}, err
	}

	if user.ID == "" {
		return web.UserResponse{}, exception.ErrNotFound("user not found")
	}
	return helper.ToUserResponse(user), nil
}

func (service *userService) UpdateUser(c context.Context, claims helper.JWTClaims, request web.UpdateUserRequest) (web.UserResponse, error) {
	if err := service.Validate.Struct(request); err != nil {
		return web.UserResponse{}, exception.ErrBadRequest(err.Error())
	}

	user, err := service.UserRepository.GetUserByID(c, claims.User.Id)
	if err != nil {
		return web.UserResponse{}, exception.ErrNotFound(err.Error())
	}

	if request.Email != "" {
		if userByEmail, _ := service.UserRepository.GetUsersByQuery(c, "email", request.Email); 
		userByEmail.ID != "" && userByEmail.ID != claims.User.Id {
			exception.ErrBadRequest("email already registered")
		}
		user.Email = request.Email
	}

	if request.Username != "" {
		if userByUsername, _ := service.UserRepository.GetUsersByQuery(c, "username", request.Username); 
		userByUsername.ID != "" && userByUsername.ID != claims.User.Id {
			exception.ErrBadRequest("username already registered")
		}
		user.Username = request.Username
	}

	if request.Phone != "" {
		if userByPhone, _ := service.UserRepository.GetUsersByQuery(c, "phone", request.Phone); 
		userByPhone.ID != "" && userByPhone.ID != claims.User.Id {
			exception.ErrBadRequest("phone already registered")
		}
		user.Phone = request.Phone
	}

	if request.RoleID != "" {
		if role := service.RoleRepository.GetRoleByID(c, request.RoleID); role.ID == "" {
			exception.ErrBadRequest("role not found")
		}
		user.Email = request.Email
	}

	user.UpdatedAt = time.Now()

	if err := service.UserRepository.UpdateUser(c, user); err != nil {
		return web.UserResponse{}, exception.ErrInternalServer(err.Error())
	}

	// KAFAK

	user, _ = service.UserRepository.GetUserByID(c, claims.User.Id)
	return helper.ToUserResponse(user), nil
}

func (service *userService) UpdateUserPassword(c context.Context, claims helper.JWTClaims, request web.UpdatePasswordRequest) (web.UserResponse, error) {	
	if err := service.Validate.Struct(request); err != nil {
		return web.UserResponse{}, exception.ErrBadRequest(err.Error())
	}

	user, err := service.UserRepository.GetUserByID(c, claims.User.Id)
	if err != nil || user.ID == "" {
		return web.UserResponse{}, exception.ErrNotFound("user does not exist")
	}

	if request.Password != request.ConfirmPassword {
		return web.UserResponse{}, exception.ErrBadRequest("password not match")
	}

	user.SetPassword(request.Password)
	
	user.UpdatedAt = time.Now()
	
	if err := service.UserRepository.UpdateUserPassword(c, user); err != nil {
		return web.UserResponse{}, err
	}

	// KAFKA

	return helper.ToUserResponse(user), nil
}

func (service *userService) DeleteUser(c context.Context, claims helper.JWTClaims) error {
	user, err := service.UserRepository.GetUserByID(c, claims.User.Id)
	if err != nil || user.ID == "" {
		return exception.ErrNotFound("user not found")
	}

	err = service.UserRepository.DeleteUser(c, claims.User.Id)
	if err != nil {
		return err
	}

	// KAFKA

	return nil
}
