package services

import (
	"errors"

	"github.com/google/uuid"
	"github.com/senodp/project-management/models"
	"github.com/senodp/project-management/repositories"
)

type BoardService interface {
	Create(board *models.Board) error
	Update(board *models.Board) error
	GetByPublicID(publicID string)(*models.Board, error)
	AddMembers(boardPublicID string, userPublicIDS []string)error
	RemoveMembers(boardPublicID string, userPublicIDs []string)error
	GetAllByUserPaginate(userID, filter, sort string, limit, 
		offset int)([]models.Board, int64, error)
}

type boardService struct {
	boardRepo repositories.BoardRepository
	userRepo repositories.UserRepository
	boardMemberRepo repositories.BoardMemberRepository
}

func NewBoardService(
	boardRepo repositories.BoardRepository,
	userRepo repositories.UserRepository,
	boardMemberRepo repositories.BoardMemberRepository,
	) BoardService {
	return &boardService{boardRepo, userRepo, boardMemberRepo}
}

func (s *boardService) Create(board *models.Board) error {
	user, err := s.userRepo.FindByPublicID(board.OwnerPublicID.String())
	if err != nil {
		return errors.New("owner not found")
	}
	board.PublicID = uuid.New()
	board.OwnerID = user.InternalID
	return s.boardRepo.Create(board)
}

func (s *boardService) Update(board *models.Board) error {
	return s.boardRepo.Update(board)
}

func (s *boardService) GetByPublicID(publicID string)(*models.Board, error){
	return s.boardRepo.FindByPublicID(publicID)
}

func (s *boardService) AddMembers(boardPublicID string, userPublicIDS []string)error{
	board, err := s.boardRepo.FindByPublicID(boardPublicID)
	if err != nil{
		return errors.New("board not found")
	}

	var userInternalIDs []uint
	for _, userPublicID := range userPublicIDS{
		user, err := s.userRepo.FindByPublicID(userPublicID)
		if err != nil{
			return errors.New("user not found"+ userPublicID)
		}
		//kumpulkan IDnya ke user internal id
		userInternalIDs = append(userInternalIDs, uint(user.InternalID))
	}
	//cek keanggotaan sekaligus
	existingMembers, err := s.boardMemberRepo.GetMembers(string(board.PublicID.String()))
	if err != nil{
		return err
	}

	//cek cepat dengna map
	memberMap := make(map[uint]bool)
	for _, member := range existingMembers{
		memberMap[uint(member.InternalID)] = true //memberMap[1] = true
	}

	var newMembersIDs []uint
	for _,userID := range userInternalIDs{
		if !memberMap[userID]{
			newMembersIDs = append(newMembersIDs, userID)
		}
	}
	if len(newMembersIDs) == 0 {
		return nil
	}
	return s.boardRepo.AddMember(uint(board.InternalID), newMembersIDs)
}

func (s *boardService) RemoveMembers (boardPublicID string, userPublicIDs []string)error{
	//cek apakah ada atau tidak
	board, err := s.boardRepo.FindByPublicID(boardPublicID)
	if err != nil{
		return errors.New("board not found")
	}

	//validasi user apakah dia user kita atau tidak
	var userInternalIDs []uint
	for _, userPublicID := range userPublicIDs{
		user, err := s.userRepo.FindByPublicID(userPublicID)
		if err != nil{
			return errors.New("User Not Found! "+ userPublicID)
		}
		//kumpulkan IDnya ke user internal id
		userInternalIDs = append(userInternalIDs, uint(user.InternalID))
	}

	//cek keanggotaan
	existingMembers, err := s.boardMemberRepo.GetMembers(string(board.PublicID.String()))
	if err != nil{
		return err
	}

	//cek cepat dengna map
	memberMap := make(map[uint]bool)
	for _, member := range existingMembers{
		memberMap[uint(member.InternalID)] = true //memberMap[1] = true
	}

	var memberToRemove []uint
	for _,userID := range userInternalIDs{
		if memberMap[userID]{
			memberToRemove = append(memberToRemove, userID)
		}
	}
	return s.boardRepo.RemoveMembers(uint(board.InternalID), memberToRemove)
}

func (s *boardService) GetAllByUserPaginate (userID, filter, sort string, limit, 
	offset int)([]models.Board, int64, error){
		return s.boardRepo.FindAllByUserPagination(userID,filter,sort,limit,offset)
	}
