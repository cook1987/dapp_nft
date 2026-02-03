package nfteventdeal

import (
	"context"
	"dapp_nft/config"
	"dapp_nft/models"
	"fmt"
	"log"
	"math/big"
	"os"
	"strings"
	"time"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
)

func SubscribeNftMarket() {
	rpcURL := os.Getenv("RPC_URL_WSS")
	if rpcURL == "" {
		log.Fatal("RPC_URL_WSS must be set")
	}
	contractAddr := os.Getenv("NFT_MARKET_ADDRESS")
	if contractAddr == "" {
		log.Fatal("missing --NFT_MARKET_ADDRESS")
	}

	content, err := os.ReadFile("CookNFTMarketplace.json")
	if err != nil {
		panic(fmt.Sprintf("读取拍卖市场 ABI 文件失败：%v", err))
	}
	// 解析 ABI
	parsedABI, err := abi.JSON(strings.NewReader(string(content)))
	if err != nil {
		log.Fatalf("failed to parse ABI: %v", err)
	}

	ctx, _ := context.WithCancel(context.Background())
	// defer cancel()

	client, err := ethclient.DialContext(ctx, rpcURL)
	if err != nil {
		log.Fatalf("failed to connect to Ethereum node: %v", err)
	}
	// defer client.Close()

	contract := common.HexToAddress(contractAddr)

	query := ethereum.FilterQuery{
		Addresses: []common.Address{contract},
	}

	logsCh := make(chan types.Log)
	sub, err := client.SubscribeFilterLogs(ctx, query, logsCh)
	if err != nil {
		log.Fatalf("failed to subscribe logs: %v", err)
	}
	fmt.Printf("Subscribed to logs of contract %s via %s\n", contract.Hex(), rpcURL)
	fmt.Printf("Listening for events...\n\n")

	// 同步历史记录
	go syncHisotry()

	go func() {
		for {
			select {
			case vLog := <-logsCh:
				// 解析日志事件
				ParseLogEvent(&vLog, parsedABI)
			case err := <-sub.Err():
				log.Printf("subscription error: %v", err)
				client.Close()
				go SubscribeNftMarket()
				return
			case <-ctx.Done():
				fmt.Println("context cancelled, exiting...")
				return
			}
		}
	}()
}

// parseLogEvent 解析日志事件，展示如何从 logs 中提取事件信息
func ParseLogEvent(vLog *types.Log, parsedABI abi.ABI) {
	if len(vLog.Topics) == 0 {
		return
	}

	// 步骤 1: 识别事件类型
	// Topics[0] 是事件签名的 keccak256 哈希值
	// 例如: Transfer(address,address,uint256) 的哈希
	eventTopic := vLog.Topics[0]

	// 尝试识别是哪个事件（通过比较 Topics[0] 和事件签名的哈希）
	var eventName string

	// 遍历 ABI 中定义的所有事件，查找匹配的事件签名
	for name, event := range parsedABI.Events {
		// 计算事件的签名哈希
		eventSigHash := crypto.Keccak256Hash([]byte(event.Sig))
		if eventSigHash == eventTopic {
			eventName = name
			break
		}
	}

	if eventName == "" {
		// 如果无法识别事件类型，打印原始信息
		fmt.Printf("[%s] Unknown Event - Block: %d, Tx: %s, Topic[0]: %s\n",
			time.Now().Format(time.RFC3339),
			vLog.BlockNumber,
			vLog.TxHash.Hex(),
			eventTopic.Hex(),
		)
		return
	}

	switch eventName {
	case "AuctionCreated":
		log.Printf("监听到事件 %s 开始处理\n", eventName)
		auctionCreatedEvent := struct {
			AuctionId   *big.Int       `json:"auctionId"`
			Seller      common.Address `json:"seller"`
			NftContract common.Address `json:"nftContract"`
			TokenId     *big.Int       `json:"tokenId"`
			StartPrice  *big.Int       `json:"startPrice"`
			EndTime     *big.Int       `json:"endTime"`
		}{}
		err := parsedABI.UnpackIntoInterface(&auctionCreatedEvent, "AuctionCreated", vLog.Data)
		if err != nil {
			log.Fatal(err)
		}

		auction := models.Auction{
			AuctionID:    uint(new(big.Int).SetBytes(vLog.Topics[1].Bytes()).Uint64()),
			Seller:       common.BytesToAddress(vLog.Topics[2].Bytes()).String(),
			TokenAddress: common.BytesToAddress(vLog.Topics[3].Bytes()).String(),
			TokenId:      uint(auctionCreatedEvent.TokenId.Uint64()),
			StartPrice:   uint(auctionCreatedEvent.StartPrice.Uint64()),
			BlockNumber:  uint(vLog.BlockNumber),
			CreateTxHash: vLog.TxHash.Hex(),
		}
		log.Println("create new Auction: ", auction)
		// todo: 获取 tokenurl
		if err := config.DB.Create(&auction).Error; err != nil {
			log.Printf("failed to create Auction: %v", err)
		}
	case "BidPlaced":
		log.Printf("监听到事件 %s 开始处理\n", eventName)
		bidEvent := struct {
			AuctionId  *big.Int       `json:"auctionId"`
			Bidder     common.Address `json:"bidder"`
			Erc20Token common.Address `json:"erc20Token"`
			Amount     *big.Int       `json:"amount"`
		}{}
		err := parsedABI.UnpackIntoInterface(&bidEvent, "BidPlaced", vLog.Data)
		if err != nil {
			log.Fatal(err)
		}
		bid := models.Bid{
			AuctionID:   uint(new(big.Int).SetBytes(vLog.Topics[1].Bytes()).Uint64()),
			Bidder:      common.BytesToAddress(vLog.Topics[2].Bytes()).String(),
			Erc20Token:  common.BytesToAddress(bidEvent.Erc20Token.Bytes()).String(),
			Price:       uint(bidEvent.Amount.Uint64()),
			BlockNumber: uint(vLog.BlockNumber),
			TxHash:      vLog.TxHash.Hex(),
		}
		log.Println("create new bid: ", bid)
		// todo: 获取 tokenurl
		if err := config.DB.Create(&bid).Error; err != nil {
			log.Printf("failed to create Bid: %v", err)
		}
	case "AuctionEnded":
		log.Printf("监听到事件 %s 开始处理\n", eventName)
		auctionEndedEvent := struct {
			AuctionId  *big.Int       `json:"auctionId"`
			Winner     common.Address `json:"winner"`
			Erc20Token common.Address `json:"erc20Token"`
			FinalPrice *big.Int       `json:"finalPrice"`
		}{}
		err := parsedABI.UnpackIntoInterface(&auctionEndedEvent, "AuctionEnded", vLog.Data)
		if err != nil {
			log.Fatal(err)
		}
		price := auctionEndedEvent.FinalPrice.Uint64()
		auctionID := uint(new(big.Int).SetBytes(vLog.Topics[1].Bytes()).Uint64())
		bidder := common.BytesToAddress(vLog.Topics[2].Bytes()).String()
		var auction models.Auction
		config.DB.Where("auction_id = ?", auctionID).Find(&auction)
		if auction.ID == 0 {
			log.Printf("auctionID: %d not exists!\n", auctionID)
			return
		}
		if price == 0 {
			// 没有人出价
			auction.Status = models.AuctionState_abortive
			auction.EndTxHash = vLog.TxHash.Hex()
			if err := config.DB.Save(&auction).Error; err != nil {
				log.Printf("failed to update auction: %v", err)
			}
		} else {
			auction.Status = models.AuctionState_deal
			auction.Price = uint(price)
			auction.EndTxHash = vLog.TxHash.Hex()

			if err := config.DB.Save(&auction).Error; err != nil {
				log.Printf("failed to update auction: %v", err)
			}

			var oldOwner models.NftOwner
			config.DB.Where("token_address = ?", auction.TokenAddress).
				Where("token_id = ?", auction.TokenId).
				Where("status = ?", models.NftOwnerStatus_have).Find(&oldOwner)
			if oldOwner.ID > 0 {
				oldOwner.AuctionOutID = auctionID
				oldOwner.Status = models.NftOwnerStatus_sold
				config.DB.Save(&oldOwner)
			}
			newOwner := models.NftOwner{
				AuctionInID:  auctionID,
				Owner:        bidder,
				TokenAddress: auction.TokenAddress,
				TokenId:      auction.TokenId,
				TokenUrl:     auction.TokenUrl,
				BlockNumber:  uint(vLog.BlockNumber),
				TxHash:       vLog.TxHash.Hex(),
			}

			if err := config.DB.Save(&newOwner).Error; err != nil {
				log.Printf("failed to save NftOwner: %v", err)
			}
		}
	default:
		log.Printf("监听到事件 %s 忽略\n", eventName)
	}
}
