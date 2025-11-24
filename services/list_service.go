package services

import (
	"github.com/google/uuid"
	"github.com/senodp/project-management/models"
	"github.com/senodp/project-management/repositories"
)

type listService struct {
	listRepo repositories.ListRepository
	boardRepo repositories.BoardRepository
	listPosRepo repositories.ListPositionRepository
}

type ListWithOrder struct {
	Position []uuid.UUID
	Lists []models.List
}

type ListService interface {
	GetByBoardID(boardPublicID string) (*ListWithOrder,error)
	GetByID(id uint) (*models.List, error)
	GetByPublicID(publicID string) (*models.List,error)
	Create(list *models.List) error
	Update(list *models.List) error
	Delete(id uint) error
	UpdatePositions(boardPublicID string, positions []uuid.UUID) error 
}