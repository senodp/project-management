package services

import (
	"errors"

	"github.com/google/uuid"
	"github.com/senodp/project-management/models"
	"github.com/senodp/project-management/repositories"
	"github.com/senodp/project-management/utils"
)

type UserService interface {
	Register(user *models.User) error
}

type userService struct{
	repo repositories.UserRepository
}

func NewUserService(repo repositories.UserRepository) UserService {
	return &userService{repo}
}

func (s *userService)Register(user *models.User) error {
	//harus mengecek email yang terdaftar apakah sudah digunakan atau belum
	//hasing password
	//set role user
	//simpan user

	existingUser, _ := s.repo.FindByEmail(user.Email)
	if existingUser.InternalID != 0 {
		return errors.New("email already registered")
	}
	
	hased, err := utils.HashPassword(user.Password)
	if err != nil {
		return err
	}

	user.Password = hased
	user.Role = "user"
	user.PublicID = uuid.New()

	return s.repo.Create(user)
}

