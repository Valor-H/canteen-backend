package user

import (
	"canteen/internal/model"
	"canteen/internal/repository/user"
	"errors"
)

type UserService interface {
	FindByCardNo(cardNo string) (*model.UserVo, error)
	FindById(userId int) (*model.UserVo, error)
	FindByNickName(nickName string) (*model.UserVo, error)
	DecreaseUserCount(userId int) error
}

type userService struct {
	userRepo user.UserRepository
}

func NewUserService(userRepo user.UserRepository) UserService {
	return &userService{userRepo: userRepo}
}

func (s *userService) FindByCardNo(cardNo string) (*model.UserVo, error) {
	if cardNo == "" {
		return nil, errors.New("卡号不能为空")
	}
	return s.userRepo.FindByCardNo(cardNo)
}

func (s *userService) FindById(userId int) (*model.UserVo, error) {
	if userId <= 0 {
		return nil, errors.New("无效的用户ID")
	}
	return s.userRepo.FindById(userId)
}

func (s *userService) FindByNickName(nickName string) (*model.UserVo, error) {
	if nickName == "" {
		return nil, errors.New("昵称不能为空")
	}
	return s.userRepo.FindByNickName(nickName)
}

func (s *userService) DecreaseUserCount(userId int) error {
	if userId <= 0 {
		return errors.New("无效的用户ID")
	}
	return s.userRepo.DecreaseCountByUserId(userId)
}