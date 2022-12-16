package models

import "time"

type User struct {
	UserId     int `gorm:"primary_key;AUTO_INCREMENT"`
	Username   string
	Password   string
	LastAccess time.Time
}

func GetAllUsers(users *[]User) {
	Db.Find(users)
}

func GetUserById(user *User, key int) {
	Db.First(user, key)
}

func GetUserByName(user *User, name string) {
	v := Db.Where("username = ?", name).First(user)
	_ = v
}

func InsertUser(user *User) {
	Db.NewRecord(user)
	Db.Create(&user)
}

func UpdateUser(user *User) {
	Db.Save(user)
}

func DeleteUser(key string) {
	Db.Where("user_id = ?", key).Delete(&User{})
}
