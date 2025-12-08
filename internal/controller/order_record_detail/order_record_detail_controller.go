package order_record_detail

import (
	ordRepo "canteen/internal/repository/order_record_detail"
	ordService "canteen/internal/service/order_record_detail"
	"database/sql"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

var (
	db            *sql.DB
	detailService ordService.OrderRecordDetailService
)

func SetDB(database *sql.DB) {
	db = database

	// 初始化repository
	repository := ordRepo.NewOrderRecordDetailRepository(db)

	// 初始化service
	detailService = ordService.NewOrderRecordDetailService(repository)
}

// AllMealSelectionStats 所有菜品选择统计
type AllMealSelectionStats struct {
	StartDate string         `json:"startDate"` // 查询开始日期
	EndDate   string         `json:"endDate"`   // 查询结束日期
	Items     map[string]int `json:"items"`     // 菜品名称及其出现次数
}

// GetAllMealSelectionStatsHandler 获取时间范围内所有菜品选择统计处理器
func GetAllMealSelectionStatsHandler(c *gin.Context) {
	// 从查询参数中获取日期范围
	startDate := c.Query("start_date")
	endDate := c.Query("end_date")

	if startDate == "" || endDate == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  400,
			"message": "开始日期(start_date)和结束日期(end_date)不能为空",
		})
		return
	}

	// 验证日期格式
	if _, err := time.Parse("2006-01-02", startDate); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  400,
			"message": "开始日期格式错误，请使用YYYY-MM-DD格式",
		})
		return
	}

	if _, err := time.Parse("2006-01-02", endDate); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  400,
			"message": "结束日期格式错误，请使用YYYY-MM-DD格式",
		})
		return
	}

	// 转换日期格式，去掉分隔符，以便匹配setmeal表中的code格式（如M20251201-X-X）
	// 将YYYY-MM-DD转换为YYYYMMDD
	startDateForCode := strings.ReplaceAll(startDate, "-", "")
	endDateForCode := strings.ReplaceAll(endDate, "-", "")

	log.Printf("查询菜品出现次数统计，日期范围: %s 到 %s (转换为code格式: %s 到 %s)",
		startDate, endDate, startDateForCode, endDateForCode)

	// 直接查询数据库获取所有菜品的数据
	query := `
		SELECT 
			s.description, COUNT(*) as count
		FROM 
			order_record o
		LEFT JOIN 
			weekly_setmeal ws ON o.setmeal_id = ws.id
		LEFT JOIN 
			setmeal s ON ws.setmeal_id = s.id
		WHERE 
			o.order_date BETWEEN ? AND ?
			AND s.description IS NOT NULL
		GROUP BY 
			s.description
		ORDER BY 
			count DESC
	`

	rows, err := db.Query(query, startDateForCode, endDateForCode)
	if err != nil {
		log.Printf("查询所有菜品统计失败: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  500,
			"message": "获取所有菜品统计失败: " + err.Error(),
		})
		return
	}
	defer rows.Close()

	items := make(map[string]int)
	totalCount := 0

	for rows.Next() {
		var description string
		var count int
		err := rows.Scan(&description, &count)
		if err != nil {
			log.Printf("扫描行失败: %v", err)
			continue
		}

		items[description] = count
		totalCount += count
		log.Printf("菜品统计: %s = %d次", description, count)
	}

	if err = rows.Err(); err != nil {
		log.Printf("行遍历错误: %v", err)
	}

	stats := AllMealSelectionStats{
		StartDate: startDate,
		EndDate:   endDate,
		Items:     items,
	}

	log.Printf("所有菜品统计完成，共%d种菜品，总计%d次", len(items), totalCount)

	// 返回成功响应
	c.JSON(http.StatusOK, gin.H{
		"status":  200,
		"message": "请求成功",
		"data":    stats,
	})
}

// CMealSelectionStats C套餐选择统计
type CMealSelectionStats struct {
	StartDate string         `json:"startDate"` // 查询开始日期
	EndDate   string         `json:"endDate"`   // 查询结束日期
	Items     map[string]int `json:"items"`     // 套餐名称及其出现次数
}

// GetCMealSelectionStatsHandler 获取时间范围内的C套餐选择统计处理器
func GetCMealSelectionStatsHandler(c *gin.Context) {
	// 从查询参数中获取日期范围
	startDate := c.Query("start_date")
	endDate := c.Query("end_date")

	if startDate == "" || endDate == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  400,
			"message": "开始日期(start_date)和结束日期(end_date)不能为空",
		})
		return
	}

	// 验证日期格式
	if _, err := time.Parse("2006-01-02", startDate); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  400,
			"message": "开始日期格式错误，请使用YYYY-MM-DD格式",
		})
		return
	}

	if _, err := time.Parse("2006-01-02", endDate); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  400,
			"message": "结束日期格式错误，请使用YYYY-MM-DD格式",
		})
		return
	}

	// 转换日期格式，去掉分隔符，以便匹配setmeal表中的code格式（如M20251201-X-X）
	// 将YYYY-MM-DD转换为YYYYMMDD
	startDateForCode := strings.ReplaceAll(startDate, "-", "")
	endDateForCode := strings.ReplaceAll(endDate, "-", "")

	log.Printf("查询菜品出现次数统计，日期范围: %s 到 %s (转换为code格式: %s 到 %s)",
		startDate, endDate, startDateForCode, endDateForCode)

	// 直接查询数据库获取C套餐的数据
	query := `
		SELECT 
			s.description, COUNT(*) as count
		FROM 
			order_record o
		LEFT JOIN 
			weekly_setmeal ws ON o.setmeal_id = ws.id
		LEFT JOIN 
			setmeal s ON ws.setmeal_id = s.id
		WHERE 
			o.order_date BETWEEN ? AND ?
			AND s.code LIKE '%-C'
			AND s.description IS NOT NULL
		GROUP BY 
			s.description
		ORDER BY 
			count DESC
	`

	rows, err := db.Query(query, startDateForCode, endDateForCode)
	if err != nil {
		log.Printf("查询C套餐统计失败: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  500,
			"message": "获取C套餐统计失败: " + err.Error(),
		})
		return
	}
	defer rows.Close()

	items := make(map[string]int)
	totalCount := 0

	for rows.Next() {
		var description string
		var count int
		err := rows.Scan(&description, &count)
		if err != nil {
			log.Printf("扫描行失败: %v", err)
			continue
		}

		items[description] = count
		totalCount += count
		log.Printf("C套餐统计: %s = %d次", description, count)
	}

	if err = rows.Err(); err != nil {
		log.Printf("行遍历错误: %v", err)
	}

	stats := CMealSelectionStats{
		StartDate: startDate,
		EndDate:   endDate,
		Items:     items,
	}

	log.Printf("C套餐统计完成，共%d种套餐，总计%d次", len(items), totalCount)

	// 返回成功响应
	c.JSON(http.StatusOK, gin.H{
		"status":  200,
		"message": "请求成功",
		"data":    stats,
	})
}

// BasicDishStats 基本菜品选择统计
type BasicDishStats struct {
	StartDate string         `json:"startDate"` // 查询开始日期
	EndDate   string         `json:"endDate"`   // 查询结束日期
	Items     map[string]int `json:"items"`     // 基本菜品名称及其出现次数
}

// GetBasicDishStatsHandler 获取时间范围内基本菜品选择统计处理器
func GetBasicDishStatsHandler(c *gin.Context) {
	// 从查询参数中获取日期范围
	startDate := c.Query("start_date")
	endDate := c.Query("end_date")

	if startDate == "" || endDate == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  400,
			"message": "开始日期(start_date)和结束日期(end_date)不能为空",
		})
		return
	}

	// 验证日期格式
	if _, err := time.Parse("2006-01-02", startDate); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  400,
			"message": "开始日期格式错误，请使用YYYY-MM-DD格式",
		})
		return
	}

	if _, err := time.Parse("2006-01-02", endDate); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  400,
			"message": "结束日期格式错误，请使用YYYY-MM-DD格式",
		})
		return
	}

	// 转换日期格式，去掉分隔符，以便匹配setmeal表中的code格式（如M20251201-X-X）
	// 将YYYY-MM-DD转换为YYYYMMDD
	startDateForCode := strings.ReplaceAll(startDate, "-", "")
	endDateForCode := strings.ReplaceAll(endDate, "-", "")

	log.Printf("查询菜品出现次数统计，日期范围: %s 到 %s (转换为code格式: %s 到 %s)",
		startDate, endDate, startDateForCode, endDateForCode)

	// 直接查询数据库获取所有菜品的数据
	query := `
		SELECT 
			s.description, COUNT(*) as count
		FROM 
			order_record o
		LEFT JOIN 
			weekly_setmeal ws ON o.setmeal_id = ws.id
		LEFT JOIN 
			setmeal s ON ws.setmeal_id = s.id
		WHERE 
			o.order_date BETWEEN ? AND ?
			AND s.description IS NOT NULL
		GROUP BY 
			s.description
		ORDER BY 
			count DESC
	`

	rows, err := db.Query(query, startDateForCode, endDateForCode)
	if err != nil {
		log.Printf("查询基本菜品统计失败: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  500,
			"message": "获取基本菜品统计失败: " + err.Error(),
		})
		return
	}
	defer rows.Close()

	items := make(map[string]int)
	totalCount := 0

	for rows.Next() {
		var description string
		var count int
		err := rows.Scan(&description, &count)
		if err != nil {
			log.Printf("扫描行失败: %v", err)
			continue
		}

		// 分割description，统计基本菜品
		// 例如：将"肉沫茄子+清炒杭白菜+鲜花椒鸡块"分割成多个基本菜品
		basicDishes := strings.Split(description, "+")
		for _, dish := range basicDishes {
			dish = strings.TrimSpace(dish)
			if dish != "" {
				items[dish] += count
				totalCount += count
				log.Printf("基本菜品统计: %s = %d次", dish, items[dish])
			}
		}
	}

	if err = rows.Err(); err != nil {
		log.Printf("行遍历错误: %v", err)
	}

	stats := BasicDishStats{
		StartDate: startDate,
		EndDate:   endDate,
		Items:     items,
	}

	log.Printf("基本菜品统计完成，共%d种菜品，总计%d次", len(items), totalCount)

	// 返回成功响应
	c.JSON(http.StatusOK, gin.H{
		"status":  200,
		"message": "请求成功",
		"data":    stats,
	})
}

// DishAppearanceStats 菜品出现次数统计（精确统计）
type DishAppearanceStats struct {
	StartDate string         `json:"startDate"` // 查询开始日期
	EndDate   string         `json:"endDate"`   // 查询结束日期
	Items     map[string]int `json:"items"`     // 菜品名称及其出现次数
}

// GetDishAppearanceStatsHandler 获取时间范围内菜品出现次数统计处理器
func GetDishAppearanceStatsHandler(c *gin.Context) {
	// 从查询参数中获取日期范围
	startDate := c.Query("start_date")
	endDate := c.Query("end_date")

	if startDate == "" || endDate == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  400,
			"message": "开始日期(start_date)和结束日期(end_date)不能为空",
		})
		return
	}

	// 验证日期格式
	if _, err := time.Parse("2006-01-02", startDate); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  400,
			"message": "开始日期格式错误，请使用YYYY-MM-DD格式",
		})
		return
	}

	if _, err := time.Parse("2006-01-02", endDate); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  400,
			"message": "结束日期格式错误，请使用YYYY-MM-DD格式",
		})
		return
	}

	// 转换日期格式，去掉分隔符，以便匹配setmeal表中的code格式（如M20251201-X-X）
	// 将YYYY-MM-DD转换为YYYYMMDD
	startDateForCode := strings.ReplaceAll(startDate, "-", "")
	endDateForCode := strings.ReplaceAll(endDate, "-", "")

	log.Printf("查询菜品出现次数统计，日期范围: %s 到 %s (转换为code格式: %s 到 %s)",
		startDate, endDate, startDateForCode, endDateForCode)

	// 按照需求：先从setmeal表中根据code获取时间范围内的菜单，然后统计这些菜单中菜品的出现次数
	query := `
		SELECT 
			d.name, COUNT(*) as appearance_count
		FROM 
			setmeal s
		JOIN 
			setmeal_dish sd ON s.id = sd.setmeal_id
		JOIN 
			dish d ON sd.dish_id = d.id
		WHERE 
			SUBSTRING(s.code, 2, 8) BETWEEN ? AND ?
		GROUP BY 
			d.id, d.name
		ORDER BY 
			appearance_count DESC
	`

	rows, err := db.Query(query, startDateForCode, endDateForCode)
	if err != nil {
		log.Printf("查询菜品出现次数统计失败: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  500,
			"message": "获取菜品出现次数统计失败: " + err.Error(),
		})
		return
	}
	defer rows.Close()

	items := make(map[string]int)
	totalCount := 0

	for rows.Next() {
		var dishName string
		var count int
		err := rows.Scan(&dishName, &count)
		if err != nil {
			log.Printf("扫描行失败: %v", err)
			continue
		}

		items[dishName] = count
		totalCount += count
		log.Printf("菜品出现次数统计: %s = %d次", dishName, count)
	}

	if err = rows.Err(); err != nil {
		log.Printf("行遍历错误: %v", err)
	}

	stats := DishAppearanceStats{
		StartDate: startDate,
		EndDate:   endDate,
		Items:     items,
	}

	log.Printf("菜品出现次数统计完成，共%d种菜品，总计%d次", len(items), totalCount)

	// 返回成功响应
	c.JSON(http.StatusOK, gin.H{
		"status":  200,
		"message": "请求成功",
		"data":    stats,
	})
}

// UserDishOrderStats 用户点餐菜品统计
type UserDishOrderStats struct {
	StartDate string         `json:"startDate"` // 查询开始日期
	EndDate   string         `json:"endDate"`   // 查询结束日期
	Items     map[string]int `json:"items"`     // 菜品名称及其被点餐次数
}

// GetUserDishOrderStatsHandler 获取时间范围内用户点餐菜品统计处理器
func GetUserDishOrderStatsHandler(c *gin.Context) {
	// 从查询参数中获取日期范围
	startDate := c.Query("start_date")
	endDate := c.Query("end_date")

	if startDate == "" || endDate == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  400,
			"message": "开始日期(start_date)和结束日期(end_date)不能为空",
		})
		return
	}

	// 验证日期格式
	if _, err := time.Parse("2006-01-02", startDate); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  400,
			"message": "开始日期格式错误，请使用YYYY-MM-DD格式",
		})
		return
	}

	if _, err := time.Parse("2006-01-02", endDate); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  400,
			"message": "结束日期格式错误，请使用YYYY-MM-DD格式",
		})
		return
	}

	// 转换日期格式，去掉分隔符，转换为与week_number字段匹配的格式（如20251224）
	startWeekNumber := strings.ReplaceAll(startDate, "-", "")
	endWeekNumber := strings.ReplaceAll(endDate, "-", "")

	log.Printf("查询用户点餐菜品统计，日期范围: %s 到 %s (转换为week_number格式: %s 到 %s)",
		startDate, endDate, startWeekNumber, endWeekNumber)

	// 根据week_number字段筛选订单，然后通过关联表获取菜品统计
	query := `
		SELECT 
			d.name, COUNT(*) as order_count
		FROM 
			order_record o
		JOIN 
			weekly_setmeal ws ON o.setmeal_id = ws.id
		JOIN 
			setmeal s ON ws.setmeal_id = s.id
		JOIN 
			setmeal_dish sd ON s.id = sd.setmeal_id
		JOIN 
			dish d ON sd.dish_id = d.id
		WHERE 
			o.week_number BETWEEN ? AND ?
		GROUP BY 
			d.id, d.name
		ORDER BY 
			order_count DESC
	`

	rows, err := db.Query(query, startWeekNumber, endWeekNumber)
	if err != nil {
		log.Printf("查询用户点餐菜品统计失败: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  500,
			"message": "获取用户点餐菜品统计失败: " + err.Error(),
		})
		return
	}
	defer rows.Close()

	items := make(map[string]int)
	totalCount := 0

	for rows.Next() {
		var dishName string
		var count int
		err := rows.Scan(&dishName, &count)
		if err != nil {
			log.Printf("扫描行失败: %v", err)
			continue
		}

		items[dishName] = count
		totalCount += count
		log.Printf("用户点餐菜品统计: %s = %d次", dishName, count)
	}

	if err = rows.Err(); err != nil {
		log.Printf("行遍历错误: %v", err)
	}

	stats := UserDishOrderStats{
		StartDate: startDate,
		EndDate:   endDate,
		Items:     items,
	}

	log.Printf("用户点餐菜品统计完成，共%d种菜品，总计%d次", len(items), totalCount)

	// 返回成功响应
	c.JSON(http.StatusOK, gin.H{
		"status":  200,
		"message": "请求成功",
		"data":    stats,
	})
}

// DishStatsComparisonItem 菜品统计对比项
type DishStatsComparisonItem struct {
	DishId          int     `json:"dishId"`          // 菜品ID
	DishName        string  `json:"dishName"`        // 菜品名称
	CategoryId      int     `json:"categoryId"`       // 菜品类型ID
	AppearanceCount int     `json:"appearanceCount"` // 菜品在菜单中的出现次数
	OrderCount      int     `json:"orderCount"`      // 用户点餐次数
	Ratio           float64 `json:"ratio"`           // 点餐次数/出现次数的比值
}

// DishStatsComparison 菜品统计对比
type DishStatsComparison struct {
	StartDate string                    `json:"startDate"` // 查询开始日期
	EndDate   string                    `json:"endDate"`   // 查询结束日期
	Items     []DishStatsComparisonItem `json:"items"`     // 菜品统计对比列表
}

// GetDishStatsComparisonHandler 获取菜品统计对比处理器
func GetDishStatsComparisonHandler(c *gin.Context) {
	// 从查询参数中获取日期范围
	startDate := c.Query("start_date")
	endDate := c.Query("end_date")

	if startDate == "" || endDate == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  400,
			"message": "开始日期(start_date)和结束日期(end_date)不能为空",
		})
		return
	}

	// 验证日期格式
	if _, err := time.Parse("2006-01-02", startDate); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  400,
			"message": "开始日期格式错误，请使用YYYY-MM-DD格式",
		})
		return
	}

	if _, err := time.Parse("2006-01-02", endDate); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  400,
			"message": "结束日期格式错误，请使用YYYY-MM-DD格式",
		})
		return
	}

	// 转换日期格式，去掉分隔符
	startDateForCode := strings.ReplaceAll(startDate, "-", "")
	endDateForCode := strings.ReplaceAll(endDate, "-", "")

	log.Printf("查询菜品统计对比，日期范围: %s 到 %s (转换为code格式: %s 到 %s)",
		startDate, endDate, startDateForCode, endDateForCode)

	// 1. 查询菜品在菜单中的出现次数（对应getDishAppearanceStats接口）
	appearanceQuery := `
		SELECT 
			d.id, d.name, d.category_id, COUNT(*) as appearance_count
		FROM 
			setmeal s
		JOIN 
			setmeal_dish sd ON s.id = sd.setmeal_id
		JOIN 
			dish d ON sd.dish_id = d.id
		WHERE 
			SUBSTRING(s.code, 2, 8) BETWEEN ? AND ?
		GROUP BY 
			d.id, d.name, d.category_id
	`

	appearanceRows, err := db.Query(appearanceQuery, startDateForCode, endDateForCode)
	if err != nil {
		log.Printf("查询菜品出现次数统计失败: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  500,
			"message": "获取菜品出现次数统计失败: " + err.Error(),
		})
		return
	}
	defer appearanceRows.Close()

	// 存储菜品在菜单中的出现次数（使用dishId作为键）
	appearanceCountMap := make(map[int]struct {
		name       string
		categoryId int
		count      int
	})
	for appearanceRows.Next() {
		var dishId, categoryId int
		var dishName string
		var count int
		err := appearanceRows.Scan(&dishId, &dishName, &categoryId, &count)
		if err != nil {
			log.Printf("扫描菜品出现次数行失败: %v", err)
			continue
		}
		appearanceCountMap[dishId] = struct {
			name       string
			categoryId int
			count      int
		}{name: dishName, categoryId: categoryId, count: count}
	}

	if err = appearanceRows.Err(); err != nil {
		log.Printf("菜品出现次数查询遍历错误: %v", err)
	}

	// 2. 查询用户点餐次数（对应getUserDishOrderStats接口）
	orderQuery := `
		SELECT 
			d.id, d.name, d.category_id, COUNT(*) as order_count
		FROM 
			order_record o
		JOIN 
			weekly_setmeal ws ON o.setmeal_id = ws.id
		JOIN 
			setmeal s ON ws.setmeal_id = s.id
		JOIN 
			setmeal_dish sd ON s.id = sd.setmeal_id
		JOIN 
			dish d ON sd.dish_id = d.id
		WHERE 
			o.week_number BETWEEN ? AND ?
		GROUP BY 
			d.id, d.name, d.category_id
	`

	orderRows, err := db.Query(orderQuery, startDateForCode, endDateForCode)
	if err != nil {
		log.Printf("查询用户点餐菜品统计失败: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  500,
			"message": "获取用户点餐菜品统计失败: " + err.Error(),
		})
		return
	}
	defer orderRows.Close()

	// 存储用户点餐次数（使用dishId作为键）
	orderCountMap := make(map[int]int)
	for orderRows.Next() {
		var dishId, categoryId int
		var dishName string
		var count int
		err := orderRows.Scan(&dishId, &dishName, &categoryId, &count)
		if err != nil {
			log.Printf("扫描用户点餐行失败: %v", err)
			continue
		}
		orderCountMap[dishId] = count
	}

	if err = orderRows.Err(); err != nil {
		log.Printf("用户点餐查询遍历错误: %v", err)
	}

// 3. 合并两个查询结果，计算比值
	var comparisonItems []DishStatsComparisonItem
	totalAppearanceDishes := len(appearanceCountMap)
	totalOrderDishes := len(orderCountMap)

	// 遍历所有在菜单中出现的菜品
	for dishId, appearanceData := range appearanceCountMap {
		appearanceCount := appearanceData.count
		dishName := appearanceData.name
		categoryId := appearanceData.categoryId
		orderCount := orderCountMap[dishId] // 如果用户没有点这道菜，则为0
		ratio := 0.0
		
		// 计算比值，避免除以0
		if appearanceCount > 0 {
			ratio = float64(orderCount) / float64(appearanceCount)
		}
		
		comparisonItems = append(comparisonItems, DishStatsComparisonItem{
			DishId:          dishId,
			DishName:        dishName,
			CategoryId:      categoryId,
			AppearanceCount: appearanceCount,
			OrderCount:      orderCount,
			Ratio:           ratio,
		})
		
		log.Printf("菜品对比统计: ID:%d, %s, 类别ID:%d, 菜单出现%d次, 用户点餐%d次, 比值=%.2f", 
			dishId, dishName, categoryId, appearanceCount, orderCount, ratio)
	}

	// 按比值降序排序（热门菜品在前）
	for i := 0; i < len(comparisonItems)-1; i++ {
		for j := i + 1; j < len(comparisonItems); j++ {
			if comparisonItems[i].Ratio < comparisonItems[j].Ratio {
				comparisonItems[i], comparisonItems[j] = comparisonItems[j], comparisonItems[i]
			}
		}
	}

	comparison := DishStatsComparison{
		StartDate: startDate,
		EndDate:   endDate,
		Items:     comparisonItems,
	}

	log.Printf("菜品统计对比完成，菜单中%d种菜品，用户点餐%d种菜品", totalAppearanceDishes, totalOrderDishes)

	// 返回成功响应
	c.JSON(http.StatusOK, gin.H{
		"status":  200,
		"message": "请求成功",
		"data":    comparison,
	})
}
