package uploadFile

import (
	"fmt"
	"log"
	"net/http"
	"path/filepath"
	"time"

	"github.com/gin-gonic/gin"
)

// UploadFileHandler 文件上传处理
func UploadFileHandler(c *gin.Context) {
	// 获取上传的文件
	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": 0, "msg": "获取文件失败: " + err.Error()})
		return
	}
	
	// 生成唯一文件名
	ext := filepath.Ext(file.Filename)
	filename := fmt.Sprintf("%d%s", time.Now().UnixNano(), ext)
	
	// 保存文件
	dst := filepath.Join("./uploads", filename)
	if err := c.SaveUploadedFile(file, dst); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": 0, "msg": "保存文件失败: " + err.Error()})
		return
	}
	
	log.Printf("文件上传成功: %s", dst)
	
	// 返回文件访问URL
	url := fmt.Sprintf("/uploads/%s", filename)
	c.JSON(http.StatusOK, gin.H{
		"status": 1,
		"msg":    "上传成功",
		"url":    url,
	})
}