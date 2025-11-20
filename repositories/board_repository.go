package repositories

import (
	"github.com/senodp/project-management/config"
	"github.com/senodp/project-management/models"
)

type BoardRepository interface {
	Create(board *models.Board) error
}

type boardRepository struct {
}

func NewBoardRepository() BoardRepository {
	//ambil alamat dari struct
	return &boardRepository{}
}

func (r *boardRepository) Create(board *models.Board) error {
	return config.DB.Create(board).Error
}