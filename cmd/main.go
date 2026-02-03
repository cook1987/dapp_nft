package main

import (
	"dapp_nft/config"
	subscribe "dapp_nft/nfteventdeal"
	"dapp_nft/routes"
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
)

func main() {
	// 加载环境变量
	if err := godotenv.Load(); err != nil {
		logrus.Warn("No .env file found, using system environment variables")
	}

	// 初始化日志
	logrus.SetFormatter(&logrus.JSONFormatter{})
	logrus.Info("Starting nftMarket application...")

	// 初始化数据库
	config.InitDatabase()

	// 设置路由
	r := routes.SetupRoutes()

	// 启动监听
	subscribe.SubscribeNftMarket()

	// 启动服务器
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	port = ":" + port

	logrus.WithField("port", port).Info("Server starting")
	if err := r.Run(port); err != nil {
		log.Fatal("Failed to start server: ", err)
	}

}
