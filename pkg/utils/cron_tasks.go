package utils

import (
	"context"
	"crypto"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"database/sql"
	"encoding/base64"
	"encoding/json"
	"encoding/pem"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/go-redis/redis/v8"
)

type License struct {
	ExpiryTime string `json:"expiry_time"`
	Signer     string `json:"signer"`
	Email      string `json:"email"`
	Signature  string `json:"signature"`
}

func loadPublicKeyFromEmbed() (*rsa.PublicKey, error) {
	block, _ := pem.Decode(publicKeyPEM)
	if block == nil {
		return nil, fmt.Errorf("failed to decode embedded public.pem")
	}
	pubAny, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return nil, err
	}
	pubKey, ok := pubAny.(*rsa.PublicKey)
	if !ok {
		return nil, fmt.Errorf("embedded key is not RSA public key")
	}
	return pubKey, nil
}
func CheckLicenseFileIntegrity(path string) bool {
	return true
}
func ValidateLicense() bool {
	exePath, err := os.Getwd()
	if err != nil {
		log.Printf("无法获取可执行路径：%v", err)
		return false
	}

	licenseFile := filepath.Join(exePath, "config", "license.lic")

	log.Printf("许可证文件路径: %s", licenseFile)

	if !CheckLicenseFileIntegrity(licenseFile) {
		return false
	}

	data, err := os.ReadFile(licenseFile)
	if err != nil {
		return false
	}

	var license License
	if err := json.Unmarshal(data, &license); err != nil {
		return false
	}

	publicKey, err := loadPublicKeyFromEmbed()
	if err != nil {
		return false
	}

	signatureBytes, err := base64.StdEncoding.DecodeString(license.Signature)
	if err != nil {
		return false
	}

	temp := License{
		ExpiryTime: license.ExpiryTime,
		Signer:     license.Signer,
		Email:      license.Email,
	}
	toVerify, err := json.Marshal(temp)
	if err != nil {
		return false
	}
	h := sha256.Sum256(toVerify)
	if err := rsa.VerifyPKCS1v15(publicKey, crypto.SHA256, h[:], signatureBytes); err != nil {
		return false
	}

	expiry, err := time.Parse(time.RFC3339, license.ExpiryTime)
	if err != nil {
		return false
	}
	if time.Now().After(expiry) {
		return false
	}
	return true
}

func DailyLicenseCheck() {
	for {
		now := time.Now()
		next := now.AddDate(0, 0, 1)
		next = time.Date(next.Year(), next.Month(), next.Day(), 2, 0, 0, 0, next.Location())
		time.Sleep(next.Sub(now))

		if !ValidateLicense() {
			log.Fatal("Daily license check failed. Exiting...")
		}
	}
}

// 2.DailyExpireOrderRecords 定时任务：每日处理过期订单
func DailyExpireOrderRecords(db *sql.DB) {
	for {
		now := time.Now()
		var next time.Time
		if now.Hour() >= 23 {
			next = time.Date(now.Year(), now.Month(), now.Day()+1, 23, 0, 0, 0, now.Location())
		} else {
			next = time.Date(now.Year(), now.Month(), now.Day(), 23, 0, 0, 0, now.Location())
		}
		sleepDuration := next.Sub(now)
		log.Printf("Next expire task scheduled at: %v", next)
		time.Sleep(sleepDuration)

		todayStr := time.Now().Format("20060102")
		todayInt, _ := strconv.Atoi(todayStr)

		tx, err := db.Begin()
		if err != nil {
			log.Printf("Failed to begin transaction: %v", err)
			continue
		}

		// Step 1: 更新订单记录并获取受影响的 user_id 列表
		queryUpdateOrders := `
				UPDATE order_record 
				SET status = '已过期', update_time = NOW() 
				WHERE week_number = ? AND status = '已报餐'
			`
		result, err := tx.Exec(queryUpdateOrders, todayInt)
		if err != nil {
			tx.Rollback()
			log.Printf("Failed to update expired orders: %v", err)
			continue
		}

		rowsAffected, _ := result.RowsAffected()
		log.Printf("Marked %d orders as expired for day %s", rowsAffected, todayStr)

		// Step 2: 将对应的 sys_user.count -1
		if rowsAffected > 0 {
			queryDecrementUserCount := `
					UPDATE sys_user u
					JOIN (
						SELECT DISTINCT user_id
						FROM order_record
						WHERE week_number = ? AND status = '已过期'
					) AS o USING (user_id)
					SET u.count = GREATEST(u.count - 1, 0)
				`

			_, err = tx.Exec(queryDecrementUserCount, todayInt)
			if err != nil {
				tx.Rollback()
				log.Printf("Failed to decrement user count: %v", err)
			} else {
				tx.Commit()
				log.Printf("Successfully decremented count for users with expired orders.")
			}
		} else {
			tx.Commit()
		}
	}
}

// 3.WeeklyGenerateSetmeal 定时任务：每周生成套餐
func WeeklyGenerateSetmeal(db *sql.DB) {
	for {
		now := time.Now()
		nextThursday := FindNextThursday(now)
		log.Printf("Next scheduled generation at: %v", nextThursday)
		time.Sleep(nextThursday.Sub(now))

		// 生成下周套餐
		if err := GenerateNextWeekSetmeals(db); err != nil {
			log.Printf("Failed to generate setmeals: %v", err)
		}
	}
}

func FindNextThursday(t time.Time) time.Time {
	// 计算距离下个周四的天数
	daysUntilThursday := (time.Thursday - t.Weekday() + 7) % 7
	nextThursday := t.AddDate(0, 0, int(daysUntilThursday))
	// 设置为9:00
	nextThursday = time.Date(nextThursday.Year(), nextThursday.Month(), nextThursday.Day(), 10, 0, 0, 0, nextThursday.Location())
	// 如果当前时间已过本周四12:00，则跳至下周
	if nextThursday.Before(t) {
		nextThursday = nextThursday.AddDate(0, 0, 7)
	}
	return nextThursday
}

// GenerateNextWeekSetmeals 生成下周的套餐记录
func GenerateNextWeekSetmeals(db *sql.DB) error {
	// 获取下周一的日期
	nextMonday := GetNextMonday(time.Now())
	dates := make([]time.Time, 6)
	for i := 0; i < 6; i++ {
		dates[i] = nextMonday.AddDate(0, 0, i)
	}

	// 开启事务
	tx, err := db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// 删除旧记录
	startWeek := dates[0].Format("20060102")
	endWeek := dates[5].Format("20060102")
	if _, err := tx.Exec("DELETE FROM weekly_setmeal WHERE week_number BETWEEN ? AND ?", startWeek, endWeek); err != nil {
		return err
	}

	// 插入新记录
	for _, date := range dates {
		weekNumber := date.Format("20060102")
		weekday := GetWeekdayZh(date)

		// 定义午餐和晚餐的备注
		lunchRemarks := []string{"套餐A", "套餐B", "套餐C"}
		//dinnerRemark := []string{"套餐A", "套餐C"}
		dinnerRemark := []string{"套餐A", "套餐C"}

		// 插入3个午餐
		for _, remark := range lunchRemarks {
			_, err := tx.Exec(`
                INSERT INTO weekly_setmeal 
                    (week_number, weekday, meal_type, setmeal_id, create_time, create_user, remark)
                VALUES (?, ?, ?, NULL, NOW(), 263, ?)`,
				weekNumber, weekday, "午餐", remark)
			if err != nil {
				return err
			}
		}

		// 插入晚餐
		for _, remark := range dinnerRemark {
			_, err := tx.Exec(`
                INSERT INTO weekly_setmeal 
                    (week_number, weekday, meal_type, setmeal_id, create_time, create_user, remark)
                VALUES (?, ?, ?, NULL, NOW(), 263, ?)`,
				weekNumber, weekday, "晚餐", remark)
			if err != nil {
				return err
			}
		}
	}
	// 提交事务
	if err := tx.Commit(); err != nil {
		return err
	}
	log.Printf("Generated setmeals for week starting %s", dates[0].Format("2006-01-02"))
	return nil
}

// GetNextMonday 获取下周一的日期
func GetNextMonday(t time.Time) time.Time {
	daysUntilMonday := (time.Monday - t.Weekday() + 7) % 7
	return t.AddDate(0, 0, int(daysUntilMonday))
}

// GetWeekdayZh 获取zn周期
func GetWeekdayZh(t time.Time) string {
	switch t.Weekday() {
	case time.Monday:
		return "周一"
	case time.Tuesday:
		return "周二"
	case time.Wednesday:
		return "周三"
	case time.Thursday:
		return "周四"
	case time.Friday:
		return "周五"
	case time.Saturday:
		return "周六"
	case time.Sunday:
		return "周日"
	}
	return ""
}

// CheckIfWeeklySetmealGenerated 检查下一周的套餐是否已生成
func CheckIfWeeklySetmealGenerated(db *sql.DB) bool {
	// 获取下周一的日期
	nextMonday := GetNextMonday(time.Now())
	dates := make([]time.Time, 6)
	for i := 0; i < 6; i++ {
		dates[i] = nextMonday.AddDate(0, 0, i)
	}

	// 查询是否已有记录
	for _, date := range dates {
		weekNumber := date.Format("20060102")
		query := `SELECT COUNT(*) FROM weekly_setmeal WHERE week_number = ?`
		var count int
		err := db.QueryRow(query, weekNumber).Scan(&count)
		if err != nil {
			log.Printf("Failed to check weekly setmeal: %v", err)
			return false
		}
		if count == 0 {
			// 如果有任何一天没有生成套餐，则返回 false
			return false
		}
	}
	return true
}
func DailyMealCacheUpdate(db *sql.DB, redisClient *redis.Client) {
	ctx := context.Background()

	for {
		now := time.Now()
		// 计算下一个凌晨5点的时间点
		next := time.Date(now.Year(), now.Month(), now.Day(), 5, 0, 0, 0, now.Location())
		if !now.Before(next) {
			next = next.Add(24 * time.Hour)
		}
		log.Printf("meal cache update scheduled at: %v", next)
		time.Sleep(time.Until(next))

		if err := UpdateDailyMealCache(ctx, db, redisClient); err != nil {
			log.Printf("Daily meal cache update failed: %v", err)
		}
	}
}

func UpdateDailyMealCache(ctx context.Context, db *sql.DB, redisClient *redis.Client) error {
	dateStr := time.Now().Format("20060102")

	query := `
		SELECT id, meal_type, remark FROM weekly_setmeal
		WHERE week_number = ? AND meal_type IN ('午餐', '晚餐')
	`

	rows, err := db.Query(query, dateStr)
	if err != nil {
		return fmt.Errorf("db query failed: %w", err)
	}
	defer rows.Close()

	mealTypeMap := map[string]string{"午餐": "lunch", "晚餐": "dinner"}
	remarkMap := map[string]string{"套餐A": "A", "套餐B": "B", "套餐C": "C"}

	found := make(map[string]bool)

	for rows.Next() {
		var id int
		var mealType, remark string
		if err := rows.Scan(&id, &mealType, &remark); err != nil {
			log.Printf("Row scan failed: %v", err)
			continue
		}

		mealTypeEn, ok1 := mealTypeMap[mealType]
		remarkEn, ok2 := remarkMap[remark]

		if !ok1 || !ok2 {
			log.Printf("Unknown mealType or remark: %s, %s", mealType, remark)
			continue
		}

		key := fmt.Sprintf("%s-%s-%s", dateStr, mealTypeEn, remarkEn)
		if err := redisClient.Set(ctx, key, id, 24*time.Hour).Err(); err != nil {
			log.Printf("Redis SET failed for key %s: %v", key, err)
		} else {
			log.Printf("Set Redis: %s => %d", key, id)
			found[key] = true
		}
	}

	if err := rows.Err(); err != nil {
		return fmt.Errorf("rows iteration error: %w", err)
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
			if err := redisClient.Set(ctx, key, 1, 24*time.Hour).Err(); err != nil {
				log.Printf("Redis SET default failed for %s: %v", key, err)
			} else {
				log.Printf("Set default Redis: %s => 1", key)
			}
		}
	}

	return nil
}