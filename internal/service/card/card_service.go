package card

import (
	"context"
	"errors"
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	"canteen/internal/model"
	"canteen/internal/repository/card"
	"canteen/internal/repository/order"
	"canteen/internal/repository/user"

	"github.com/go-redis/redis/v8"
)

type CardService interface {
	ProcessConsumTransaction(req model.ConsumTransaction, deviceID string) (*model.ConsumResponse, error)
	GetServerTime() time.Time
	ProcessOffLineRequest(req model.OffLineRequest) error
}

type ConsumResponse struct {
	Status     int
	Message    string
	Name       string
	CardNo     string
	Money      int
	Subsidy    float64
	Times      int
	Integral   float64
	InTime     string
	OutTime    string
	Cumulative string
	Amount     string
	VoiceID    string
	Text       string
}

type cardService struct {
	userRepo  user.UserRepository
	orderRepo order.OrderRepository
	cardRepo  card.CardRepository
	redis     *redis.Client
}

func NewCardService(userRepo user.UserRepository, orderRepo order.OrderRepository, cardRepo card.CardRepository, redisClient *redis.Client) CardService {
	return &cardService{
		userRepo:  userRepo,
		orderRepo: orderRepo,
		cardRepo:  cardRepo,
		redis:     redisClient,
	}
}

func (s *cardService) ProcessConsumTransaction(req model.ConsumTransaction, deviceID string) (*model.ConsumResponse, error) {
	ctx := context.Background()
	deviceToRemark := map[string]string{
		"0180800116": "A",
		"0127448632": "B",
		"0158577664": "C",
	}

	log.Printf("TAG: 核销开始")
	log.Printf("传入卡号=%s", req.CardNo)

	// 查询用户信息
	user, err := s.cardRepo.FindUserByCardNo(req.CardNo)
	if err != nil {
		log.Printf("TAG: 查询用户信息失败: %v", err)
		return nil, fmt.Errorf("查询用户信息失败: %v", err)
	}

	log.Printf("TAG: 获取到用户信息 user_id=%d,名称=%s,卡号=%s", user.Id, user.NickName, user.CardNo)

	now := time.Now()
	dateStr := now.Format("20060102")
	weekdays := [...]string{"周日", "周一", "周二", "周三", "周四", "周五", "周六"}
	weekday := weekdays[now.Weekday()]

	// 判断餐类型
	var mealType string
	if now.Hour() >= 11 && now.Hour() < 14 {
		mealType = "午餐"
	} else {
		mealType = "晚餐"
	}

	// 客户处理逻辑
	if user.DeptId == 219 {
		log.Printf("TAG: 客户刷卡 dept_id=219")

		mealID, err := s.getMealIDFromRedis(ctx, deviceID, deviceToRemark, dateStr, now)
		if err != nil {
			return nil, fmt.Errorf("获取套餐ID失败: %v", err)
		}

		// 创建临时订单
		order := &model.OrderRecord{
			UserId:     user.Id,
			MealId:     mealID,
			Status:     "临时用餐",
			MealType:   mealType,
			WeekNumber: dateStr,
			OrderDate:  now.Format("2006-01-02"),
			Weekday:    weekday,
		}

		// 开始事务
		if err := s.createTempOrderAndDecreaseCount(order, user.Id, user.Count); err != nil {
			return nil, err
		}

		return &model.ConsumResponse{
			Status:     1,
			Message:    "核销成功:" + mealType,
			Name:       user.NickName,
			CardNo:     req.CardNo,
			Money:      0,
			Subsidy:    0.00,
			Times:      user.Count - 1,
			Integral:   0.00,
			InTime:     "",
			OutTime:    "",
			Cumulative: "",
			Amount:     req.Amount,
			VoiceID:    "核销成功",
			Text:       user.NickName + ":" + mealType + "核销成功",
		}, nil
	}

	// 员工处理逻辑
	if mealType == "晚餐" {
		ok, msg := s.checkDinnerTime(ctx, user.DeptId, now)
		if !ok {
			return nil, errors.New(msg)
		}
	}

	// 查询订单记录
	order, err := s.cardRepo.FindOrderRecord(user.Id, mealType, dateStr, weekday)
	if err != nil {
		log.Printf("TAG: 查询订单失败: %v", err)
		return nil, fmt.Errorf("查询订单失败: %v", err)
	}

	log.Printf("TAG: 查询订单结果: %+v", order)

	// 检查是否需要创建临时订单
	isUnordered := order.Id == 0 || order.Status == "" || order.MealId == 0
	if isUnordered {
		log.Printf("TAG: 未报餐，创建临时订单")

		mealID, err := s.getMealIDFromRedis(ctx, deviceID, deviceToRemark, dateStr, now)
		if err != nil {
			return nil, fmt.Errorf("获取套餐ID失败: %v", err)
		}

		tempOrder := &model.OrderRecord{
			UserId:     user.Id,
			MealId:     mealID,
			Status:     "临时用餐",
			MealType:   mealType,
			WeekNumber: dateStr,
			OrderDate:  now.Format("2006-01-02"),
			Weekday:    weekday,
		}

		if err := s.createTempOrderAndDecreaseCount(tempOrder, user.Id, user.Count); err != nil {
			return nil, err
		}

		return &model.ConsumResponse{
			Status:     1,
			Message:    "核销成功:" + mealType,
			Name:       user.NickName,
			CardNo:     req.CardNo,
			Money:      0,
			Subsidy:    0.00,
			Times:      user.Count - 1,
			Integral:   0.00,
			InTime:     "",
			OutTime:    "",
			Cumulative: "",
			Amount:     req.Amount,
			VoiceID:    "核销成功",
			Text:       user.NickName + ":" + mealType + "核销成功",
		}, nil
	}

	// 检查是否重复刷卡
	if order.Status == "已领取" || order.Status == "临时用餐" {
		log.Printf("TAG: 重复刷卡，订单已领取")
		return nil, fmt.Errorf("该卡今天%s重复刷卡取餐！", mealType)
	}

	// 检查窗口是否正确
	window, ok := deviceToRemark[deviceID]
	if !ok {
		log.Printf("TAG: 未知设备ID: %s", deviceID)
		return nil, fmt.Errorf("未知设备，无法取餐")
	}

	// 除周六外，其他日期不可刷其他套餐
	if weekday != "周六" {
		mealTypeMap := map[string]string{"午餐": "lunch", "晚餐": "dinner"}
		redisKey := fmt.Sprintf("%s-%s-%s", dateStr, mealTypeMap[mealType], window)

		cachedMealIDStr, err := s.redis.Get(ctx, redisKey).Result()
		if err != nil {
			log.Printf("TAG: Redis 获取失败, key=%s, err=%v", redisKey, err)
			return nil, fmt.Errorf("窗口配置读取失败")
		}

		cachedMealID, err := strconv.Atoi(cachedMealIDStr)
		if err != nil {
			log.Printf("TAG: Redis 缓存的 MealID 解析失败: %v", err)
			return nil, fmt.Errorf("系统配置异常")
		}

		if cachedMealID != order.MealId {
			log.Printf("TAG: 用户刷错窗口, 正确套餐ID=%d, 当前窗口套餐ID=%d", order.MealId, cachedMealID)
			return nil, fmt.Errorf("请前往正确的窗口刷卡取餐")
		}
	}

	// 更新订单状态并减少用户次数
	if err := s.updateOrderStatusAndDecreaseCount(order.Id, "已领取", user.Id, user.Count); err != nil {
		return nil, err
	}

	log.Printf("TAG: 核销成功: %s，用户: %s", mealType, user.NickName)

	return &model.ConsumResponse{
		Status:     1,
		Message:    "核销成功:" + mealType,
		Name:       user.NickName,
		CardNo:     req.CardNo,
		Money:      0,
		Subsidy:    0.00,
		Times:      user.Count - 1,
		Integral:   0.00,
		InTime:     "",
		OutTime:    "",
		Cumulative: "",
		Amount:     req.Amount,
		VoiceID:    "核销成功",
		Text:       user.NickName + ":" + mealType + "核销成功",
	}, nil
}

// 从Redis获取套餐ID
func (s *cardService) getMealIDFromRedis(ctx context.Context, deviceID string, deviceToRemark map[string]string, dateStr string, now time.Time) (int, error) {
	var mealTypeEn string
	if now.Hour() >= 11 && now.Hour() < 14 {
		mealTypeEn = "lunch"
	} else {
		mealTypeEn = "dinner"
	}

	remarkEn, ok := deviceToRemark[deviceID]
	if !ok {
		// 设备ID未映射到套餐类型，默认A
		remarkEn = "A"
	}

	key := fmt.Sprintf("%s-%s-%s", dateStr, mealTypeEn, remarkEn)
	setmealIDStr, err := s.redis.Get(ctx, key).Result()
	if err != nil {
		if err == redis.Nil {
			// 缓存不存在，使用默认值
			log.Printf("Redis key not found: %s, using default ID 1", key)
			return 1, nil
		}
		log.Printf("Redis GET failed for key %s: %v", key, err)
		return 0, err
	}

	mealId, err := strconv.Atoi(setmealIDStr)
	if err != nil {
		log.Printf("Convert setmealID failed for key %s: %v", key, err)
		return 1, nil
	}

	return mealId, nil
}

// 检查晚餐时间
func (s *cardService) checkDinnerTime(ctx context.Context, userDeptId int, now time.Time) (bool, string) {
	const redisKey = "canteen:dinner_config"
	var (
		flexibleDeptIdsStr, fixedDeptIdsStr, flexibleDinnerStart, fixedDinnerStart string
	)

	configs, err := s.redis.HGetAll(ctx, redisKey).Result()
	if err != nil || len(configs) == 0 {
		log.Println("TAG: Redis 缓存未命中，查询数据库")
		flexibleDeptIdsStr, fixedDeptIdsStr, flexibleDinnerStart, fixedDinnerStart, err = s.cardRepo.GetCanteenConfigs()
		if err != nil {
			log.Printf("TAG: 查询晚餐时间配置失败: %v", err)
			return false, "系统配置错误"
		}

		s.redis.HSet(ctx, redisKey, "flexible_dept_id", flexibleDeptIdsStr)
		s.redis.HSet(ctx, redisKey, "fixed_dept_id", fixedDeptIdsStr)
		s.redis.HSet(ctx, redisKey, "flexible_dinner_start_time", flexibleDinnerStart)
		s.redis.HSet(ctx, redisKey, "fixed_dinner_start_time", fixedDinnerStart)
		s.redis.Expire(ctx, redisKey, 120*time.Hour)
	} else {
		log.Println("TAG: 命中 Redis 缓存配置")
		flexibleDeptIdsStr = configs["flexible_dept_id"]
		fixedDeptIdsStr = configs["fixed_dept_id"]
		flexibleDinnerStart = configs["flexible_dinner_start_time"]
		fixedDinnerStart = configs["fixed_dinner_start_time"]
	}

	deptIdStr := strconv.Itoa(userDeptId)
	flexibleDeptIds := strings.Split(flexibleDeptIdsStr, ",")
	fixedDeptIds := strings.Split(fixedDeptIdsStr, ",")

	var isFlexible, isFixed bool
	for _, id := range flexibleDeptIds {
		if id == deptIdStr {
			isFlexible = true
			break
		}
	}
	if !isFlexible {
		for _, id := range fixedDeptIds {
			if id == deptIdStr {
				isFixed = true
				break
			}
		}
	}

	if !isFlexible && !isFixed {
		log.Printf("TAG: 部门未配置用餐规则 dept_id=%d", userDeptId)
		return false, "部门未配置用餐规则"
	}

	var dinnerStart string
	if isFlexible {
		dinnerStart = flexibleDinnerStart
		log.Printf("TAG: 弹性用餐部门，晚餐开始时间: %s", dinnerStart)
	} else {
		dinnerStart = fixedDinnerStart
		log.Printf("TAG: 固定用餐部门，晚餐开始时间: %s", dinnerStart)
	}

	startParts := strings.Split(dinnerStart, ":")
	startHour, _ := strconv.Atoi(startParts[0])
	startMin, _ := strconv.Atoi(startParts[1])

	dinnerStartTime := time.Date(now.Year(), now.Month(), now.Day(), startHour, startMin, 0, 0, now.Location())
	dinnerEndTime := time.Date(now.Year(), now.Month(), now.Day(), 21, 0, 0, 0, now.Location())

	if now.Before(dinnerStartTime) || now.After(dinnerEndTime) {
		log.Printf("TAG: 当前时间 %v 不在晚餐时间段 %v ~ %v", now, dinnerStartTime, dinnerEndTime)
		return false, "不在就餐时间范围内"
	}
	return true, ""
}

// 创建临时订单并减少用户次数
func (s *cardService) createTempOrderAndDecreaseCount(order *model.OrderRecord, userId int, count int) error {
	// 使用数据库事务
	if err := s.cardRepo.CreateOrderRecord(order); err != nil {
		log.Printf("TAG: 创建临时订单失败: %v", err)
		return fmt.Errorf("创建订单失败: %v", err)
	}

	if err := s.cardRepo.UpdateUserCount(userId, count-1); err != nil {
		log.Printf("TAG: 扣除次数失败: %v", err)
		return fmt.Errorf("扣次数失败: %v", err)
	}

	return nil
}

// 更新订单状态并减少用户次数
func (s *cardService) updateOrderStatusAndDecreaseCount(orderId int, status string, userId int, count int) error {
	// 更新订单状态
	if err := s.cardRepo.UpdateOrderStatus(orderId, status); err != nil {
		log.Printf("TAG: 更新订单失败: %v", err)
		return fmt.Errorf("更新失败: %v", err)
	}

	// 减少用户次数
	if err := s.cardRepo.UpdateUserCount(userId, count-1); err != nil {
		log.Printf("TAG: 扣除次数失败: %v", err)
		return fmt.Errorf("扣次数失败: %v", err)
	}

	return nil
}

func (s *cardService) GetServerTime() time.Time {
	return time.Now()
}

func (s *cardService) ProcessOffLineRequest(req model.OffLineRequest) error {
	// 处理离线请求
	log.Printf("Processing offline request: %+v", req)
	return nil
}
