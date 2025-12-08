package model

// OrderRecordDetail 点餐详情，包含从多个表关联查询的信息
type OrderRecordDetail struct {
	Id            int    `json:"id"`             // 点菜号
	UserId        int    `json:"userId"`         // 用户号
	SetmealId     int    `json:"setmealId"`      // 周菜单号
	DishId        int    `json:"dishId"`         // 菜品号
	DishCode      string `json:"dishCode"`       // 菜品编号
	Description   string `json:"description"`    // 菜品描述
	OrderDate     string `json:"orderDate"`      // 点餐日期
	MealType      string `json:"mealType"`       // 餐别（午餐/晚餐）
	WeekNumber    string `json:"weekNumber"`     // 周数
	Weekday       string `json:"weekday"`        // 星期几
	Status        string `json:"status"`          // 状态
}