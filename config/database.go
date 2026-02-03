package config

import (
	"dapp_nft/models"
	"fmt"
	"log"
	"os"

	"github.com/sirupsen/logrus"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

// 初始化数据库链接
func InitDatabase() {
	var err error

	// 从环境变量获取 MySql 连接配置
	dbHost := getEnv("DB_HOST", "localhost")
	dbPort := getEnv("DB_PORT", "3306")
	dbUser := getEnv("DB_USER", "root")
	dbPassword := getEnv("DB_PASSWORD", "")
	dbName := getEnv("DB_NAME", "golang_nft")

	// 构造mysql 连接字符串
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		dbUser, dbPassword, dbHost, dbPort, dbName)

	// 连接mysql 数据库
	DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		log.Fatal("Failed to connect to Mysql database: ", err)
	}

	// 自动迁移数据库表结构
	err = DB.AutoMigrate(
		&models.Auction{},
		&models.Bid{},
		&models.NftOwner{},
	)

	if err != nil {
		log.Fatal("Failed to migrate database: ", err)
	}

	log.Println("MySQL database connected and migrated successfully")

	// 检查表是否存在
	if DB.Migrator().HasTable(&models.Auction{}) {
		logrus.Info("✓ Auction table created successfully")
	}
	if DB.Migrator().HasTable(&models.Bid{}) {
		logrus.Info("✓ Bid table created successfully")
	}
	if DB.Migrator().HasTable(&models.NftOwner{}) {
		logrus.Info("✓ NftOwner table created successfully")
	}
}

// 获取数据库连接实例
func GetDB() *gorm.DB {
	if DB == nil {
		log.Fatal("Database not initialized. Call InitDB first.")
	}
	return DB
}

// 获取环境变量，如果不存在则返回默认值
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
