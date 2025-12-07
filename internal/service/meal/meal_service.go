package meal

import (
	"context"
	"fmt"
	"log"
	"time"

	"canteen/internal/model"
	"canteen/internal/repository/meal"
	"github.com/go-redis/redis/v8"
)

type MealService interface {
	GetDishesByIds(dishIds []int16) ([]model.Dish, error)
	GetSetmealsByWeekNumber(weekNumber string) ([]model.WeekMeal, error)
	UpdateDailyMealCache() error
	GenerateWeeklySetmeals() error
	GenerateNextWeekSetmeals() error
	CheckIfWeeklySetmealGenerated() bool
}

type mealService struct {
	mealRepo meal.MealRepository
	redis    *redis.Client
}

func NewMealService(mealRepo meal.MealRepository, redisClient *redis.Client) MealService {
	return &mealService{
		mealRepo: mealRepo,
		redis:    redisClient,
	}
}

func (s *mealService) GetDishesByIds(dishIds []int16) ([]model.Dish, error) {
	return s.mealRepo.FindDishesByIds(dishIds)
}

func (s *mealService) GetSetmealsByWeekNumber(weekNumber string) ([]model.WeekMeal, error) {
	return s.mealRepo.FindSetmealsByWeekNumber(weekNumber)
}

func (s *mealService) UpdateDailyMealCache() error {
	ctx := context.Background()
	dateStr := time.Now().Format("20060102")

	setmeals, err := s.mealRepo.FindSetmealsByWeekNumber(dateStr)
	if err != nil {
		return fmt.Errorf("db query failed: %w", err)
	}

	mealTypeMap := map[string]string{"午餐": "lunch", "晚餐": "dinner"}
	remarkMap := map[string]string{"套餐A": "A", "套餐B": "B", "套餐C": "C"}

	found := make(map[string]bool)

	for _, setmeal := range setmeals {
		// 这里需要从数据库查询获取remark，简化处理
		remark := "套餐A" // 假设值，实际应从数据库获取
		mealTypeEn, ok1 := mealTypeMap[setmeal.MealType]
		remarkEn, ok2 := remarkMap[remark]

		if !ok1 || !ok2 {
			log.Printf("Unknown mealType or remark: %s, %s", setmeal.MealType, remark)
			continue
		}

		key := fmt.Sprintf("%s-%s-%s", dateStr, mealTypeEn, remarkEn)
		if err := s.redis.Set(ctx, key, setmeal.MealId, 24*time.Hour).Err(); err != nil {
			log.Printf("Redis SET failed for key %s: %v", key, err)
		} else {
			log.Printf("Set Redis: %s => %d", key, setmeal.MealId)
			found[key] = true
		}
	}

	defaultKeys := []string{
		fmt.Sprintf("%s-lunch-A", dateStr),
		fmt.Sprintf("%s-lunch-B", dateStr),
		fmt.Sprintf("%s-lunch-C", dateStr),
		fmt.Sprintf("%s-dinner-A", dateStr),
		fmt.Sprintf("%s-dinner-C", dateStr),
	}

	for _, key := range defaultKeys {
		if !found[key] {
			if err := s.redis.Set(ctx, key, 1, 24*time.Hour).Err(); err != nil {
				log.Printf("Redis SET default failed for %s: %v", key, err)
			} else {
				log.Printf("Set default Redis: %s => 1", key)
			}
		}
	}

	return nil
}

func (s *mealService) GenerateWeeklySetmeals() error {
	nextMonday := getNextMonday(time.Now())
	dates := make([]time.Time, 6)
	for i := 0; i < 6; i++ {
		dates[i] = nextMonday.AddDate(0, 0, i)
	}

	return s.mealRepo.GenerateWeeklySetmeals(dates)
}

func (s *mealService) GenerateNextWeekSetmeals() error {
	return s.GenerateWeeklySetmeals()
}

func (s *mealService) CheckIfWeeklySetmealGenerated() bool {
	// 获取下周一的日期
	nextMonday := getNextMonday(time.Now())
	dates := make([]time.Time, 6)
	for i := 0; i < 6; i++ {
		dates[i] = nextMonday.AddDate(0, 0, i)
	}

	// 查询是否已有记录
	for _, date := range dates {
		weekNumber := date.Format("20060102")
		setmeals, err := s.mealRepo.FindSetmealsByWeekNumber(weekNumber)
		if err != nil {
			log.Printf("Failed to check weekly setmeal: %v", err)
			return false
		}
		if len(setmeals) == 0 {
			// 如果有任何一天没有生成套餐，则返回 false
			return false
		}
	}
	return true
}

// getNextMonday 获取下周一的日期
func getNextMonday(t time.Time) time.Time {
	daysUntilMonday := (time.Monday - t.Weekday() + 7) % 7
	return t.AddDate(0, 0, int(daysUntilMonday))
}