package nfteventdeal

import (
	"context"
	"dapp_nft/config"
	"dapp_nft/models"
	"fmt"
	"log"
	"math"
	"math/big"
	"os"
	"strings"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
)

func syncHisotry() {
	rpcURL := os.Getenv("RPC_URL_HTTP")
	if rpcURL == "" {
		log.Fatal("RPC_URL_HTTP must be set")
	}
	contractAddr := os.Getenv("NFT_MARKET_ADDRESS")
	if contractAddr == "" {
		log.Fatal("missing --NFT_MARKET_ADDRESS")
	}
	var latestBlockNumber_auction uint
	var latestBlockNumber_bid uint
	var latestBlockNumber_nftowner uint
	config.DB.Model(&models.Auction{}).Select("max(block_number)").Find(&latestBlockNumber_auction)
	config.DB.Model(&models.Bid{}).Select("max(block_number)").Find(&latestBlockNumber_bid)
	config.DB.Model(&models.NftOwner{}).Select("max(block_number)").Find(&latestBlockNumber_nftowner)
	maxBlockNumber := math.Max(math.Max(float64(latestBlockNumber_auction), float64(latestBlockNumber_bid)), float64(latestBlockNumber_nftowner))

	query := ethereum.FilterQuery{
		FromBlock: big.NewInt(int64(maxBlockNumber)),
		// ToBlock:   big.NewInt(2394201),
		Addresses: []common.Address{
			common.HexToAddress(contractAddr),
		},
		// Topics: [][]common.Hash{
		//  {},
		//  {},
		// },
	}
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	client, err := ethclient.DialContext(ctx, rpcURL)
	if err != nil {
		log.Fatalf("failed to connect to Ethereum node: %v", err)
	}
	defer client.Close()

	logs, err := client.FilterLogs(context.Background(), query)
	if err != nil {
		log.Fatal(err)
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

	for _, vLog := range logs {
		ParseLogEvent(&vLog, parsedABI)
	}

}
