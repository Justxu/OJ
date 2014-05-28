package models

type User struct {
	Id       int64 `xorm:"pk"`
	Problems int64 //Number of solved problems
}
