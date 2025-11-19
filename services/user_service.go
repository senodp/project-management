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
	//daftarkan login ke interface
	Login (email,password string)(*models.User,error)
	GetByID (id uint) (*models.User, error)
	GetByPublicID (id string) (*models.User, error)
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

//fungsi untuk login
func (s *userService) Login (email,password string)(*models.User,error){
	user, err := s.repo.FindByEmail(email)
	if err != nil{
		return nil,errors.New("invalid credential")
	}
	if !utils.CheckPasswordHash(password, user.Password){
		return nil, errors.New("invalid credential")
	}
	return user, nil
}

func (s *userService) GetByID (id uint) (*models.User, error){
	return s.repo.FindByID(id)
}

func (s *userService) GetByPublicID (id string) (*models.User, error){
	return s.repo.FindByPublicID(id)
}


