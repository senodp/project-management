package main

import (
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/senodp/project-management/config"
	"github.com/senodp/project-management/controllers"
	"github.com/senodp/project-management/database/seed"
	"github.com/senodp/project-management/repositories"
	"github.com/senodp/project-management/routes"
	"github.com/senodp/project-management/services"
)

func main() {
	config.LoadEnv()
	config.ConnectDB()

	seed.SeedAdmin()
	app := fiber.New()

	userRepo := repositories.NewUserRepository()
	userService := services.NewUserService(userRepo)
	userController := controllers.NewUserController(userService)

	routes.Setup(app,userController)

	port := config.AppConfig.AppPort
	log.Println("Server is running on port :", port)
	log.Fatal(app.Listen(":" + port))
}