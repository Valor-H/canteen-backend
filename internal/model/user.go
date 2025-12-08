package model

type Login struct {
	User     string `json:"user"`
	Password string `json:"password"`
}

type User struct {
	Id       string
	Name     string
	NickName string
}

type UserVo struct {
	UserId   int    `json:"userId"`
	NickName string `json:"nickName"`
	Count    int    `json:"count"`
	DeptId   int    `json:"deptId"`
	CardNo   string `json:"cardNo"`
}