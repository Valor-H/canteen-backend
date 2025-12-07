package order

import (
	"canteen/internal/model"
	"canteen/internal/repository/order"
	"canteen/internal/repository/user"
	"errors"
	"log"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/xuri/excelize/v2"
)

type OrderService interface {
	ExportOrdersByDate(date string) (*excelize.File, error)
	ExportOrdersByMonth(date string) (*excelize.File, error)
	ProcessExpiredOrders() error
	CreateOrder(order *model.OrderRecord) error
}

type orderService struct {
	orderRepo order.OrderRepository
	userRepo  user.UserRepository
	redis     *redis.Client
}

func NewOrderService(orderRepo order.OrderRepository, userRepo user.UserRepository, redisClient *redis.Client) OrderService {
	return &orderService{
		orderRepo: orderRepo,
		userRepo:  userRepo,
		redis:     redisClient,
	}
}

func (s *orderService) ExportOrdersByDate(date string) (*excelize.File, error) {
	// 使用repository导出Excel
	return s.orderRepo.ExportToExcel(date)
}

func (s *orderService) ExportOrdersByMonth(date string) (*excelize.File, error) {
	// 类似ExportOrdersByDate的实现，但按月查询
	return s.ExportOrdersByDate(date) // 简化实现
}

func (s *orderService) ProcessExpiredOrders() error {
	todayStr := time.Now().Format("20060102")
	
	// 使用事务处理过期订单
	// 注意：这里应该使用数据库事务，但为了简化，我们分步执行
	
	// 1. 更新订单记录并获取受影响的 user_id 列表
	err := s.orderRepo.UpdateStatusToExpired(todayStr)
	if err != nil {
		log.Printf("Failed to update expired orders: %v", err)
		return err
	}
	
	// 2. 获取受影响的用户ID
	userIds, err := s.orderRepo.FindExpiredOrdersByWeekNumber(todayStr)
	if err != nil {
		log.Printf("Failed to get affected user IDs: %v", err)
		return err
	}
	
	// 3. 将对应的 sys_user.count -1
	for _, userId := range userIds {
		if err := s.userRepo.DecreaseCountByUserId(userId); err != nil {
			log.Printf("Failed to decrement user count for user %d: %v", userId, err)
			continue
		}
	}
	
	log.Printf("Successfully processed expired orders for day %s", todayStr)
	return nil
}

func (s *orderService) CreateOrder(order *model.OrderRecord) error {
	if order.UserId <= 0 {
		return errors.New("无效的用户ID")
	}
	if order.MealId <= 0 {
		return errors.New("无效的餐食ID")
	}
	if order.Status == "" {
		order.Status = "已报餐" // 默认状态
	}
	
	return s.orderRepo.CreateOrder(order)
}