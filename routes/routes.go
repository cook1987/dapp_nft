package routes

import (
	"dapp_nft/controllers"
	"dapp_nft/middleware"

	"github.com/gin-gonic/gin"
)

// 设置路由
func SetupRoutes() *gin.Engine {
	r := gin.New()

	// 使用中间件
	r.Use(middleware.LoggerMiddleware())
	r.Use(middleware.ErrorHandlerMiddleware())
	r.Use(gin.Recovery())

	// 创建控制器实例
	acutionController := &controllers.AuctionController{}
	bidController := &controllers.BidController{}
	nftOwnerController := &controllers.NftownerController{}

	// api 路由组
	api := r.Group("/api/v1")
	api.GET("/auctionPage", acutionController.GetAuctionPage)
	api.GET("/bidListOfAuction", bidController.GetBidListOfAuction)
	api.GET("/getBidStatic", bidController.GetBidStatic)
	api.GET("/getNftOfOwner", nftOwnerController.GetNftOfOwner)

	// 健康检查
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status":  "ok",
			"message": "Blog API is running",
		})
	})

	return r
}
