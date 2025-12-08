package order_record_detail

import (
	"canteen/internal/model"
	"database/sql"
	"log"
)

type OrderRecordDetailRepository interface {
	FindByDateRange(startDate, endDate string) ([]model.OrderRecordDetail, error)
}

type orderRecordDetailRepository struct {
	db *sql.DB
}

func NewOrderRecordDetailRepository(db *sql.DB) OrderRecordDetailRepository {
	return &orderRecordDetailRepository{db: db}
}

// FindByDateRange 查询指定时间范围内的点餐详情
func (r *orderRecordDetailRepository) FindByDateRange(startDate, endDate string) ([]model.OrderRecordDetail, error) {
	log.Printf("查询点餐详情，日期范围: %s 至 %s", startDate, endDate)

	// 修复表连接查询
	query := `
		SELECT 
			o.id,
			o.user_id,
			o.setmeal_id,
			s.id as dish_id,
			s.code as dish_code,
			s.description,
			o.order_date,
			o.meal_type,
			o.week_number,
			o.weekday,
			o.status
		FROM 
			order_record o
		LEFT JOIN 
			weekly_setmeal ws ON o.setmeal_id = ws.id
		LEFT JOIN 
			setmeal s ON ws.setmeal_id = s.id
		WHERE 
			o.order_date BETWEEN ? AND ?
		ORDER BY 
			o.order_date, o.user_id
	`

	log.Printf("执行查询: %s", query)
	rows, err := r.db.Query(query, startDate, endDate)
	if err != nil {
		log.Printf("查询失败: %v", err)
		return nil, err
	}
	defer rows.Close()

	var details []model.OrderRecordDetail
	for rows.Next() {
		var detail model.OrderRecordDetail
		// 使用sql.NullString处理可能为NULL的字段
		var dishId sql.NullInt64
		var dishCode sql.NullString
		var description sql.NullString
		
		err := rows.Scan(
			&detail.Id,
			&detail.UserId,
			&detail.SetmealId,
			&dishId,
			&dishCode,
			&description,
			&detail.OrderDate,
			&detail.MealType,
			&detail.WeekNumber,
			&detail.Weekday,
			&detail.Status,
		)
		if err != nil {
			log.Printf("扫描行失败: %v", err)
			return nil, err
		}

		// 处理可能为NULL的值
		if dishId.Valid {
			detail.DishId = int(dishId.Int64)
		}
		if dishCode.Valid {
			detail.DishCode = dishCode.String
		}
		if description.Valid {
			detail.Description = description.String
		}

		details = append(details, detail)
		log.Printf("找到记录: ID=%d, UserID=%d, OrderDate=%s, DishId=%d, DishCode=%s, Description=%s",
			detail.Id, detail.UserId, detail.OrderDate, detail.DishId, detail.DishCode, detail.Description)
	}

	if err = rows.Err(); err != nil {
		log.Printf("行遍历错误: %v", err)
		return nil, err
	}

	log.Printf("共找到 %d 条记录", len(details))
	return details, nil
}
