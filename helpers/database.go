package helpers

import (
	"ewallet-wallet/internal/models"
	"fmt"

	"github.com/sirupsen/logrus"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB

func SetupMySQL() {
	dbuser := GetEnv("DB_USER", "")
	dbpass := GetEnv("DB_PASSWORD", "")
	dbhost := GetEnv("DB_HOST", "127.0.0.1")
	dbport := GetEnv("DB_PORT", "3036")
	dbname := GetEnv("DB_NAME", "")

	createDBDsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/", dbuser, dbpass, dbhost, dbport)
	database, err := gorm.Open(mysql.Open(createDBDsn), &gorm.Config{})
	if err != nil {
		logrus.Fatal("failed to create database", err)
	}

	result := database.Exec("CREATE DATABASE IF NOT EXISTS " + dbname + ";")

	if result.Error != nil {
		logrus.Fatal("failed create database")
	}

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=true&loc=Local",
		dbuser,
		dbpass,
		dbhost,
		dbport,
		dbname,
	)

	DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		logrus.Fatal("failed to connect to database", err)
	}

	logrus.Info("successfully connect to database")

	DB.AutoMigrate(&models.Wallet{}, &models.WalletTransaction{}, &models.WalletLink{})
}
