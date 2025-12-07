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
	Id       int
	NickName string
	Count    int
	DeptId   int
	CardNo   string
}