package controller

import (
	"log"
	"net/http"

	"github.com/MONTplusa/ProjectSekaiDifficultyCalculation/config"
	"github.com/MONTplusa/ProjectSekaiDifficultyCalculation/models"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
)

func GetRouter() *gin.Engine {
	engine := gin.Default()
	store := cookie.NewStore([]byte("secret"))
	engine.Use(sessions.Sessions("mysession", store))
	engine.LoadHTMLGlob("templates/*")
	engine.Use(errHandler())
	engine.GET("/", home)
	adminCheckGroup := engine.Group("/", checkAdmin())
	adminCheckGroup.GET("/form_new_song", form_NewSong)
	adminCheckGroup.POST("/add_song", addSong)
	adminCheckGroup.GET("/import", importData)
	adminCheckGroup.GET("/import_song", importSongData)
	adminCheckGroup.GET("/form_create_user", form_createUser)
	adminCheckGroup.POST("/create_user", createUser)
	loginCheckGroup := engine.Group("/", checkLogin())
	loginCheckGroup.GET("/form", form)
	loginCheckGroup.POST("/add", add)
	loginCheckGroup.GET("/show/:mode", show)

	engine.GET("/login", form_login)
	engine.POST("/login", login)
	return engine
}

func errHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()
		err := c.Errors.ByType(gin.ErrorTypePublic).Last()
		if err != nil {
			meta := c.Errors.ByType(gin.ErrorTypePublic).Last().Meta
			log.Print(err.Err)
			log.Print(meta)
			httperror := err.Err.(*HttpError)
			data := NewBaseData(c)
			data.Data = gin.H{"error": httperror}
			if c.Request.Method == "GET" {
				c.HTML(httperror.StatusCode, "error.htm", data)
			} else if c.Request.Method == "POST" {
				c.JSON(http.StatusBadRequest, gin.H{"meta": meta})
			} else {
				c.HTML(http.StatusInternalServerError, "error.htm", data)
			}
			c.Abort()
		}
	}
}

func checkAdmin() gin.HandlerFunc {
	return func(c *gin.Context) {
		session := sessions.Default(c)
		id := session.Get("UserId").(int)
		var user models.User
		models.GetUserById(&user, id)
		if user.Username != config.Config.AdminName {
			c.Redirect(http.StatusFound, "/login")
			c.Abort()
		} else {
			c.Next()
		}
	}
}
func checkLogin() gin.HandlerFunc {
	return func(c *gin.Context) {
		session := sessions.Default(c)
		id := session.Get("UserId")
		if id == nil {
			c.Redirect(http.StatusFound, "/login")
			c.Abort()
		} else {
			c.Next()
		}
	}
}
