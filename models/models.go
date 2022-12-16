package models

import (
	"fmt"
	"log"

	"github.com/MONTplusa/ProjectSekaiDifficultyCalculation/config"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

var Db *gorm.DB

func init() {
	var err error
	dbConnectInfo := fmt.Sprintf(
		`%s:%s@tcp(%s:%s)/%s?charset=utf8&parseTime=True&loc=Local`,
		config.Config.DbUserName,
		config.Config.DbUserPassword,
		config.Config.DbHost,
		config.Config.DbPort,
		config.Config.DbName,
	)

	// DBに接続
	Db, err = gorm.Open(config.Config.DbDriverName, dbConnectInfo)
	if err != nil {
		log.Fatalln(err)
	}

	// テーブルをAutoMigrate
	Db.Set("gorm:table_options", "ENGINE = InnoDB").AutoMigrate(&Play{})
	Db.Set("gorm:table_options", "ENGINE = InnoDB").AutoMigrate(&Song{})
	Db.Set("gorm:table_options", "ENGINE = InnoDB").AutoMigrate(&User{})
	Db.Set("gorm:table_options", "ENGINE = InnoDB").AutoMigrate(&Achievement{})
}
