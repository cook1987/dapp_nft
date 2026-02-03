package models

import (
	"time"
)

type Auction struct {
	ID           uint      `json:"id" gorm:"primaryKey"`
	AuctionID    uint      `json:"auction_id" gorm:"not null;comment:拍卖id"`
	Seller       string    `json:"seller" gorm:"not null;size:42"`
	TokenAddress string    `json:"token_address" gorm:"not null;size:42"`
	TokenId      uint      `json:"token_id" gorm:"not null"`
	TokenUrl     string    `json:"token_url" gorm:"not null"`
	Price        uint      `json:"price" gorm:"not null;comment:成交价格（美元）"`
	StartPrice   uint      `json:"start_price" gorm:"not null;comment:起拍价价格（美元）"`
	Status       uint      `json:"status" gorm:"default:1;comment:状态 1:上架 2:成交 3:流拍"`
	BlockNumber  uint      `json:"block_number" gorm:"not null;comment:区块号"`
	CreateTxHash string    `json:"create_tx_hash" gorm:"not null;uniqueIndex:idx_auction_ctx_hash;comment:拍卖创建交易hash;size:66"`
	EndTxHash    string    `json:"end_tx_hash" gorm:"index:idx_auction_etx_hash;comment:拍卖结束交易hash;size:66"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

const (
	AuctionState_on       = iota + 1 // 1：上架
	AuctionState_deal                // 2:成交
	AuctionState_abortive            // 3:流拍
)
