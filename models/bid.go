package models

import (
	"time"
)

type Bid struct {
	ID          uint      `json:"id" gorm:"primaryKey"`
	AuctionID   uint      `json:"auction_id" gorm:"not null;comment:拍卖id"`
	Bidder      string    `json:"bidder" gorm:"not null;size:42"`
	Price       uint      `json:"price" gorm:"not null;comment:价格（美元）"`
	Erc20Token  string    `json:"erc20_token" gorm:"size:42"`
	BlockNumber uint      `json:"block_number" gorm:"not null;comment:区块号"`
	TxHash      string    `json:"tx_hash" gorm:"not null;uniqueIndex:idx_bid_tx_hash;size:66"`
	CreatedAt   time.Time `json:"created_at"`
}
