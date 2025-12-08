package user

import (
	userRepo "canteen/internal/repository/user"
	"canteen/internal/service/user"
	"database/sql"
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

var (
	db          *sql.DB
	userService user.UserService
)

func SetDB(database *sql.DB) {
	db = database

	// 初始化repository
	userRepository := userRepo.NewUserRepository(db)

	// 初始化service
	userService = user.NewUserService(userRepository)
}

// GetUserHandler 获取用户信息处理器
func GetUserHandler(c *gin.Context) {
	// 从URL参数中获取用户ID
	userIdStr := c.Param("user_id")
	if userIdStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  400,
			"message": "用户ID不能为空",
		})
		return
	}

	// 将字符串转换为整数
	userId, err := strconv.Atoi(userIdStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  400,
			"message": "用户ID格式错误",
		})
		return
	}

	// 通过用户服务获取用户信息
	user, err := userService.FindById(userId)
	if err != nil {
		log.Printf("查询用户失败: %v", err)
		c.JSON(http.StatusNotFound, gin.H{
			"status":  404,
			"message": "用户不存在",
			"error":   err.Error(),
		})
		return
	}

	// 返回成功响应
	c.JSON(http.StatusOK, gin.H{
		"status":  200,
		"message": "请求成功",
		"data":    user,
	})
}

// GetUserByNickNameHandler 根据昵称获取用户信息处理器
func GetUserByNickNameHandler(c *gin.Context) {
	// 从查询参数中获取用户昵称
	nickName := c.Query("nick_name")
	if nickName == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  400,
			"message": "用户昵称(nick_name)不能为空",
		})
		return
	}

	// 通过用户服务获取用户信息
	user, err := userService.FindByNickName(nickName)
	if err != nil {
		log.Printf("查询用户失败: %v", err)
		c.JSON(http.StatusNotFound, gin.H{
			"status":  404,
			"message": "用户不存在",
			"error":   err.Error(),
		})
		return
	}

	// 返回成功响应
	c.JSON(http.StatusOK, gin.H{
		"status":  200,
		"message": "请求成功",
		"data":    user,
	})
}
