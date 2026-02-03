package models

import (
	"time"
)

type NftOwner struct {
	ID           uint      `json:"id" gorm:"primaryKey"`
	AuctionInID  uint      `json:"auction_in_id" gorm:"comment:买入时的拍卖id"`
	AuctionOutID uint      `json:"auction_out_id" gorm:"comment:卖出时的拍卖id"`
	Owner        string    `json:"owner" gorm:"not null;size:42"`
	TokenAddress string    `json:"token_address" gorm:"not null;size:42"`
	TokenId      uint      `json:"token_id" gorm:"not null"`
	TokenUrl     string    `json:"token_url" gorm:"not null"`
	Status       uint      `json:"status" gorm:"default:1;comment:状态 1:拥有 2:卖掉"`
	BlockNumber  uint      `json:"block_number" gorm:"not null;comment:买入时区块号"`
	TxHash       string    `json:"tx_hash" gorm:"not null;uniqueIndex:idx_nftowner_tx_hash;comment:买入时交易hash;size:66"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

const (
	NftOwnerStatus_have = iota + 1 // 1：拥有
	NftOwnerStatus_sold            // 2: 卖掉
)
