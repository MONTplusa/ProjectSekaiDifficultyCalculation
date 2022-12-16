package config

import (
	"log"
	"os"

	"gopkg.in/ini.v1"
)

type ConfigList struct {
	ConfigDB
	ConfigAPI
	ConfigWebSecurity
}

type ConfigDB struct {
	DbDriverName   string
	DbName         string
	DbUserName     string
	DbUserPassword string
	DbHost         string
	DbPort         string
}
type ConfigAPI struct {
	ServerPort int
}
type ConfigWebSecurity struct {
	AdminName       string
	AdminPassword   string
	CookieSecretKey string
}

var Config ConfigList

func init() {
	cfg, err := ini.Load("config.ini")
	if err != nil {
		log.Printf("Failed to read file: %v", err)
		os.Exit(1)
	}

	Config = ConfigList{}
	Config.DbDriverName = cfg.Section("db").Key("db_driver_name").String()
	Config.DbName = cfg.Section("db").Key("db_name").String()
	Config.DbUserName = cfg.Section("db").Key("db_user_name").String()
	Config.DbUserPassword = cfg.Section("db").Key("db_user_password").String()
	Config.DbHost = cfg.Section("db").Key("db_host").String()
	Config.DbPort = cfg.Section("db").Key("db_port").String()
	Config.ServerPort = cfg.Section("api").Key("server_port").MustInt()
	Config.AdminName = cfg.Section("security").Key("admin_name").String()
	Config.AdminPassword = cfg.Section("security").Key("admin_password").String()
	Config.CookieSecretKey = cfg.Section("security").Key("cookie_secret_key").String()
}
