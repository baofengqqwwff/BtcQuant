package goex

import "net/http"

type Order struct {
	Price,
	Amount,
	AvgPrice,
	DealAmount,
	Fee float64
	OrderID2  string
	OrderID   int
	OrderTime int
	Status    TradeStatus
	Currency  CurrencyPair
	Side      TradeSide
}

type Trade struct {
	Tid    int64   `json:"tid"`
	Type   string  `json:"type"`
	Amount float64 `json:"amount,string"`
	Price  float64 `json:"price,string"`
	Date   int64   `json:"date_ms"`
}

type SubAccount struct {
	Currency Currency
	Amount,
	ForzenAmount,
	LoanAmount float64
}

type Account struct {
	Exchange    string
	Asset       float64 //总资产
	NetAsset    float64 //净资产
	SubAccounts map[Currency]SubAccount
}

type Ticker struct {
	Symbol string `json:"symbol"`
	Last float64 `json:"last"`
	Buy  float64 `json:"buy"`
	Sell float64 `json:"sell"`
	High float64 `json:"high"`
	Low  float64 `json:"low"`
	Vol  float64 `json:"vol"`
	Date uint64  `json:"date"` // 单位:秒(second)
}

type DepthRecord struct {
	Price,
	Amount float64
}

type DepthRecords []DepthRecord

func (dr DepthRecords) Len() int {
	return len(dr)
}

func (dr DepthRecords) Swap(i, j int) {
	dr[i], dr[j] = dr[j], dr[i]
}

func (dr DepthRecords) Less(i, j int) bool {
	return dr[i].Price < dr[j].Price
}

type Depth struct {
	AskList,
	BidList DepthRecords
}

type APIConfig struct {
	HttpClient *http.Client
	ApiUrl,
	AccessKey,
	SecretKey string
}

type Kline struct {
	Symbol string
	Period string
	Timestamp int64
	Open,
	Close,
	High,
	Low,
	Vol float64
}

type FutureKline struct {
	*Kline
	Vol2 float64 //个数
}

type FutureSubAccount struct {
	Currency      Currency
	AccountRights float64 //账户权益
	KeepDeposit   float64 //保证金
	ProfitReal    float64 //已实现盈亏
	ProfitUnreal  float64
	RiskRate      float64 //保证金率
}

type FutureAccount struct {
	FutureSubAccounts map[Currency]FutureSubAccount
}

type FutureOrder struct {
	Price        float64
	Amount       float64
	AvgPrice     float64
	DealAmount   float64
	OrderID      int64
	OrderTime    int64
	Status       TradeStatus
	Currency     CurrencyPair
	OType        int     //1：开多 2：开空 3：平多 4： 平空
	LeverRate    int     //倍数
	Fee          float64 //手续费
	ContractName string
}

type FuturePosition struct {
	BuyAmount      float64
	BuyAvailable   float64
	BuyPriceAvg    float64
	BuyPriceCost   float64
	BuyProfitReal  float64
	CreateDate     int64
	LeverRate      int
	SellAmount     float64
	SellAvailable  float64
	SellPriceAvg   float64
	SellPriceCost  float64
	SellProfitReal float64
	Symbol         CurrencyPair //btc_usd:比特币,ltc_usd:莱特币
	ContractType   string
	ContractId     int64
	ForceLiquPrice float64 //预估爆仓价
}
