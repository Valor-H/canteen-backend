package tempDirect

import (
	"canteen/internal/model"
	"canteen/internal/service/meal"
	"canteen/internal/service/order"
	mealRepo "canteen/internal/repository/meal"
	orderRepo "canteen/internal/repository/order"
	"canteen/internal/infrastructure/cache"
	"database/sql"
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/xuri/excelize/v2"
)

var (
	db *sql.DB
	mealService meal.MealService
	orderService order.OrderService
)

func SetDB(database *sql.DB) {
	db = database
	
	// 初始化repositories
	mealRepository := mealRepo.NewMealRepository(db)
	orderRepository := orderRepo.NewOrderRepository(db)
	
	// 初始化services
	mealService = meal.NewMealService(mealRepository, cache.RedisClient())
	orderService = order.NewOrderService(orderRepository, nil, cache.RedisClient()) // userRepo设为nil，暂时不使用
}

// InsertMealSQL 插入餐食SQL
func InsertMealSQL(c *gin.Context) {
	var req model.CreateMeal
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"Status": 0, "Msg": "请求参数错误: " + err.Error()})
		return
	}
	
	if len(req.DishIds) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"Status": 0, "Msg": "菜品列表不能为空"})
		return
	}
	
	// 使用服务层处理
	dishes, err := mealService.GetDishesByIds(req.DishIds)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"Status": 0, "Msg": "查询失败: " + err.Error()})
		return
	}
	
	c.JSON(http.StatusOK, gin.H{"Status": 1, "Msg": "查询成功", "Data": dishes})
}

// ExportOrdersByDate 按日期导出订单
func ExportOrdersByDate(c *gin.Context) {
	date := c.Query("date") // 获取查询参数
	if len(date) != 8 {
		c.JSON(http.StatusBadRequest, gin.H{"status": 0, "msg": "参数错误：请传入格式为 yyyyMMdd 的日期"})
		return
	}
	log.Printf("Received date parameter: %s", date)
	
	// 使用服务层处理
	file, err := orderService.ExportOrdersByDate(date)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": 0, "msg": "导出失败: " + err.Error()})
		return
	}
	
	// 设置响应头
	c.Header("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
	c.Header("Content-Disposition", "attachment; filename=orders-"+date+".xlsx")
	
	// 写入响应
	file.Write(c.Writer)
}

// ExportOrdersByMonth 按月份导出订单
func ExportOrdersByMonth(c *gin.Context) {
	date := c.Query("month") // 获取查询参数
	if len(date) != 6 {
		c.JSON(http.StatusBadRequest, gin.H{"status": 0, "msg": "参数错误：请传入格式为 yyyyMM 的月份"})
		return
	}
	
	// 使用服务层处理
	file, err := orderService.ExportOrdersByMonth(date)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": 0, "msg": "导出失败: " + err.Error()})
		return
	}
	
	// 设置响应头
	c.Header("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
	c.Header("Content-Disposition", "attachment; filename=orders-"+date+".xlsx")
	
	// 写入响应
	file.Write(c.Writer)
}

// UploadWeekMenuHandler 上传周菜单
func UploadWeekMenuHandler(c *gin.Context) {
	// 获取上传的文件
	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": 0, "msg": "获取文件失败: " + err.Error()})
		return
	}
	
	// 打开文件
	src, err := file.Open()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": 0, "msg": "打开文件失败: " + err.Error()})
		return
	}
	defer src.Close()
	
	// 使用excelize读取文件
	xlsx, err := excelize.OpenReader(src)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": 0, "msg": "解析Excel文件失败: " + err.Error()})
		return
	}
	
	// 获取第一个工作表
	sheets := xlsx.GetSheetList()
	if len(sheets) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"status": 0, "msg": "Excel文件中没有工作表"})
		return
	}
	
	// 读取数据
	rows, err := xlsx.GetRows(sheets[0])
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": 0, "msg": "读取工作表数据失败: " + err.Error()})
		return
	}
	
	// 处理数据（简化处理）
	log.Printf("读取到 %d 行数据", len(rows))
	
	// 这里应该解析Excel数据并插入数据库
	// 为了简化，我们只是记录日志
	
	c.JSON(http.StatusOK, gin.H{"status": 1, "msg": "上传成功"})
}

// DateImport 日期导入
func DateImport(c *gin.Context) {
	var req struct {
		Date string `json:"date"`
	}
	
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": 0, "msg": "请求参数错误: " + err.Error()})
		return
	}
	
	if len(req.Date) != 8 {
		c.JSON(http.StatusBadRequest, gin.H{"status": 0, "msg": "日期格式错误，请使用 yyyyMMdd 格式"})
		return
	}
	
	// 验证日期
	_, err := strconv.Atoi(req.Date)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": 0, "msg": "日期格式错误，请使用 yyyyMMdd 格式"})
		return
	}
	
	// 使用服务层处理
	setmeals, err := mealService.GetSetmealsByWeekNumber(req.Date)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": 0, "msg": "查询失败: " + err.Error()})
		return
	}
	
	c.JSON(http.StatusOK, gin.H{"status": 1, "msg": "查询成功", "data": setmeals})
}

// DishDetail 菜品详情
func DishDetail(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": 0, "msg": "无效的ID"})
		return
	}
	
	// 使用服务层处理
	dishes, err := mealService.GetDishesByIds([]int16{int16(id)})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": 0, "msg": "查询失败: " + err.Error()})
		return
	}
	
	if len(dishes) == 0 {
		c.JSON(http.StatusNotFound, gin.H{"status": 0, "msg": "菜品不存在"})
		return
	}
	
	c.JSON(http.StatusOK, gin.H{"status": 1, "msg": "查询成功", "data": dishes[0]})
}