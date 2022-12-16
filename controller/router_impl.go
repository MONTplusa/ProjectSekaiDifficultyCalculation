package controller

import (
	"bufio"
	"fmt"
	"math"
	"net/http"
	"os"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/MONTplusa/ProjectSekaiDifficultyCalculation/models"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

const (
	FILEPATH              = "./data/tries.csv"
	SONGPATH              = "./data/songs.csv"
	VOLATILITY_CLEAR      = 0.05
	VOLATILITY_FULLCOMBO  = 0.05
	VOLATILITY_ALLPERFECT = 0.05
)

var difficultyname = [5]string{"Easy", "Normal", "Hard", "Expert", "Master"}

type BaseData struct {
	Links    []Link
	Username string
	Title    string
	Data     gin.H
}

type Link struct {
	Url   string
	Title string
}

type Plate struct {
	Rate           string
	DifficultyName string
	Level          int
	Name           string
	Try            int
}

func NewBaseData(c *gin.Context) BaseData {
	bd := BaseData{}
	bd.Links = []Link{}
	bd.Links = append(bd.Links, Link{Url: "/", Title: "ホーム"})
	bd.Links = append(bd.Links, Link{Url: "/show/GENERAL", Title: "総合難易度"})
	bd.Links = append(bd.Links, Link{Url: "/show/CLEAR", Title: "クリア難易度"})
	bd.Links = append(bd.Links, Link{Url: "/show/FC", Title: "FC難易度"})
	bd.Links = append(bd.Links, Link{Url: "/show/AP", Title: "AP難易度"})
	bd.Links = append(bd.Links, Link{Url: "/form", Title: "プレイ結果を追加"})
	bd.Links = append(bd.Links, Link{Url: "/login", Title: "ログイン"})
	bd.Links = append(bd.Links, Link{Url: "/logout", Title: "ログアウト"})
	session := sessions.Default(c)
	name := session.Get("Username")
	if name == nil {
		bd.Username = ""
	} else {
		bd.Username = name.(string)
	}
	bd.Data = gin.H{}
	return bd
}

func add(c *gin.Context) {
	fmt.Fprintln(os.Stderr, "/add")
	name := c.PostForm("Song")
	cleared, _ := strconv.Atoi(c.PostForm("Cleared"))
	perfect, _ := strconv.Atoi(c.PostForm("PERFECT"))
	great, _ := strconv.Atoi(c.PostForm("GREAT"))
	good, _ := strconv.Atoi(c.PostForm("GOOD"))
	bad, _ := strconv.Atoi(c.PostForm("BAD"))
	miss, _ := strconv.Atoi(c.PostForm("MISS"))
	difficulty, _ := strconv.Atoi(c.PostForm("Difficulty"))
	addPlay(c, name, cleared, perfect, great, good, bad, miss, difficulty, "", false)
	if c.Errors.Last() == nil {
		session := sessions.Default(c)
		username := session.Get("Username").(string)
		file, _ := os.OpenFile(FILEPATH, os.O_APPEND|os.O_WRONLY, 0644)
		newLine := fmt.Sprintf("%s,%d,%d,%d,%d,%d,%d,%d,%s", name, perfect, great, good, bad, miss, cleared, difficulty, username)
		fmt.Fprintln(file, newLine)
		file.Close()
	}
}

func addPlay(c *gin.Context, name string, cleared, perfect, great, good, bad, miss, difficulty int, username string, autoUpdateSong bool) {
	var song models.Song
	exist := models.GetSongByNameAndDiff(&song, name, difficulty)
	if !exist {
		if autoUpdateSong {
			song.Name = name
			song.Notes = perfect + great + good + bad + miss
			models.InsertSong(&song)
		} else {
			err := InvalidRequest
			err.Detail = "未対応の曲であるか、名称が異なります。"
			c.Error(&err).SetType(gin.ErrorTypePublic).SetMeta("")
			return
		}
	}
	if cleared == 1 && song.Notes != perfect+great+good+bad+miss {
		if autoUpdateSong {
			song.Notes = perfect + great + good + bad + miss
			song.Difficulty = difficulty
			models.UpdateSong(&song)
		} else {
			err := InvalidRequest
			err.Detail = "ノーツ数が一致しません。もういちど確認してください。"
			c.Error(&err).SetType(gin.ErrorTypePublic).SetMeta("")
			return
		}
	}
	session := sessions.Default(c)
	userId := 0
	if username == "" {
		userId = session.Get("UserId").(int)
	} else {
		var user models.User
		models.GetUserByName(&user, username)
		if user.Username == "" {
			return
		}
		userId = user.UserId
	}
	play := models.NewPlay(song.SongId, cleared, perfect, great, good, bad, miss, userId)
	models.InsertPlay(&play)
	fmt.Fprintln(os.Stderr, play)
}

func addSong(c *gin.Context) {
	name := c.PostForm("Song")
	notes, _ := strconv.Atoi(c.PostForm("Notes"))
	level, _ := strconv.Atoi(c.PostForm("Level"))
	difficulty, _ := strconv.Atoi(c.PostForm("Difficulty"))
	var song models.Song
	models.GetSongByNameAndDiff(&song, name, difficulty)
	song.Name = name
	song.Notes = notes
	song.Difficulty = difficulty
	song.Level = level
	if song.Name == "" {
		models.InsertSong(&song)
	} else {
		models.UpdateSong(&song)
	}
	if c.Errors.Last() == nil {
		file, _ := os.OpenFile(SONGPATH, os.O_APPEND|os.O_WRONLY, 0644)
		newLine := fmt.Sprintf("%s,%d,%d,%d", name, level, notes, difficulty)
		fmt.Fprintln(file, newLine)
		file.Close()
	}
}

func createUser(c *gin.Context) {
	name := c.PostForm("Username")
	password := c.PostForm("Password")
	if name == "" || password == "" {
		err := InvalidRequest
		err.Detail = "名前またはパスワードが入力されていません"
		c.Error(&err).SetType(gin.ErrorTypePublic).SetMeta("")
		return
	}
	var user models.User
	models.GetUserByName(&user, name)
	if user.Username != "" {
		err := InvalidRequest
		err.Detail = "既に同じ名前のユーザーが存在します"
		c.Error(&err).SetType(gin.ErrorTypePublic).SetMeta("")
		return
	}
	hashed, _ := bcrypt.GenerateFromPassword([]byte(password), 10)
	user.Username = name
	user.Password = string(hashed)
	user.LastAccess = time.Now()
	models.InsertUser(&user)
}

func form(c *gin.Context) {
	songs := []models.Song{}
	models.GetAllSongs(&songs)
	numSongs := len(songs)
	songnames := []string{}
	for i := 0; i < numSongs; i++ {
		if songs[i].Difficulty == 4 {
			songnames = append(songnames, songs[i].Name)
		}
	}
	data := NewBaseData(c)
	data.Data = gin.H{"songnames": songnames}
	c.HTML(http.StatusOK, "form.htm", data)
}

func form_createUser(c *gin.Context) {
	data := NewBaseData(c)
	c.HTML(http.StatusOK, "form_create_user.htm", data)
}

func form_login(c *gin.Context) {
	data := NewBaseData(c)
	c.HTML(http.StatusOK, "form_login.htm", data)

}

func form_NewSong(c *gin.Context) {
	data := NewBaseData(c)
	c.HTML(http.StatusOK, "form_new_song.htm", data)
}

func home(c *gin.Context) {
	data := NewBaseData(c)
	c.HTML(http.StatusOK, "home.htm", data)
}

func login(c *gin.Context) {
	name := c.PostForm("Username")
	password := c.PostForm("Password")
	var user models.User
	models.GetUserByName(&user, name)
	if user.Username == "" || user.Password == "" {
		err := InvalidRequest
		err.Detail = "ログインに失敗しました。名前またはパスワードをもう一度確認してください。"
		c.Error(&err).SetType(gin.ErrorTypePublic).SetMeta("")
		return
	} else {
		err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
		if err != nil {
			newerr := InvalidRequest
			newerr.Detail = "ログインに失敗しました。名前またはパスワードをもう一度確認してください。"
			c.Error(&newerr).SetType(gin.ErrorTypePublic).SetMeta(err.Error())
			return
		}
		session := sessions.Default(c)
		session.Set("UserId", user.UserId)
		session.Set("Username", user.Username)
		session.Save()
		c.Redirect(http.StatusMovedPermanently, "/show/GENERAL")
	}
}

func importData(c *gin.Context) {
	songs := []models.Song{}
	songIndexes := map[int]int{} // key:SongId
	models.GetAllSongs(&songs)
	numSongs := len(songs)
	for i := 0; i < numSongs; i++ {
		songIndexes[songs[i].SongId] = i
	}
	file, _ := os.Open(FILEPATH)
	scanner := bufio.NewScanner(file)
	for loop := 0; loop < 1<<20; loop++ {
		scanner.Scan()
		stringInput := scanner.Text()
		if stringInput == "" {
			break
		}
		splitedInput := strings.Split(stringInput, ",")
		if loop == 0 {
			splitedInput[0] = splitedInput[0][3:]
		}
		name := splitedInput[0]
		perfect, _ := strconv.Atoi(splitedInput[1])
		great, _ := strconv.Atoi(splitedInput[2])
		good, _ := strconv.Atoi(splitedInput[3])
		bad, _ := strconv.Atoi(splitedInput[4])
		miss, _ := strconv.Atoi(splitedInput[5])
		cleared, _ := strconv.Atoi(splitedInput[6])
		difficulty, _ := strconv.Atoi(splitedInput[7])
		username := splitedInput[8]
		addPlay(c, name, cleared, perfect, great, good, bad, miss, difficulty, username, false)
	}
	err := c.Errors.ByType(gin.ErrorTypePublic).Last()
	if err == nil {
		c.JSON(http.StatusAccepted, gin.H{})
	}

}

func importSongData(c *gin.Context) {
	file, _ := os.Open(SONGPATH)
	scanner := bufio.NewScanner(file)
	for loop := 0; loop < 1<<20; loop++ {
		scanner.Scan()
		stringInput := scanner.Text()
		if stringInput == "" {
			break
		}
		splitedInput := strings.Split(stringInput, ",")
		if loop == 0 {
			splitedInput[0] = splitedInput[0][3:]
		}
		name := splitedInput[0]
		level, _ := strconv.Atoi(splitedInput[1])
		notes, _ := strconv.Atoi(splitedInput[2])
		difficulty, _ := strconv.Atoi(splitedInput[3])
		var song models.Song
		models.GetSongByNameAndDiff(&song, name, difficulty)
		song.Name = name
		song.Notes = notes
		song.Difficulty = difficulty
		song.Level = level
		if song.Name == "" {
			models.InsertSong(&song)
		} else {
			models.UpdateSong(&song)
		}
	}
	c.JSON(http.StatusAccepted, gin.H{"result": "OK"})

}

func show(c *gin.Context) {
	mode := c.Param("mode")
	performances := map[int]models.Performance{} // key:SongId
	fmt.Fprintf(os.Stderr, "show_mode: %s\n", mode)
	songs := []models.Song{}
	songIndexes := map[int]int{} // key:SongId
	models.GetAllSongs(&songs)
	numSongs := len(songs)
	data := NewBaseData(c)
	for i := 0; i < numSongs; i++ {
		songIndexes[songs[i].SongId] = i
	}
	/*file, _ := os.Open(FILEPATH)
	scanner := bufio.NewScanner(file)*/
	plays := []models.Play{}
	session := sessions.Default(c)
	userId := session.Get("UserId").(int)
	models.GetAllPlaysByUserId(&plays, userId)
	numPlays := len(plays)
	for loop := 0; loop < numPlays; loop++ {
		/*
			scanner.Scan()
			stringInput := scanner.Text()
			if stringInput == "" {
				break
			}
			splitedInput := strings.Split(stringInput, ",")
			if loop == 0 {
				splitedInput[0] = splitedInput[0][3:]
			}
		*/
		play := (plays)[loop]
		if _, exist := songIndexes[play.SongId]; !exist {
			fmt.Fprintf(os.Stderr, "WARNING:SongId '%d' is unknown. ignored.(PlayId = %d)\n", play.SongId, play.PlayId)
			continue
		}
		performance, exist := performances[play.SongId]
		if !exist {
			performance = models.Performance{}
			performance.SongId = play.SongId
		}
		performance = update(performance, play)
		performances[play.SongId] = performance
	}
	performanceSlice := []models.Performance{}
	for _, v := range performances {
		if v.Try != 0 {
			performanceSlice = append(performanceSlice, v)

		}
	}
	if mode == "CLEAR" {
		data.Title = "Clear Difficulties"
		sort.Slice(performanceSlice, func(i int, j int) bool { return performanceSlice[i].ClearRate > performanceSlice[j].ClearRate })
	} else if mode == "FC" {
		data.Title = "FC Difficulties"
		sort.Slice(performanceSlice, func(i int, j int) bool { return performanceSlice[i].FullComboRate > performanceSlice[j].FullComboRate })
	} else if mode == "AP" {
		data.Title = "AP Difficulties"
		sort.Slice(performanceSlice, func(i int, j int) bool {
			return performanceSlice[i].AllPerfectRate > performanceSlice[j].AllPerfectRate
		})
	} else if mode == "GENERAL" {
		data.Title = "General Difficulties"
		sort.Slice(performanceSlice, func(i int, j int) bool { return performanceSlice[i].GeneralRate > performanceSlice[j].GeneralRate })
	} else {
		err := PageNotFound
		err.Detail = "ページが見つかりませんでした。urlが正しいかもう一度お確かめください。"
		c.Error(&err).SetType(gin.ErrorTypePublic).SetMeta("")
		return
	}
	plates := []Plate{}
	for _, v := range performanceSlice {
		song := songs[songIndexes[v.SongId]]
		namePlate := song.Name
		ratePlate := ""
		if mode == "CLEAR" {
			ratePlate = fmt.Sprintf("%.2f", v.ClearRate/1.25+5)
		} else if mode == "FC" {
			ratePlate = fmt.Sprintf("%.2f", v.FullComboRate/1.25+5)
		} else if mode == "AP" {
			ratePlate = fmt.Sprintf("%.2f", v.AllPerfectRate/1.25+5)
		} else {
			ratePlate = fmt.Sprintf("%.2f", v.GeneralRate/1.25+5)
		}
		/*if v.difficulty == 4 {
			fmt.Fprintf(w, "%s", namePlate)
		} else if v.difficulty == 3 {
			fmt.Fprintf(w, "%s", namePlate)
		} else {
			fmt.Fprintf(w, "%s", namePlate)
		}*/
		plate := Plate{}
		plate.Rate = ratePlate
		plate.DifficultyName = difficultyname[song.Difficulty]
		plate.Level = song.Level
		plate.Name = namePlate
		plate.Try = v.Try
		plates = append(plates, plate)
		// fmt.Fprintf(w, "clear: %6f,\tfc: %6f\n\n", v.clearRate/1.25+5, v.fullComboRate/1.25+5)
	}
	data.Data = gin.H{"plates": plates}
	c.HTML(http.StatusOK, "show.htm", data)
}

func update(performance models.Performance, play models.Play) models.Performance {
	clearScore := play.ToClearScore()
	fullComboScore := play.ToFullComboScore()
	allPerfectScore := play.ToAllPerfectScore()
	performance.Try++
	// fmt.Printf("clearScore: %f\n", clearScore)
	// fmt.Printf("fullComboScore: %f\n", fullComboScore)
	mult := 1.0
	if performance.Try == 1 {
		mult = 20
		/*} else if song.try == 2 {
			mult = 11
		} else if song.try == 3 {
			mult = 9
		} else if song.try == 4 {
			mult = 7
		} else if song.try == 5 {
			mult = 6
		} else if song.try == 6 {
			mult = 5*/
	} else if performance.Try <= 10 {
		mult = 20 / float64(performance.Try)
	} else if performance.Try <= 100 {
		mult = 20 / math.Sqrt(float64(10*performance.Try))
	}
	newClearRate := mult*VOLATILITY_CLEAR*(clearScore/-200) + (1-mult*VOLATILITY_CLEAR)*performance.ClearRate
	newFullComboRate := mult*VOLATILITY_FULLCOMBO*(fullComboScore/-200) + (1-mult*VOLATILITY_FULLCOMBO)*performance.FullComboRate
	newAllPerfectRate := mult*VOLATILITY_ALLPERFECT*(allPerfectScore/-200) + (1-mult*VOLATILITY_ALLPERFECT)*performance.AllPerfectRate
	// fmt.Printf("clear%f -> %f\n", performance.clearRate, newClearRate)
	// fmt.Printf("fc%f -> %f\n", performance.fullComboRate, newFullComboRate)
	performance.ClearRate = newClearRate
	performance.FullComboRate = newFullComboRate
	performance.AllPerfectRate = newAllPerfectRate
	performance.GeneralRate = (newClearRate + newFullComboRate + newAllPerfectRate) / 3

	return performance
}
