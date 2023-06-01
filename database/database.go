package database

import (
	"fmt"

	"github.com/Base-Technology/base-backend-lite/conf"
	"github.com/pkg/errors"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var db *gorm.DB

func InitDatabase() error {
	var err error
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/base-lite?charset=utf8mb4&parseTime=True&loc=Local",
		conf.Conf.DBConf.Username,
		conf.Conf.DBConf.Password,
		conf.Conf.DBConf.IP,
		conf.Conf.DBConf.Port,
	)
	db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		return errors.Errorf("open database error, %v", err)
	}

	db.AutoMigrate(&User{})
	db.AutoMigrate(&Post{})
	db.AutoMigrate(&Image{})
	db.AutoMigrate(&Comment{})
	db.AutoMigrate(&Follow{})
	db.AutoMigrate(&Like{})
	db.AutoMigrate(&Collect{})
	db.AutoMigrate(&FriendRequest{})
	return nil
}

func GetInstance() *gorm.DB {
	return db
}
