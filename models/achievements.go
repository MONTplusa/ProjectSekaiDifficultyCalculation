package models

import "time"

type Achievement struct {
	AchievementId int `gorm:"primary_key"`
	SongId        int
	Type          int // 1:cleared,2:FC,3:AP
	CreatedAt     time.Time
	UpdatedAt     time.Time
}

func GetAllAchievements(achievements *[]Achievement) {
	Db.Find(achievements)
}

func GetSingleAchievement(achievement *Achievement, key string) {
	Db.First(achievement, key)
}

func InsertAchievement(achievement *Achievement) {
	Db.NewRecord(achievement)
	Db.Create(achievement)
}

func DeleteAchievement(key string) {
	Db.Where("achievement_id = ?", key).Delete(&Achievement{})
}
