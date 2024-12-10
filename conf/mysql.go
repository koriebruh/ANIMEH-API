package conf

import (
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"koriebruh/find/domain"
	"log/slog"
)

func InitDB(conf *Config) *gorm.DB {

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		conf.Mysql.User, // Username
		conf.Mysql.Pass, // Password
		conf.Mysql.Host, // Hostname
		conf.Mysql.Port, // Port
		conf.Mysql.Name, // Database Name
	)

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})

	if err != nil {
		slog.Error("failed make connection to database", err)
	}

	if err = db.AutoMigrate(
		&domain.User{},
		&domain.Favorite{},
	); err != nil {
		slog.Error("failed auto migrate ", err)
	}

	slog.Info("success migrate")
	return db
}
