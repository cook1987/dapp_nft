package controllers

import (
	"dapp_nft/config"
	"dapp_nft/models"
	"dapp_nft/utils"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type BidController struct{}

// 查询某个拍卖的出价历史记录
func (ac *BidController) GetBidListOfAuction(c *gin.Context) {
	auctionId, err := strconv.ParseUint(c.Query("auctionId"), 10, 32)
	if err != nil {
		utils.BadRequest(c, "Invalid Auction ID")
		return
	}

	// 检查拍卖是否存在
	var auction models.Auction
	if err := config.DB.Where("auction_id = ?", auctionId).Find(&auction).Error; err != nil {
		utils.NotFound(c, "Auction not found")
		return
	}

	var bids []models.Bid

	// 获取数据
	if err := config.DB.Where("auction_id = ?", auctionId).Find(&bids).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	utils.Success(c, &bids)
}

// 平台的统计数据（拍卖总数，出价总数）
func (ac *BidController) GetBidStatic(c *gin.Context) {

	var auctionCount int64
	var bidCount int64
	// 获取总数
	if err := config.DB.Model(&models.Auction{}).Count(&auctionCount).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if err := config.DB.Model(&models.Bid{}).Count(&bidCount).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	utils.Success(c, gin.H{
		"auctionCount": auctionCount,
		"bidCount":     bidCount,
	})
}
