package main

import (
	"strconv"
	"time"

	"github.com/MONTplusa/ProjectSekaiDifficultyCalculation/config"
	"github.com/MONTplusa/ProjectSekaiDifficultyCalculation/controller"
	"github.com/MONTplusa/ProjectSekaiDifficultyCalculation/models"
	"golang.org/x/crypto/bcrypt"
)

func main() {
	r := controller.GetRouter()
	admin := config.Config.AdminName
	password := config.Config.AdminPassword
	hashed, _ := bcrypt.GenerateFromPassword([]byte(password), 10)
	var user models.User
	user.LastAccess = time.Now()
	user.Password = string(hashed)
	models.GetUserByName(&user, admin)
	if user.Username == "" {
		user.Username = admin
		models.InsertUser(&user)
	} else {
		models.UpdateUser(&user)

	}
	port := strconv.Itoa(config.Config.ServerPort)
	r.Run(":" + port)
}
