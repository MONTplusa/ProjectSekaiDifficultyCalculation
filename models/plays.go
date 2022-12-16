package models

import (
	"math"
	"time"
)

type Play struct {
	PlayId    int `gorm:"primary_key;AUTO_INCREMENT"`
	SongId    int
	Perfect   int
	Great     int
	Good      int
	Bad       int
	Miss      int
	Cleared   int
	PlayerId  int
	CreatedAt time.Time
	UpdatedAt time.Time
}

func NewPlay(songId, cleared, perfect, great, good, bad, miss, playerId int) Play {
	play := Play{}
	play.SongId = songId
	play.PlayerId = playerId
	play.Cleared = cleared
	if cleared == 0 {
		return play
	}
	play.Perfect = perfect
	play.Great = great
	play.Good = good
	play.Bad = bad
	play.Miss = miss
	return play
}

func GetAllPlays(plays *[]Play) {
	Db.Find(plays)
}

func GetAllPlaysByUserId(plays *[]Play, userId int) {
	Db.Where("player_id = ?", userId).Find(plays)
}

func InsertPlay(play *Play) {
	Db.NewRecord(play)
	Db.Create(play)
}

func DeletePlay(key string) {
	Db.Where("PlayId = ?", key).Delete(Play{})
}

func (play Play) ToClearScore() float64 {
	score := 0.0
	if play.Cleared == 0 {
		return -1000
	}
	damage := 1000.0
	damage -= 1.2 * float64(play.Great)
	damage -= 13 * float64(play.Good)
	damage -= 35 * float64(play.Bad)
	damage -= 63 * float64(play.Miss)
	score = damageToScore(damage)*2.0 - 1000
	return score
}

func (play Play) ToFullComboScore() float64 {
	score := 0.0
	if play.Cleared == 0 {
		return -1000
	}
	minus := 0
	minus += play.Great
	minus += play.Good
	minus += play.Bad
	minus += play.Miss
	if minus == 0 {
		return 1000
	}
	damage := 0.0
	damage += 6 * float64(play.Great)
	damage += 55 * float64(play.Good)
	damage += 80 * float64(play.Bad)
	damage += 60 * float64(play.Miss)
	if damage > 1000 {
		damage = 1000
	}
	score = -1.95*damageToScore(damage) + 950
	return score
}

func (play Play) ToAllPerfectScore() float64 {
	score := 0.0
	if play.Cleared == 0 {
		return -1000
	}
	minus := 0
	minus += play.Great
	minus += play.Good
	minus += play.Bad
	minus += play.Miss
	if minus == 0 {
		return 1000
	}
	damage := 0.0
	damage += 55 * float64(play.Great)
	damage += 105 * float64(play.Good)
	damage += 150 * float64(play.Bad)
	damage += 90 * float64(play.Miss)
	if damage > 1000 {
		damage = 1000
	}
	score = -1.95*damageToScore(damage) + 950
	return score
}

func damageToScore(damage float64) float64 {
	if damage < 0 {
		damage = 0
	}
	return 1000 - math.Sqrt((1000-damage)/1000)*(1000-damage)
}
