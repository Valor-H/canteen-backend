package meal

import (
	"canteen/internal/model"
	"database/sql"
	"time"
)

type MealRepository interface {
	FindSetmealsByWeekNumber(weekNumber string) ([]model.WeekMeal, error)
	FindDishesByIds(dishIds []int16) ([]model.Dish, error)
	GenerateWeeklySetmeals(dates []time.Time) error
	DeleteWeeklySetmeals(startWeek, endWeek string) error
	InsertWeeklySetmeal(weekNumber string, weekday string, mealType string, remark string) error
}

type mealRepository struct {
	db *sql.DB
}

func NewMealRepository(db *sql.DB) MealRepository {
	return &mealRepository{db: db}
}

func (r *mealRepository) FindSetmealsByWeekNumber(weekNumber string) ([]model.WeekMeal, error) {
	query := `
		SELECT id, week_number, weekday, meal_type FROM weekly_setmeal
		WHERE week_number = ? AND meal_type IN ('午餐', '晚餐')
	`
	
	rows, err := r.db.Query(query, weekNumber)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	
	var setmeals []model.WeekMeal
	for rows.Next() {
		var setmeal model.WeekMeal
		if err := rows.Scan(&setmeal.MealId, &setmeal.WeekNumber, &setmeal.WeekDay, &setmeal.MealType); err != nil {
			continue
		}
		setmeals = append(setmeals, setmeal)
	}
	
	return setmeals, rows.Err()
}

func (r *mealRepository) FindDishesByIds(dishIds []int16) ([]model.Dish, error) {
	if len(dishIds) == 0 {
		return []model.Dish{}, nil
	}
	
	placeholders, args := buildInClause(dishIds)
	query := `select DISTINCT setmeal_id from setmeal t1 left join setmeal_dish t2 on t1.id = t2.setmeal_id
                           left join dish t3 on t2.dish_id = t3.id where dish_id in(` + placeholders + `)`
	
	rows, err := r.db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	
	var dishes []model.Dish
	for rows.Next() {
		var dish model.Dish
		if err := rows.Scan(&dish.Id, &dish.Name); err != nil {
			continue
		}
		dishes = append(dishes, dish)
	}
	
	return dishes, rows.Err()
}

func (r *mealRepository) GenerateWeeklySetmeals(dates []time.Time) error {
	// 开启事务
	tx, err := r.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()
	
	// 删除旧记录
	startWeek := dates[0].Format("20060102")
	endWeek := dates[len(dates)-1].Format("20060102")
	if _, err := tx.Exec("DELETE FROM weekly_setmeal WHERE week_number BETWEEN ? AND ?", startWeek, endWeek); err != nil {
		return err
	}
	
	// 插入新记录
	for _, date := range dates {
		weekNumber := date.Format("20060102")
		weekday := getWeekdayZh(date)
		
		// 定义午餐和晚餐的备注
		lunchRemarks := []string{"套餐A", "套餐B", "套餐C"}
		dinnerRemark := []string{"套餐A", "套餐C"}
		
		// 插入3个午餐
		for _, remark := range lunchRemarks {
			if err := r.InsertWeeklySetmeal(weekNumber, weekday, "午餐", remark); err != nil {
				return err
			}
		}
		
		// 插入晚餐
		for _, remark := range dinnerRemark {
			if err := r.InsertWeeklySetmeal(weekNumber, weekday, "晚餐", remark); err != nil {
				return err
			}
		}
	}
	
	// 提交事务
	return tx.Commit()
}

func (r *mealRepository) InsertWeeklySetmeal(weekNumber string, weekday string, mealType string, remark string) error {
	_, err := r.db.Exec(`
		INSERT INTO weekly_setmeal 
			(week_number, weekday, meal_type, setmeal_id, create_time, create_user, remark)
		VALUES (?, ?, ?, NULL, NOW(), 263, ?)`,
		weekNumber, weekday, mealType, remark)
	return err
}

func (r *mealRepository) DeleteWeeklySetmeals(startWeek, endWeek string) error {
	_, err := r.db.Exec("DELETE FROM weekly_setmeal WHERE week_number BETWEEN ? AND ?", startWeek, endWeek)
	return err
}

// buildInClause 构造 SQL 中 IN 子句的 (?, ?, ...) 和参数列表
func buildInClause(ids []int16) (string, []interface{}) {
	placeholders := ""
	args := make([]interface{}, 0, len(ids))
	for i, id := range ids {
		if i == 0 {
			placeholders += "?"
		} else {
			placeholders += ",?"
		}
		args = append(args, id)
	}
	return placeholders, args
}

func getWeekdayZh(t time.Time) string {
	switch t.Weekday() {
	case time.Monday:
		return "周一"
	case time.Tuesday:
		return "周二"
	case time.Wednesday:
		return "周三"
	case time.Thursday:
		return "周四"
	case time.Friday:
		return "周五"
	case time.Saturday:
		return "周六"
	case time.Sunday:
		return "周日"
	}
	return ""
}