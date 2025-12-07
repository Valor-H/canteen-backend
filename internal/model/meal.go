package model

type WeekMeal struct {
	WeekNumber int32  `json:"weekNumber"`
	WeekDay    string `json:"weekDay"`
	MealType   string `json:"mealType"`
	MealId     int16  `json:"mealId"`
}

type Meal struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

type Dish struct {
	Id   int16  `json:"id"`
	Name string `json:"name"`
}

type CreateMeal struct {
	Name    string  `json:"name"`
	DishIds []int16 `json:"dishIds"`
}

// 周菜单发布
type WeekMenu struct {
	Day      string            `json:"day"`
	MealType string            `json:"meal_type"`
	Window   string            `json:"window"`
	Dishes   map[string]string `json:"dishes"`
}