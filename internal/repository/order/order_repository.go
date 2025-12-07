package order

import (
	"canteen/internal/model"
	"database/sql"
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/xuri/excelize/v2"
)

type OrderRepository interface {
	FindByWeekNumber(weekNumber string) ([]model.OrderRecord, error)
	UpdateStatusToExpired(weekNumber string) error
	FindExpiredOrdersByWeekNumber(weekNumber string) ([]int, error)
	CreateOrder(order *model.OrderRecord) error
	FindOrdersForExport(weekNumber string) ([]ExportOrderRecord, error)
	ExportToExcel(date string) (*excelize.File, error)
}

type orderRepository struct {
	db *sql.DB
}

type ExportOrderRecord struct {
	WorkNo    string
	Name      string
	Dept      string
	MealType  string
	Date      string
	Weekday   string
	Status    string
}

func NewOrderRepository(db *sql.DB) OrderRepository {
	return &orderRepository{db: db}
}

func (r *orderRepository) FindByWeekNumber(weekNumber string) ([]model.OrderRecord, error) {
	query := `SELECT id, user_id, status, meal_id FROM order_record WHERE week_number = ?`
	
	rows, err := r.db.Query(query, weekNumber)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	
	var orders []model.OrderRecord
	for rows.Next() {
		var order model.OrderRecord
		if err := rows.Scan(&order.Id, &order.UserId, &order.Status, &order.MealId); err != nil {
			continue
		}
		orders = append(orders, order)
	}
	
	return orders, rows.Err()
}

func (r *orderRepository) UpdateStatusToExpired(weekNumber string) error {
	_, err := r.db.Exec(
		"UPDATE order_record SET status = '已过期', update_time = NOW() WHERE week_number = ? AND status = '已报餐'",
		weekNumber,
	)
	return err
}

func (r *orderRepository) FindExpiredOrdersByWeekNumber(weekNumber string) ([]int, error) {
	rows, err := r.db.Query(
		"SELECT DISTINCT user_id FROM order_record WHERE week_number = ? AND status = '已报餐'",
		weekNumber,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	
	var userIds []int
	for rows.Next() {
		var userId int
		if err := rows.Scan(&userId); err != nil {
			continue
		}
		userIds = append(userIds, userId)
	}
	
	return userIds, rows.Err()
}

func (r *orderRepository) CreateOrder(order *model.OrderRecord) error {
	query := `INSERT INTO order_record (user_id, meal_id, status, create_time, meal_type, week_number, order_date, weekday) 
	          VALUES (?, ?, ?, ?, ?, ?, ?, ?)`
	_, err := r.db.Exec(query, order.UserId, order.MealId, order.Status, time.Now(), 
		order.MealType, order.WeekNumber, order.OrderDate, order.Weekday)
	return err
}

func (r *orderRepository) FindOrdersForExport(weekNumber string) ([]ExportOrderRecord, error) {
	query := `
		SELECT 
			s.user_name AS 工号,
			s.nick_name AS 姓名,
			sd.dept_name AS 部门,
			ord.meal_type AS 餐别,
			DATE_FORMAT(STR_TO_DATE(CAST(ord.week_number AS CHAR), '%Y%m%d'), '%Y/%c/%e') AS 日期,
			ord.weekday AS 星期,
			ord.status AS 状态
		FROM order_record ord
		LEFT JOIN sys_user s ON ord.user_id = s.user_id
		LEFT JOIN sys_dept sd ON s.dept_id = sd.dept_id
		WHERE ord.week_number = ?
		ORDER BY sd.dept_name, ord.status
	`
	
	rows, err := r.db.Query(query, weekNumber)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	
	var orders []ExportOrderRecord
	for rows.Next() {
		var order ExportOrderRecord
		if err := rows.Scan(&order.WorkNo, &order.Name, &order.Dept, &order.MealType, &order.Date, &order.Weekday, &order.Status); err != nil {
			log.Printf("Row scan failed: %v", err)
			continue
		}
		orders = append(orders, order)
	}
	
	return orders, rows.Err()
}

// ExportToExcel 将订单数据导出到Excel
func (r *orderRepository) ExportToExcel(date string) (*excelize.File, error) {
	// 验证日期格式
	if len(date) != 8 {
		return nil, fmt.Errorf("参数错误：请传入格式为 yyyyMMdd 的日期")
	}
	
	_, err := strconv.Atoi(date)
	if err != nil {
		return nil, fmt.Errorf("日期格式错误")
	}
	
	// 查询数据
	orders, err := r.FindOrdersForExport(date)
	if err != nil {
		return nil, fmt.Errorf("查询失败: %v", err)
	}
	
	// 创建Excel文件
	f := excelize.NewFile()
	sheet := "报餐记录"
	f.SetSheetName("Sheet1", sheet)
	
	// 设置表头
	headers := []string{"工号", "姓名", "部门", "餐别", "日期", "星期", "状态"}
	for i, header := range headers {
		cell := fmt.Sprintf("%s1", string(rune('A'+i)))
		f.SetCellValue(sheet, cell, header)
	}
	
	// 填充数据
	for i, order := range orders {
		rowNum := i + 2 // 从第2行开始
		f.SetCellValue(sheet, fmt.Sprintf("A%d", rowNum), order.WorkNo)
		f.SetCellValue(sheet, fmt.Sprintf("B%d", rowNum), order.Name)
		f.SetCellValue(sheet, fmt.Sprintf("C%d", rowNum), order.Dept)
		f.SetCellValue(sheet, fmt.Sprintf("D%d", rowNum), order.MealType)
		f.SetCellValue(sheet, fmt.Sprintf("E%d", rowNum), order.Date)
		f.SetCellValue(sheet, fmt.Sprintf("F%d", rowNum), order.Weekday)
		f.SetCellValue(sheet, fmt.Sprintf("G%d", rowNum), order.Status)
	}
	
	return f, nil
}