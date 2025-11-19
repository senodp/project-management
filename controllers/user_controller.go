package controllers

import (
	"math"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
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

//buat fungsi login
func (c *UserController) Login(ctx *fiber.Ctx)error{
	//buat struct unutk menerima input dari client / end point body
	var body struct{
		Email string `json:"email"`
		Password string `json:"password"`
	}

	if err := ctx.BodyParser(&body);err != nil{
		return utils.BadRequest(ctx, "Invalid Request", err.Error())
	}

	user, err := c.service.Login(body.Email,body.Password)
	if err != nil{
		return utils.Unauthorized(ctx, "Login Failed", err.Error())
	}

	token ,_ := utils.GenerateToken(user.InternalID, user.Role, user.Email, user.PublicID)
	refreshToken ,_ := utils.GenerateRefreshToken(user.InternalID) 

	var userResp models.UserResponse
	_ = copier.Copy(&userResp,&user)
	return utils.Success(ctx, "Login Succesful!", fiber.Map{
		"access_token": token,
		"refresh_token": refreshToken,
		"user": userResp, 
	})
}

func (c *UserController) GetUser(ctx *fiber.Ctx) error{
	id := ctx.Params("id")
	user, err := c.service.GetByPublicID(id)
	if err != nil{
		return utils.NotFound(ctx, "Data Not Found", err.Error())
	}

	var userResp models.UserResponse
	err = copier.Copy(&userResp, &user)
	if err != nil{
		return utils.BadRequest(ctx, "Internal Server Error", err.Error())
	}
	return  utils.Success(ctx, "Data Berhasil Ditemukan", userResp)
}

func (c *UserController) GetUserPagination (ctx *fiber.Ctx) error {
	// /users/page?page=1&limit=10&sort=-id&filter=alfiani
	// 100 /10 = 10 page
	page,_ := strconv.Atoi(ctx.Query("page", "1"))
	limit,_ := strconv.Atoi(ctx.Query("limit", "10"))
	offset := (page - 1) * limit

	filter := ctx.Query("filter", "")
	sort := ctx.Query("sort", "")

	user, total, err := c.service.GetAllpagination(filter, sort, limit, offset)
	if err != nil{
		return utils.BadRequest(ctx, "Gagal mengambil data!", err.Error())
	}

	var userResp []models.UserResponse
	_ = copier.Copy(&userResp, &user)

	meta := utils.PaginationMeta{
		Page : page,
		Limit : limit,
		Total : int(total),
		TotalPage :  int(math.Ceil(float64(total)/float64(limit))),
		Filter : filter,
		Sort: sort,
	}

	if total == 0{
		return utils.NotFoundPagination(ctx, "Data pengguna tidak ditemukan!", userResp, meta)
	}

	return utils.SuccessPagination(ctx, "Data ditemukan", userResp, meta)
}

func (c *UserController) UpdateUser (ctx *fiber.Ctx) error {
	id := ctx.Params("id")
	publicID, err := uuid.Parse(id)
	if err != nil{
		return utils.BadRequest(ctx, "Invalid ID Format", err.Error())
	}

	var user models.User
	if err := ctx.BodyParser(&user); err != nil {
		return utils.BadRequest(ctx, "Gagal Parsing Data", err.Error())
	}

	user.PublicID = publicID

	if err := c.service.Update(&user); err != nil {
		return utils.BadRequest(ctx, "Gagal Update Data", err.Error())
	}

	userUpdated , err := c.service.GetByPublicID(id)
	if err != nil{
		return utils.InternalServerError(ctx, "Gagal Ambil Data!", err.Error())
	}

	var userResp models.UserResponse
	err = copier.Copy(&userResp,&userUpdated)
	if err != nil{
		return utils.InternalServerError(ctx, "Error Parsing Data", err.Error())
	}

	return utils.Success(ctx, "Berhasil Update Data!", userResp)
}

func (c *UserController) DeleteUser(ctx *fiber.Ctx) error{
	id, _ := strconv.Atoi(ctx.Params("id"))
	if err := c.service.Delete(uint(id)); err != nil{
		return utils.InternalServerError(ctx, "Gagal Menghapus Data", err.Error())
	}
	return utils.Success(ctx, "Berhasil Menghapus Data!", id)
}