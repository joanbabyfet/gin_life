package models

import (
	"fmt"

	"github.com/jinzhu/gorm"
	"github.com/lexkong/log"
	"github.com/spf13/viper"
)

type Database struct {
	Self *gorm.DB
}

var (
	DB *Database
)

func (db *Database) Init() {
	DB = &Database{
		Self: GetSelfDB(),
	}
}

func openDB(username, password, addr, name, charset string) *gorm.DB {
	config := fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=%s&parseTime=%t&loc=%s",
		username,
		password,
		addr,
		name,
		charset,
		true,
		"Local")

	db, err := gorm.Open("mysql", config)

	if err != nil {
		log.Errorf(err, "Database connection failed. Database name: %s", name)
	}
	setupDB(db)
	return db
}

func setupDB(db *gorm.DB) {
	db.LogMode(viper.GetBool("gormlog"))
	//db.LogMode(true)
	//db.DB().SetMaxOpenConns(20000) // 用于设置最大打开的连接数，默认值为0表示不限制.设置最大的连接数，可以避免并发太高导致连接mysql出现too many connections的错误。
	db.DB().SetMaxIdleConns(0) // 用于设置闲置的连接数.设置闲置的连接数则当开启的一个连接使用完成后可以放在池里等候下一次使用。
}

func InitSelfDB() *gorm.DB {
	return openDB(
		viper.GetString("db.username"),
		viper.GetString("db.password"),
		viper.GetString("db.addr"),
		viper.GetString("db.name"),
		viper.GetString("db.charset"))
}

func GetSelfDB() *gorm.DB {
	return InitSelfDB()
}

func (db *Database) Close() {
	DB.Self.Close()
}
