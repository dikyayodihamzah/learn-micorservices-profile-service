package main

import (
	"log"
	"time"

	"github.com/go-playground/validator"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"gitlab.com/learn-micorservices/user-service/config"
	"gitlab.com/learn-micorservices/user-service/controller"
	"gitlab.com/learn-micorservices/user-service/repository"
	"gitlab.com/learn-micorservices/user-service/service"
)

func controllers() {
	time.Local = time.UTC

	serverConfig := config.NewServerConfig()
	db := config.NewDB
	validate := validator.New()

	userRepository := repository.NewUserRepository(db)
	roleRepository := repository.NewRoleRepository(db)
	userService := service.NewUserService(userRepository, roleRepository, validate)
	userController := controller.NewUserController(userService)

	app := fiber.New()
	app.Use(recover.New())
	app.Use(cors.New(cors.Config{
		AllowOrigins:     "*",
		AllowMethods:     "*",
		AllowHeaders:     "*",
		AllowCredentials: true,
	}))

	userController.NewUserRouter(app)

	err := app.Listen(serverConfig.URI)
	log.Println(err)
}

func main() {
	time.Local = time.UTC
	controllers()
}
