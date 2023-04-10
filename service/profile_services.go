package service

import (
	"context"
	"time"

	"github.com/go-playground/validator"
	"gitlab.com/learn-micorservices/profile-service/exception"
	"gitlab.com/learn-micorservices/profile-service/helper"
	"gitlab.com/learn-micorservices/profile-service/model/web"
	"gitlab.com/learn-micorservices/profile-service/repository"
)

type ProfileService interface {
	GetCurrentProfile(c context.Context, claims helper.JWTClaims) (web.ProfileResponse, error)
	UpdateProfile(c context.Context, claims helper.JWTClaims, request web.UpdateProfileRequest) (web.ProfileResponse, error)
	UpdatePassword(c context.Context, claims helper.JWTClaims, request web.UpdatePasswordRequest) (web.ProfileResponse, error)
}

type profileService struct {
	ProfileRepository repository.ProfileRepository
	RoleRepository    repository.RoleRepository
	Validate          *validator.Validate
}

func NewProfileService(profileRepository repository.ProfileRepository, roleRepository repository.RoleRepository, validate *validator.Validate) ProfileService {
	return &profileService{
		ProfileRepository: profileRepository,
		RoleRepository:    roleRepository,
		Validate:          validate,
	}
}

func (service *profileService) GetCurrentProfile(c context.Context, claims helper.JWTClaims) (web.ProfileResponse, error) {
	user, err := service.ProfileRepository.GetProfileByID(c, claims.User.ID)
	if err != nil {
		return web.ProfileResponse{}, err
	}

	if user.ID == "" {
		return web.ProfileResponse{}, exception.ErrNotFound("user not found")
	}
	return helper.ToProfileResponse(user), nil
}

func (service *profileService) UpdateProfile(c context.Context, claims helper.JWTClaims, request web.UpdateProfileRequest) (web.ProfileResponse, error) {
	if err := service.Validate.Struct(request); err != nil {
		return web.ProfileResponse{}, exception.ErrBadRequest(err.Error())
	}

	user, err := service.ProfileRepository.GetProfileByID(c, claims.User.ID)
	if err != nil {
		return web.ProfileResponse{}, exception.ErrNotFound(err.Error())
	}

	if request.Name != "" {
		user.Name = request.Name
	}

	if request.Username != "" {
		if userByUsername, _ := service.ProfileRepository.GetProfilesByQuery(c, "email", request.Username); userByUsername.ID != "" && userByUsername.ID != claims.User.ID {
			exception.ErrBadRequest("username already registered")
		}
		user.Username = request.Username
	}

	if request.Email != "" {
		if userByEmail, _ := service.ProfileRepository.GetProfilesByQuery(c, "email", request.Email); userByEmail.ID != "" && userByEmail.ID != claims.User.ID {
			exception.ErrBadRequest("email already registered")
		}
		user.Email = request.Email
	}

	if request.Phone != "" {
		if userByPhone, _ := service.ProfileRepository.GetProfilesByQuery(c, "phone", request.Phone); userByPhone.ID != "" && userByPhone.ID != claims.User.ID {
			exception.ErrBadRequest("phone already registered")
		}
		user.Phone = request.Phone
	}

	user.UpdatedAt = time.Now()

	if err := service.ProfileRepository.UpdateProfile(c, user); err != nil {
		return web.ProfileResponse{}, exception.ErrInternalServer(err.Error())
	}

	// KAFKA
	helper.ProduceToKafka(user, "PUT.USER", helper.KafkaTopic)

	userRes, _ := service.ProfileRepository.GetProfileByID(c, user.ID)
	return helper.ToProfileResponse(userRes), nil
}

func (service *profileService) UpdatePassword(c context.Context, claims helper.JWTClaims, request web.UpdatePasswordRequest) (web.ProfileResponse, error) {
	if err := service.Validate.Struct(request); err != nil {
		return web.ProfileResponse{}, exception.ErrBadRequest(err.Error())
	}

	user, err := service.ProfileRepository.GetProfileByID(c, claims.User.ID)
	if err != nil || user.ID == "" {
		return web.ProfileResponse{}, exception.ErrNotFound("user does not exist")
	}

	if request.Password != request.ConfirmPassword {
		return web.ProfileResponse{}, exception.ErrBadRequest("password did not match")
	}

	user.SetPassword(request.Password)

	user.UpdatedAt = time.Now()

	if err := service.ProfileRepository.UpdatePassword(c, user); err != nil {
		return web.ProfileResponse{}, err
	}

	// KAFKA
	helper.ProduceToKafka(user, "PUT.USER", helper.KafkaTopic)

	userRes, _ := service.ProfileRepository.GetProfileByID(c, user.ID)
	return helper.ToProfileResponse(userRes), nil
}
