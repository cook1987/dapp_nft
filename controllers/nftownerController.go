package controllers

import (
	"dapp_nft/config"
	"dapp_nft/models"
	"dapp_nft/utils"
	"net/http"

	"github.com/gin-gonic/gin"
)

type NftownerController struct{}

// 查询某个钱包地址拥有的所有NFT Token列表
func (ac *NftownerController) GetNftOfOwner(c *gin.Context) {
	address := c.Query("address")
	if address == "" {
		utils.BadRequest(c, "Invalid address")
		return
	}

	var nftOwners []models.NftOwner

	// 获取数据
	if err := config.DB.Where("owner = ?", address).Find(&nftOwners).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	utils.Success(c, &nftOwners)

}
