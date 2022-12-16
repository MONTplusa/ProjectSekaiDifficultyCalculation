package models

import "time"

type Song struct {
	SongId     int `gorm:"primary_key;AUTO_INCREMENT"`
	Name       string
	Difficulty int // 0:Easy,...,4:Master
	Notes      int
	Level      int // 公式難易度 ~36
	CreatedAt  time.Time
	UpdatedAt  time.Time
}

func GetAllSongs(songs *[]Song) {
	Db.Find(&songs)
}

/*
名前と難易度から曲を取得する

返り値 true:存在 false:存在しない
*/
func GetSongByNameAndDiff(song *Song, name string, difficulty int) bool {
	Db.Where("name = ?", name).Where("difficulty = ?", difficulty).First(song)
	return song.Name != ""
}

func InsertSong(song *Song) {
	Db.NewRecord(song)
	Db.Create(song)
}

func UpdateSong(song *Song) {
	Db.Save(song)
}

func DeleteSong(key string) {
	Db.Where("song_id = ?", key).Delete(Song{})
}
