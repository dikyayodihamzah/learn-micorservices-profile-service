package controller

import (
	"github.com/gofiber/fiber/v2"
	"gitlab.com/learn-micorservices/profile-service/config"
	"gitlab.com/learn-micorservices/profile-service/exception"
	"gitlab.com/learn-micorservices/profile-service/helper"

	"gitlab.com/learn-micorservices/profile-service/middleware"
	"gitlab.com/learn-micorservices/profile-service/model/web"
	"gitlab.com/learn-micorservices/profile-service/service"
)

type ProfileController interface {
	NewProfileRouter(app *fiber.App)
}
type profileController struct {
	ProfileService service.ProfileService
}

func NewProfileController(profileService service.ProfileService) ProfileController {
	return &profileController{
		ProfileService: profileService,
	}
}

func (controller *profileController) NewProfileRouter(app *fiber.App) {
	user := app.Group(config.EndpointPrefix)

	user.Get("/ping", func(ctx *fiber.Ctx) error {
		return ctx.Status(fiber.StatusOK).JSON(web.WebResponse{
			Code:    fiber.StatusOK,
			Status:  true,
			Message: "ok",
		})
	})

	user.Use(middleware.IsAuthenticated)
	user.Get("/", controller.GetCurrentProfile)
	user.Put("/update", controller.UpdateProfile)
	user.Put("/change-password", controller.UpdatePassword)
}

func (controller *profileController) GetCurrentProfile(ctx *fiber.Ctx) error {
	claims := ctx.Locals("claims").(helper.JWTClaims)

	user, err := controller.ProfileService.GetCurrentProfile(ctx.Context(), claims)
	if err != nil {
		return exception.ErrorHandler(ctx, err)
	}

	return ctx.Status(fiber.StatusOK).JSON(web.WebResponse{
		Code:    fiber.StatusOK,
		Status:  true,
		Message: "success",
		Data:    user,
	})
}

func (controller *profileController) UpdateProfile(ctx *fiber.Ctx) error {
	claims := ctx.Locals("claims").(helper.JWTClaims)

	request := new(web.UpdateProfileRequest)
	if err := ctx.BodyParser(request); err != nil {
		return exception.ErrorHandler(ctx, err)
	}

	user, err := controller.ProfileService.UpdateProfile(ctx.Context(), claims, *request)
	if err != nil {
		return exception.ErrorHandler(ctx, err)
	}

	return ctx.Status(fiber.StatusOK).JSON(web.WebResponse{
		Code:    fiber.StatusOK,
		Status:  true,
		Message: "success",
		Data:    user,
	})
}

func (controller *profileController) UpdatePassword(ctx *fiber.Ctx) error {
	claims := ctx.Locals("claims").(helper.JWTClaims)

	request := new(web.UpdatePasswordRequest)
	if err := ctx.BodyParser(request); err != nil {
		return exception.ErrorHandler(ctx, err)
	}

	user, err := controller.ProfileService.UpdatePassword(ctx.Context(), claims, *request)
	if err != nil {
		return exception.ErrorHandler(ctx, err)
	}

	return ctx.Status(fiber.StatusOK).JSON(web.WebResponse{
		Code:    fiber.StatusOK,
		Status:  true,
		Message: "success",
		Data:    user,
	})
}
