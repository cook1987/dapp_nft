## 一、项目运行环境：<br/>
go版本：go1.25.4<br/>
mysql版本: 8.0<br/>
## 二、依赖安装步骤
gorm: go get -u gorm.io/gorm<br/>
mysql: go get -u gorm.io/driver/mysql<br/>
gin: go get -u github.com/gin-gonic/gin<br/>
日志：go get github.com/sirupsen/logrus<br/>
以太坊：go get github.com/ethereum/go-ethereum<br/>

## 三、文件列表
- cmd
- - main.go  --项目启动文件
- config
- - database.go  --数据库配置及数据表初始化文件
- controllers
- - auctionController.go  --拍卖相关接口
- - bidController.go  --出价接口
- - nftownerController.go  --nft所有者相关接口
- middleware
- - logger.go
- models
- - auction.go
- - big.go
- - nftowner.go
- nfteventdeal
- - nftsuscribe.go  --监听链上事件方法
- - synchistory.go  --项目启动时同步遗漏历史记录方法
- routes
- - routes.go  --路由配置
- utils
- - response.go


### 异常处理代码位置：
#### 1、网络链接断掉：nftsuscribe.go  第 74 行：创建goroutine重新连接网络
```shell
	case err := <-sub.Err():
		log.Printf("subscription error: %v", err)
		client.Close()
		go SubscribeNftMarket()
```

#### 1、服务器重启：nftsuscribe.go  第 66 行：同步历史记录（查询出表中最大的 block_number，从此处开始同步之后的所有事件）
```shell
go syncHisotry()
```

## 四、启动方式
homework04> go run .\cmd\main.go
## 五、接口列表
### 1. 拍卖列表接口：GET http://localhost:8080/api/v1/auctionPage
BODY:<br/>

{
    "page": 1,
    "pageSize": 20
}
<br/>
响应：<br/>
{
	"code": 200,
	"message": "success",
	"data": {
		"data": [
			{
				"id": 3,
				"auction_id": 8,
				"seller": "0x256e72b29D48F1F98792EF1fF854c36043be1b6B",
				"token_address": "0x38D2c9D2425A13Be961ED9eD1C45834B8195FbFB",
				"token_id": 3,
				"token_url": "",
				"price": 0,
				"start_price": 3,
				"status": 1,
				"created_at": "2026-02-03T18:11:12.929+08:00",
				"updated_at": "2026-02-03T18:11:12.929+08:00"
			},
			{
				"id": 2,
				"auction_id": 7,
				"seller": "0x256e72b29D48F1F98792EF1fF854c36043be1b6B",
				"token_address": "0x38D2c9D2425A13Be961ED9eD1C45834B8195FbFB",
				"token_id": 3,
				"token_url": "",
				"price": 0,
				"start_price": 3,
				"status": 3,
				"created_at": "2026-02-03T16:22:37.147+08:00",
				"updated_at": "2026-02-03T17:48:48.902+08:00"
			},
			{
				"id": 1,
				"auction_id": 6,
				"seller": "0x256e72b29D48F1F98792EF1fF854c36043be1b6B",
				"token_address": "0x38D2c9D2425A13Be961ED9eD1C45834B8195FbFB",
				"token_id": 2,
				"token_url": "",
				"price": 2,
				"start_price": 2,
				"status": 2,
				"created_at": "2026-02-03T16:21:36.527+08:00",
				"updated_at": "2026-02-03T17:48:25.327+08:00"
			}
		],
		"pagination": {
			"page": 1,
			"page_size": 20,
			"pages": 1,
			"total": 3
		}
	}
}

### 2. 查询某个拍卖的出价历史记录接口：GET http://localhost:8080/api/v1/bidListOfAuction?auctionId=8
<br/>
响应：<br/>
{
	"code": 200,
	"message": "success",
	"data": [
		{
			"id": 3,
			"auction_id": 8,
			"bidder": "0x76dA744c3D93118218EBD1607ee75E8b9F724292",
			"price": 3,
			"erc20_token": "0x9381CAdC5F5541cCDDF09434a51Ce77243e9Cdc8",
			"created_at": "2026-02-03T18:13:00.233+08:00"
		},
		{
			"id": 4,
			"auction_id": 8,
			"bidder": "0x76dA744c3D93118218EBD1607ee75E8b9F724292",
			"price": 4,
			"erc20_token": "0x9381CAdC5F5541cCDDF09434a51Ce77243e9Cdc8",
			"created_at": "2026-02-03T18:14:12.401+08:00"
		},
		{
			"id": 5,
			"auction_id": 8,
			"bidder": "0x76dA744c3D93118218EBD1607ee75E8b9F724292",
			"price": 5,
			"erc20_token": "0x9381CAdC5F5541cCDDF09434a51Ce77243e9Cdc8",
			"created_at": "2026-02-03T18:14:48.164+08:00"
		},
		{
			"id": 6,
			"auction_id": 8,
			"bidder": "0xcc59f5EF0BD83AE0E17Df52A84c1c9C356571226",
			"price": 6,
			"erc20_token": "0x9381CAdC5F5541cCDDF09434a51Ce77243e9Cdc8",
			"created_at": "2026-02-03T18:17:00.402+08:00"
		},
		{
			"id": 7,
			"auction_id": 8,
			"bidder": "0x76dA744c3D93118218EBD1607ee75E8b9F724292",
			"price": 7,
			"erc20_token": "0x9381CAdC5F5541cCDDF09434a51Ce77243e9Cdc8",
			"created_at": "2026-02-03T18:18:00.384+08:00"
		}
	]
}

### 3. 平台的统计数据（拍卖总数，出价总数）接口：GET http://localhost:8080/api/v1/getBidStatic

响应：
{
	"code": 200,
	"message": "success",
	"data": {
		"auctionCount": 3,
		"bidCount": 7
	}
}

### 4. 查询某个钱包地址拥有的所有NFT Token列表接口：GET http://localhost:8080/api/v1/getNftOfOwner?address=0xcc59f5EF0BD83AE0E17Df52A84c1c9C356571226

响应：
{
	"code": 200,
	"message": "success",
	"data": [
		{
			"id": 1,
			"auction_in_id": 6,
			"auction_out_id": 6,
			"owner": "0xcc59f5EF0BD83AE0E17Df52A84c1c9C356571226",
			"token_address": "0x38D2c9D2425A13Be961ED9eD1C45834B8195FbFB",
			"token_id": 2,
			"token_url": "",
			"status": 2,
			"created_at": "2026-02-03T17:48:25.338+08:00",
			"updated_at": "2026-02-03T17:48:25.341+08:00"
		}
	]
}
