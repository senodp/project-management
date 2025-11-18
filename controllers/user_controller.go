package controllers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/jinzhu/copier"
	"github.com/senodp/project-management/models"
	"github.com/senodp/project-management/services"
	"github.com/senodp/project-management/utils"
)

//buat struct
type UserController struct{
	service services.UserService
}

func NewUserController(s services.UserService)*UserController{
	return &UserController{service: s}
}

func (c *UserController) Register(ctx *fiber.Ctx) error {
	user := new(models.User)

	if err := ctx.BodyParser(user); err != nil {
		return utils.BadRequest(ctx, "Gagal Parsing data", err.Error())
	}

	if err := c.service.Register(user); err != nil {
		return utils.BadRequest(ctx, "Registrasi Gagal", err.Error())
	}

	var userResp models.UserResponse
	_ = copier.Copy(&userResp,&user)
	return utils.Success(ctx, "Register Success!", userResp)
}