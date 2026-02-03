package controllers

import (
	"dapp_nft/config"
	"dapp_nft/models"
	"dapp_nft/utils"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm/clause"
)

type AuctionController struct{}

type AuctionPageRequest struct {
	Page         uint   `json:"page"`
	PageSize     uint   `json:"pageSize"`
	TokenAddress string `json:"tokenAddress"`
	Status       uint   `json:"status" binding:"min=0,max=3"`
	CreatedSort  uint   `json:"createSort" binding:"min=0,max=1"`
}

func NewAuctionPageRequestBuilder() *AuctionPageRequest {
	return &AuctionPageRequest{
		Page:        1,
		PageSize:    20,
		CreatedSort: 1,
	}
}

// 拍卖列表接口，支持排序和过滤条件
func (ac *AuctionController) GetAuctionPage(c *gin.Context) {
	// 分页参数
	req := NewAuctionPageRequestBuilder()
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, err.Error())
		return
	}
	if req.Page < 1 {
		req.Page = 1
	}
	if req.PageSize < 1 || req.PageSize > 100 {
		req.PageSize = 10
	}
	offset := (req.Page - 1) * req.PageSize

	var auctions []models.Auction

	db := config.GetDB()

	if db == nil {
		logrus.Error("database not available")
	}

	query := db.Model(&models.Auction{})

	// 动态添加查询条件
	if req.TokenAddress != "" {
		query = query.Where("token_address = ?", req.TokenAddress)
	}

	if req.Status > 0 {
		query = query.Where("status = ?", req.Status)
	}

	var total int64

	// 获取总数
	if err := query.Count(&total).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// 获取数据
	if err := query.
		Order(clause.OrderByColumn{Column: clause.Column{Name: "created_at"}, Desc: (req.CreatedSort == 1)}).
		Limit(int(req.PageSize)).
		Offset(int(offset)).
		Find(&auctions).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	utils.Success(c, gin.H{
		"data": auctions,
		"pagination": gin.H{
			"page":      req.Page,
			"page_size": req.PageSize,
			"total":     total,
			"pages":     (total + int64(req.PageSize) - 1) / int64(req.PageSize),
		},
	})

}
