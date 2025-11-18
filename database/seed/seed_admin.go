package seed

import (
	"log"

	"github.com/google/uuid"
	"github.com/senodp/project-management/config"
	"github.com/senodp/project-management/models"
	"github.com/senodp/project-management/utils"
)

func SeedAdmin() {
	password, _ := utils.HashPassword("admin")

	admin := models.User{
		Name: "Administrator",
		Email: "administrator@gmail.com",
		Password: password,
		Role: "admin",
		PublicID: uuid.New(),
	}
	if err:=config.DB.FirstOrCreate(&admin,models.User{Email:admin.Email}).Error; err != nil{
		log.Println("Failed too seed Admin")
	} else {
		log.Println("Admin user seeded")
	}
}