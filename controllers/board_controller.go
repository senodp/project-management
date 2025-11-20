package controllers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/senodp/project-management/models"
	"github.com/senodp/project-management/services"
	"github.com/senodp/project-management/utils"
)

type BoardController struct {
	service services.BoardService
}

func NewBoardController (s services.BoardService) *BoardController {
	return &BoardController{service: s}
}

func (c *BoardController) CreateBoard(ctx *fiber.Ctx) error {
	board := new(models.Board)

	if err := ctx.BodyParser(board); err != nil {
		return utils.BadRequest(ctx, "Gagal membaca request", err.Error())
	}

	if err := c.service.Create(board); err != nil {
		return utils.BadRequest(ctx, "Gagal Menyimpan Data", err.Error())
	}

	return utils.Success(ctx, "Board berhasil dibuat", board)
}