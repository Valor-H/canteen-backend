package card

import (
	"canteen/internal/model"
	"canteen/internal/service/card"
	"canteen/internal/service/user"
	userRepo "canteen/internal/repository/user"
	orderRepo "canteen/internal/repository/order"
	cardRepo "canteen/internal/repository/card"
	"canteen/internal/infrastructure/cache"
	"database/sql"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

var (
	db *sql.DB
	cardService card.CardService
	userService user.UserService
)

func SetDB(database *sql.DB) {
	db = database
	
	// 初始化repositories
	userRepository := userRepo.NewUserRepository(db)
	orderRepository := orderRepo.NewOrderRepository(db)
	cardRepository := cardRepo.NewCardRepository(db)
	
	// 初始化services
	userService = user.NewUserService(userRepository)
	cardService = card.NewCardService(userRepository, orderRepository, cardRepository, cache.RedisClient())
}

// ConsumTransactionHandler 核销接口
func ConsumTransactionHandler(c *gin.Context) {
	deviceID := c.GetHeader("Device-ID")
	
	var req model.ConsumTransaction
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Printf("TAG: 参数绑定失败: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"Status": 0, "Msg": "请求参数错误: " + err.Error()})
		return
	}
	
	response, err := cardService.ProcessConsumTransaction(req, deviceID)
	if err != nil {
		log.Printf("TAG: 处理消费交易失败: %v", err)
		c.JSON(http.StatusOK, gin.H{"Status": 0, "Msg": err.Error()})
		return
	}
	
	// 返回完整响应
	c.JSON(http.StatusOK, gin.H{
		"Status":         response.Status,
		"Msg":            response.Message,
		"Name":           response.Name,
		"CardNo":         response.CardNo,
		"Money":          response.Money,
		"Subsidy":        response.Subsidy,
		"Times":          response.Times,
		"Integral":       response.Integral,
		"InTime":         response.InTime,
		"OutTime":        response.OutTime,
		"CumulativeTime": response.Cumulative,
		"Amount":         response.Amount,
		"VoiceID":        response.VoiceID,
		"Text":           response.Text,
	})
}

// ServerTimeHandler 服务器时间接口
func ServerTimeHandler(c *gin.Context) {
	deviceID := c.GetHeader("Device-ID")
	// 检查设备ID（保留原有逻辑）
	deviceToRemark := map[string]string{
		"0180800116": "A",
		"0127448632": "B",
		"0158577664": "C",
	}
	if deviceToRemark[deviceID] == "" {
		log.Printf("deviceID=%s", deviceID)
	}
	
	serverTime := cardService.GetServerTime()
	number := string(rune('0' + ((int(serverTime.Weekday()) + 6) % 7)))
	// 服务器时间格式：yyyyMMddHHmmssd
	formattedTime := serverTime.Format("20060102150405") + number
	
	c.JSON(http.StatusOK, gin.H{
		"Status":     1,
		"Msg":        "",
		"Time":       formattedTime,
		"faceinit":   0,
		"faceAction": 0,
		"menuno":     0,
	})
}

// OffLineHandler 离线处理接口
func OffLineHandler(c *gin.Context) {
	var req model.OffLineRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"Status": 0, "Msg": "请求参数错误: " + err.Error()})
		return
	}
	
	err := cardService.ProcessOffLineRequest(req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"Status": 0, "Msg": "处理离线请求失败: " + err.Error()})
		return
	}
	
	c.JSON(http.StatusOK, gin.H{"Status": 1, "Msg": "处理成功"})
}