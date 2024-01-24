package mysql

import (
	"fmt"
	"time"

	"go_backend/model"

	"github.com/spf13/viper"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
)

func InitMySQL() (*gorm.DB, error) {
	dsn := fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=%s&parseTime=%s&loc=%s",
		viper.GetString("mysql.username"),
		viper.GetString("mysql.password"),
		viper.GetString("mysql.host"),
		viper.GetString("mysql.dbName"),
		viper.GetString("mysql.charset"),
		viper.GetString("mysql.parseTime"),
		viper.GetString("mysql.loc"))

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			SingularTable: viper.GetBool("mysql.singularTable"),
			TablePrefix:   viper.GetString("mysql.tablePrefix"),
		},
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		return nil, err
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, err
	}

	sqlDB.SetMaxIdleConns(viper.GetInt("mysql.maxIdleConns"))
	sqlDB.SetMaxOpenConns(viper.GetInt("mysql.maxOpenConns"))
	sqlDB.SetConnMaxLifetime(time.Duration(viper.GetInt("mysql.maxLifetime")) * time.Hour)

	if err := db.AutoMigrate(&model.User{}); err != nil {
		return nil, err
	}

	return db, nil
}
