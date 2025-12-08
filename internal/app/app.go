package app

import (
	"context"
	"log"

	"canteen/internal/controller/card"
	"canteen/internal/controller/tempDirect"
	"canteen/internal/controller/user"
	"canteen/internal/controller/order_record_detail"
	"canteen/internal/infrastructure/cache"
	"canteen/internal/infrastructure/database"
	"canteen/pkg/utils"
	"database/sql"
)

type Application struct {
	db *sql.DB
}

// NewApplication 创建应用实例
func NewApplication() *Application {
	return &Application{}
}

// Initialize 初始化应用
func (app *Application) Initialize() error {
	// 初始化Redis连接
	cache.InitRedis()

	// 初始化数据库连接
	app.db = database.InitDb()

	// 注入数据库连接到控制器
	card.SetDB(app.db)
	tempDirect.SetDB(app.db)
	user.SetDB(app.db)
	order_record_detail.SetDB(app.db)

	// 更新每日餐食缓存
	log.Println("Updating daily meal cache on startup...")
	if err := utils.UpdateDailyMealCache(context.Background(), app.db, cache.RedisClient()); err != nil {
		log.Printf("Failed to update daily meal cache on startup: %v", err)
	}

	return nil
}

// StartBackgroundTasks 启动后台任务
func (app *Application) StartBackgroundTasks() {
	// 启动定时任务
	go utils.DailyLicenseCheck()
	go utils.DailyExpireOrderRecords(app.db)
	go utils.WeeklyGenerateSetmeal(app.db)
	go utils.DailyMealCacheUpdate(app.db, cache.RedisClient())
}

// Shutdown 关闭应用
func (app *Application) Shutdown() error {
	if app.db != nil {
		if err := app.db.Close(); err != nil {
			log.Printf("Failed to close database connection: %v", err)
			return err
		}
		log.Println("Database connection closed successfully")
	}
	return nil
}
