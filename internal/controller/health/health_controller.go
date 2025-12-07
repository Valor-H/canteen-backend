package health

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

// HealthCheckHandler 健康检查接口
func HealthCheckHandler(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status":  "ok",
		"time":    time.Now(),
		"message": "service is running",
	})
}