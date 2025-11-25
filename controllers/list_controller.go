package controllers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/senodp/project-management/models"
	"github.com/senodp/project-management/services"
	"github.com/senodp/project-management/utils"
)

type ListController struct {
	services services.ListService
}

func NewListController(s services.ListService) *ListController {
	return &ListController{services: s}
}

func (c *ListController) CreateList(ctx *fiber.Ctx) error {
	list := new(models.List)
	if err := ctx.BodyParser(list); err != nil {
		return utils.BadRequest(ctx, "Gagal Membaca Request", err.Error())
	}
	if err := c.services.Create(list); err != nil {
		return utils.BadRequest(ctx, "Gagal Membuat List", err.Error())
	}
	return utils.Success(ctx, "List Berhasil Di buat", list)
}

func (c *ListController) UpdateList(ctx *fiber.Ctx) error {
	publicID := ctx.Params("id")
	list := new(models.List)

	if err := ctx.BodyParser(list); err != nil {
		return utils.BadRequest(ctx, "Gagal parsing data", err.Error())
	}

	if _, err := uuid.Parse(publicID); err != nil {
		return utils.BadRequest(ctx, "ID tidak valid", err.Error())
	}

	existingList, err := c.services.GetByPublicID(publicID)
	if err != nil {
		return utils.NotFound(ctx, "List Tidak Di Temukan", err.Error())
	}
	list.InternalID = existingList.InternalID
	list.PublicID = existingList.PublicID

	if err := c.services.Update(list); err != nil {
		return utils.BadRequest(ctx, "Gagal Update List", err.Error())
	}

	updatedList, err := c.services.GetByPublicID(publicID)
	if err != nil {
		return utils.NotFound(ctx, "List Tidak Di Temukan", err.Error())
	}

	return  utils.Success(ctx, "Berhasil Memperbarui List", updatedList)
}

func (c *ListController) GetListOnBoard(ctx *fiber.Ctx) error {
	boardPublicID := ctx.Params("board_id")
	if _, err := uuid.Parse(boardPublicID); err != nil {
		return utils.BadRequest(ctx, "ID tidak valid", err.Error())
	}

	lists, err := c.services.GetByBoardID(boardPublicID)
	if err != nil {
		return utils.NotFound(ctx, "List Tidak Di temmukan", err.Error())
	}

	return utils.Success(ctx, "Data Berhasil Di Ambil",lists)
}

func (c *ListController) DeleteList(ctx *fiber.Ctx) error {
	publicID := ctx.Params("id")
	if _, err := uuid.Parse(publicID); err != nil {
		return utils.BadRequest(ctx, "ID tidak valid", err.Error())
	}

	list,err :=  c.services.GetByPublicID(publicID)
	if err != nil {
		return utils.NotFound(ctx,"List Tidak Di Temukan",err.Error())
	}

	if err :=  c.services.Delete(uint(list.InternalID)); err != nil {
		return utils.InternalServerError(ctx, "Gagal Menghapus List",err.Error())
	}

	return utils.Success(ctx,"List Berhasil Di Hapus",publicID)

}

