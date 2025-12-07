package model

type ConsumTransaction struct {
	Order    string `json:"Order"`
	CardNo   string `json:"CardNo"`
	CardMode int    `json:"CardMode"`
	Mode     int    `json:"Mode"`
	PayType  int    `json:"PayType"`
	Amount   string `json:"Amount"`
	Menus    []Menu `json:"Menus"`
}

type Menu struct {
	MenuID string `json:"MenuID"`
	Count  string `json:"Count"`
}

type OrderRecord struct {
	Id         int
	UserId     int
	Status     string
	MealId     int
	MealType   string // 添加餐类型字段
	WeekNumber string // 添加周数字段
	OrderDate  string // 添加订单日期字段
	Weekday    string // 添加星期字段
}

type OffLineRequest struct {
	DeviceNumber int    `json:"DeviceNumber"`
	Order        string `json:"Order"`
	PayType      int    `json:"PayType"`
	CardMode     int    `json:"CardMode"`
	Time         string `json:"Time"`
	CardNo       string `json:"CardNo"`
	Money        string `json:"Money"`
}