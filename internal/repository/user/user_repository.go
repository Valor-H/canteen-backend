package user

import (
	"canteen/internal/model"
	"database/sql"
)

type UserRepository interface {
	FindByCardNo(cardNo string) (*model.UserVo, error)
	DecreaseCountByUserId(userId int) error
}

type userRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) UserRepository {
	return &userRepository{db: db}
}

func (r *userRepository) FindByCardNo(cardNo string) (*model.UserVo, error) {
	var user model.UserVo
	err := r.db.QueryRow("SELECT user_id, dept_id, nick_name, count, card_no FROM sys_user WHERE card_no = ?", cardNo).
		Scan(&user.Id, &user.DeptId, &user.NickName, &user.Count, &user.CardNo)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *userRepository) DecreaseCountByUserId(userId int) error {
	_, err := r.db.Exec("UPDATE sys_user SET count = GREATEST(count - 1, 0) WHERE user_id = ?", userId)
	return err
}