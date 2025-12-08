package card

import (
	"canteen/internal/model"
	"database/sql"
)

type CardRepository interface {
	FindUserByCardNo(cardNo string) (*model.UserVo, error)
	FindOrderRecord(userId int, mealType string, weekNumber string, weekday string) (*model.OrderRecord, error)
	CreateOrderRecord(order *model.OrderRecord) error
	UpdateOrderStatus(orderId int, status string) error
	UpdateUserCount(userId int, count int) error
	GetCanteenConfigs() (flexibleDeptId, fixedDeptId, flexibleDinnerStart, fixedDinnerStart string, err error)
}

type cardRepository struct {
	db *sql.DB
}

func NewCardRepository(db *sql.DB) CardRepository {
	return &cardRepository{db: db}
}

func (r *cardRepository) FindUserByCardNo(cardNo string) (*model.UserVo, error) {
	var user model.UserVo
	err := r.db.QueryRow(
		"SELECT user_id, dept_id, nick_name, count, card_no FROM sys_user WHERE card_no = ?",
		cardNo,
	).Scan(&user.UserId, &user.DeptId, &user.NickName, &user.Count, &user.CardNo)
	return &user, err
}

func (r *cardRepository) FindOrderRecord(userId int, mealType string, weekNumber string, weekday string) (*model.OrderRecord, error) {
	var order model.OrderRecord
	err := r.db.QueryRow(`
		SELECT o.id, o.status, o.setmeal_id, o.user_id 
		FROM order_record o 
		WHERE o.user_id = ? AND o.meal_type = ? AND o.week_number = ? AND o.weekday = ?
	`, userId, mealType, weekNumber, weekday).Scan(&order.Id, &order.Status, &order.MealId, &order.UserId)
	return &order, err
}

func (r *cardRepository) CreateOrderRecord(order *model.OrderRecord) error {
	_, err := r.db.Exec(`
		INSERT INTO order_record 
		(user_id, week_number, order_date, weekday, meal_type, setmeal_id, quantity, status, create_time, update_time)
		VALUES (?, ?, ?, ?, ?, ?, 1, ?, NOW(), NOW())
	`,
		order.UserId,
		order.WeekNumber,
		order.OrderDate,
		order.Weekday,
		order.MealType,
		order.MealId,
		order.Status,
	)
	return err
}

func (r *cardRepository) UpdateOrderStatus(orderId int, status string) error {
	_, err := r.db.Exec("UPDATE order_record SET status = ? WHERE id = ?", status, orderId)
	return err
}

func (r *cardRepository) UpdateUserCount(userId int, count int) error {
	_, err := r.db.Exec("UPDATE sys_user SET count = ? WHERE user_id = ?", count, userId)
	return err
}

func (r *cardRepository) GetCanteenConfigs() (flexibleDeptId, fixedDeptId, flexibleDinnerStart, fixedDinnerStart string, err error) {
	err = r.db.QueryRow(`
		SELECT 
			(SELECT config_value FROM canteen_config WHERE config_key='flexible_dept_id'),
			(SELECT config_value FROM canteen_config WHERE config_key='fixed_dept_id'),
			(SELECT config_value FROM canteen_config WHERE config_key='flexible_dinner_start_time'),
			(SELECT config_value FROM canteen_config WHERE config_key='fixed_dinner_start_time')
	`).Scan(&flexibleDeptId, &fixedDeptId, &flexibleDinnerStart, &fixedDinnerStart)
	return flexibleDeptId, fixedDeptId, flexibleDinnerStart, fixedDinnerStart, err
}
