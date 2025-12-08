package router

import (
	"canteen/internal/controller/card"
	"canteen/internal/controller/health"
	"canteen/internal/controller/order_record_detail"
	"canteen/internal/controller/tempDirect"
	"canteen/internal/controller/uploadFile"
	"canteen/internal/controller/user"
	"canteen/internal/infrastructure/logging"
	"net/http"
	"strings"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

// BlockInvalidRequests 非法请求拦截中间件
func BlockInvalidRequests() gin.HandlerFunc {
	validPrefixes := []string{
		"/api/v1/",
		"/hxz/v1/",
		"/temp/v1/",
		"/user/v1/",
		"/order/v1/",
	}
	blockedPaths := map[string]bool{
		"/hxz/v1/test": true,
	}

	return func(c *gin.Context) {
		path := c.Request.URL.Path
		ip := c.ClientIP()
		ua := c.Request.UserAgent()
		method := c.Request.Method
		timestamp := time.Now().Format("2006-01-02 15:04:05")

		illegalLogger := logging.GetIllegalLogger()

		// 日志内容封装
		logInvalid := func(reason string) {
			illegalLogger.Printf("[非法请求] [%s] IP: %s  Method: %s  Path: %s  UA: %s  原因: %s",
				timestamp, ip, method, path, ua, reason)
		}

		if blockedPaths[path] {
			logInvalid("命中明确禁止路径")
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
				"status":  403,
				"message": "非法请求路径被拦截",
			})
			return
		}

		isValid := false
		for _, prefix := range validPrefixes {
			if strings.HasPrefix(path, prefix) {
				isValid = true
				break
			}
		}

		if !isValid {
			logInvalid("不在合法路径前缀列表中")
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
				"status":  403,
				"message": "非法请求路径被拦截",
			})
			return
		}

		c.Next()
	}
}

// RegisterRoutes 注册路由
func RegisterRoutes(router *gin.Engine) {
	commonApi := router.Group("/api")
	commonGroup := commonApi.Group("/v1")
	{
		commonGroup.GET("/health", health.HealthCheckHandler)
		commonGroup.GET("/exportDayRcord", tempDirect.ExportOrdersByDate)
		commonGroup.GET("/exportMonthRecord", tempDirect.ExportOrdersByMonth)
		commonGroup.POST("/uploadWeekMenu", tempDirect.UploadWeekMenuHandler)
		commonGroup.POST("/dateImport", tempDirect.DateImport)
		commonGroup.GET("/DishDetail/:id", tempDirect.DishDetail)
	}

	userApi := router.Group("/user")
	userGroup := userApi.Group("/v1")
	{
		userGroup.GET("/getUser/:user_id", user.GetUserHandler)
		userGroup.GET("/getUserByNickName", user.GetUserByNickNameHandler)
	}

	orderApi := router.Group("/order")
	orderGroup := orderApi.Group("/v1")
	{
		orderGroup.GET("/getCMealSelectionStats", order_record_detail.GetCMealSelectionStatsHandler)
		orderGroup.GET("/getAllMealSelectionStats", order_record_detail.GetAllMealSelectionStatsHandler)
		orderGroup.GET("/getBasicDishStats", order_record_detail.GetBasicDishStatsHandler)
		orderGroup.GET("/getDishAppearanceStats", order_record_detail.GetDishAppearanceStatsHandler)
		orderGroup.GET("/getUserDishOrderStats", order_record_detail.GetUserDishOrderStatsHandler)
		orderGroup.GET("/getDishStatsComparison", order_record_detail.GetDishStatsComparisonHandler)
	}

	cardApi := router.Group("/hxz")
	cardGroup := cardApi.Group("/v1")
	{
		cardGroup.POST("/ConsumTransactions", card.ConsumTransactionHandler)
		cardGroup.POST("/ServerTime", card.ServerTimeHandler)
		cardGroup.POST("/OffLines", card.OffLineHandler)
	}

	tempApi := router.Group("/temp")
	tempGroup := tempApi.Group("/v1")
	{
		tempGroup.POST("/InsertMealSQL", tempDirect.InsertMealSQL)
		tempGroup.POST("/upload", uploadFile.UploadFileHandler)
	}
}

// Config 配置路由
func Config() *gin.Engine {
	router := gin.New()

	// 使用日志中间件
	router.Use(logging.LoggerMiddleware())

	// 恢复 panic 的中间件
	router.Use(gin.Recovery())

	// CORS 配置
	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},                                       // 允许所有域名跨域访问，如果需要可以改为特定域名
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"}, // 允许的请求方法
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"}, // 允许的请求头
		ExposeHeaders:    []string{"Content-Length"},                          // 允许的响应头
		AllowCredentials: true,                                                // 是否允许携带凭证（如 Cookie）
		MaxAge:           12 * time.Hour,                                      // 设置缓存时间
	}))

	// 拦截非法请求
	router.Use(BlockInvalidRequests())

	// 注册路由
	RegisterRoutes(router)

	return router
}
