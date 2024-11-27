package database

import "gorm.io/gorm"

//////////////////////////////////////////////////////////////////////////////////////////////////////////////////////// 用户标识表

const tableUser = "user"
const columnUserID = "id"
const columnUserUsername = "username"
const columnUserNickname = "nickname"
const columnUserAvatar = "avatar"
const columnUserEmail = "email"
const columnUserPhone = "phone"

type User struct {
	gorm.Model
	ID       int    `gorm:"column:id; primaryKey; not null; unique; autoIncrement;"`
	Username string `gorm:"column:username;not null; unique;"`
	Nickname string `gorm:"column:nickname;"`
	Avatar   string `gorm:"column:avatar;"`
	Email    string `gorm:"column:email;unique; default:''"`
	Phone    string `gorm:"column:phone;unique; default:''"`
}

func (User) TableName() string {
	return tableUser
}
