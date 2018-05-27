package hitbtc

import (
	"errors"
	. "BTCCandle/GoEx"
	//"log"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"time"
	"log"
	"strconv"
	"net/url"
)

const (
	EXCHANGE_NAME = "hitbtc.com"

	API_BASE_URL = "https://api.hitbtc.com/"
	API_V2       = "api/2/"
	TICKER_URI   = "public/ticker/"
	SYMBOLS_URL  = "public/symbol"
)

type Hitbtc struct {
	accessKey,
	secretKey string
	httpClient *http.Client
}

func New(client *http.Client, accessKey, secretKey string) *Hitbtc {
	return &Hitbtc{accessKey, secretKey, client}
}

func (hitbtc *Hitbtc) GetExchangeName() string {
	return EXCHANGE_NAME
}

func (hitbtc *Hitbtc) GetTicker(currency CurrencyPair) (*Ticker, error) {
	currency = hitbtc.adaptCurrencyPair(currency)
	curr := currency.ToSymbol("")
	tickerUri := API_BASE_URL + API_V2 + TICKER_URI + curr
	bodyDataMap, err := HttpGet(hitbtc.httpClient, tickerUri)
	//log.Println("Hitbtc bodyDataMap:", tickerUri, bodyDataMap, err)

	if err != nil {
		//log.Println(err)
		return nil, err
	}
	if result, isok := bodyDataMap["error"].(map[string]interface{}); isok == true {
		//log.Println("bodyDataMap[\"error\"]", result)
		return nil, errors.New(result["message"].(string) + ", " + result["description"].(string))
	}

	tickerMap := bodyDataMap
	var ticker Ticker

	timestamp := time.Now().Unix()
	ticker.Date = uint64(timestamp)
	ticker.Last = ToFloat64(tickerMap["last"])
	ticker.Buy = ToFloat64(tickerMap["bid"])
	ticker.Sell = ToFloat64(tickerMap["ask"])
	ticker.Low = ToFloat64(tickerMap["low"])
	ticker.High = ToFloat64(tickerMap["high"])
	ticker.Vol = ToFloat64(tickerMap["volume"])

	//log.Println("Hitbtc", currency, "ticker:", ticker)

	return &ticker, nil
}

func (hitbtc *Hitbtc) GetTickers(currency CurrencyPair) (*Ticker, error) {
	return hitbtc.GetTicker(currency)
}

func (hitbtc *Hitbtc) GetSymbols() ([]CurrencyPair, error) {
	symbolListInterface, err := HttpGet3(hitbtc.httpClient, API_BASE_URL+API_V2+SYMBOLS_URL, nil)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	symbols := []CurrencyPair{}
	for _, symbolInterface := range symbolListInterface {
		symbolMapInterface := symbolInterface.(map[string]interface{})

		symbol := NewCurrencyPair2(symbolMapInterface["baseCurrency"].(string) + "_" + symbolMapInterface["quoteCurrency"].(string))

		symbols = append(symbols, symbol)
	}

	return symbols, nil
}

func (hitbtc *Hitbtc) GetTickerInBuf(currency CurrencyPair) (*Ticker, error) {
	return hitbtc.GetTicker(currency)
}

func (hitbtc *Hitbtc) LimitBuy(amount, price string, currency CurrencyPair) (*Order, error) {
	panic("not implements")
}

func (hitbtc *Hitbtc) LimitSell(amount, price string, currency CurrencyPair) (*Order, error) {
	panic("not implements")
}

func (hitbtc *Hitbtc) MarketBuy(amount, price string, currency CurrencyPair) (*Order, error) {
	panic("not implements")
}

func (hitbtc *Hitbtc) MarketSell(amount, price string, currency CurrencyPair) (*Order, error) {
	panic("not implements")
}

func (hitbtc *Hitbtc) CancelOrder(orderId string, currency CurrencyPair) (bool, error) {
	panic("not implements")
}

func (hitbtc *Hitbtc) GetOneOrder(orderId string, currency CurrencyPair) (*Order, error) {
	panic("not implements")
}
func (hitbtc *Hitbtc) GetUnfinishOrders(currency CurrencyPair) ([]Order, error) {
	panic("not implements")
}

func (hitbtc *Hitbtc) GetOrderHistorys(currency CurrencyPair, currentPage, pageSize int) ([]Order, error) {
	panic("not implements")
}

func (hitbtc *Hitbtc) GetAccount() (*Account, error) {
	var ret []interface{}
	err := hitbtc.doRequest("GET", "trading/balance", &ret)
	//log.Println(ret)
	if err != nil {
		return nil, err
	}

	acc := new(Account)
	acc.SubAccounts = make(map[Currency]SubAccount, 1)

	for _, v := range ret {
		vv := v.(map[string]interface{})
		currency := NewCurrency(vv["currency"].(string), "")
		acc.SubAccounts[currency] = SubAccount{
			Currency:     currency,
			Amount:       ToFloat64(vv["available"]),
			ForzenAmount: ToFloat64(vv["reserved"])}
	}

	return acc, nil
}

func (hitbtc *Hitbtc) GetDepth(size int, currency CurrencyPair) (*Depth, error) {
	panic("not implement")
}
func (hitbtc *Hitbtc) adaptCurrencyPair(pair CurrencyPair) CurrencyPair {
	var currencyA Currency
	var currencyB Currency

	if pair.CurrencyA == BCC {
		currencyA = BCH
	} else {
		currencyA = pair.CurrencyA
	}
	//currencyB = pair.BaseCurrency
	if pair.CurrencyB == USDT {
		currencyB = USD
	} else {
		currencyB = pair.CurrencyB
	}

	return NewCurrencyPair(currencyA, currencyB)
}

func (hitbtc *Hitbtc) GetKlineRecords(currencyPair CurrencyPair, period, size, since int) ([]Kline, error) {
	//https://api.hitbtc.com/api/2/public/candles/ETHBTC?period=M30
	if since != 0 {
		log.Println("please notice that hitbtc GetKlineRecords since is not been used")
	}
	v := url.Values{}
	v.Set("limit",strconv.Itoa(size))
	var _period string
	switch period {
	case KLINE_PERIOD_1MIN:
		v.Set("period", "M1")
		_period = "1m"
	case KLINE_PERIOD_5MIN:
		v.Set("period", "M3")
		_period = "5m"
	case KLINE_PERIOD_15MIN:
		v.Set("period", "M15")
		_period = "15m"
	case KLINE_PERIOD_30MIN:
		v.Set("period", "M30")
		_period = "30m"
	case KLINE_PERIOD_60MIN:
		v.Set("period", "H1")
		_period = "60m"
	case KLINE_PERIOD_1DAY:
		v.Set("period", "D1")
		_period = "1day"
	case KLINE_PERIOD_1WEEK:
		v.Set("period", "D7")
		_period = "1week"
	default:
		return nil, errors.New("no this period")
	}
	klineInterfaces, err := HttpGet3(hitbtc.httpClient, "https://api.hitbtc.com/api/2/public/candles/"+strings.ToUpper(currencyPair.ToSymbol(""))+"?"+v.Encode(), nil)
	if err!= nil{
		return nil,err
	}

	klines :=[]Kline{}
	for _,klineInterface := range klineInterfaces{
		klineMapInterface := klineInterface.(map[string]interface{})

		kline := Kline{}
		kline.Symbol = strings.ToUpper(currencyPair.ToSymbol("_"))
		kline.Period = _period
		dateTime,_ := time.Parse("2006-01-02T15:04:05.000Z",klineMapInterface["timestamp"].(string))
		kline.Timestamp = dateTime.Unix()*1000
		kline.Open,_ = strconv.ParseFloat(klineMapInterface["open"].(string),64)
		kline.Close,_ = strconv.ParseFloat(klineMapInterface["close"].(string),64)
		kline.High,_ = strconv.ParseFloat(klineMapInterface["max"].(string),64)
		kline.Low,_ = strconv.ParseFloat(klineMapInterface["min"].(string),64)
		kline.Vol,_ = strconv.ParseFloat(klineMapInterface["volume"].(string),64)
		klines = append(klines,kline)
	}
	return klines,err
}

//非个人，整个交易所的交易记录
func (hitbtc *Hitbtc) GetTrades(currencyPair CurrencyPair, since int64) ([]Trade, error) {
	panic("not implements")
}

func (hitbtc *Hitbtc) doRequest(reqMethod, uri string, ret interface{}) error {
	url := API_BASE_URL + API_V2 + uri
	req, _ := http.NewRequest(reqMethod, url, strings.NewReader(""))
	req.SetBasicAuth(hitbtc.accessKey, hitbtc.secretKey)
	resp, err := hitbtc.httpClient.Do(req)
	if err != nil {
		return err
	}

	defer resp.Body.Close()

	bodyData, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	println(string(bodyData))

	if resp.StatusCode != 200 {
		return errors.New(fmt.Sprintf("HttpStatusCode:%d ,Desc:%s", resp.StatusCode, string(bodyData)))
	}

	err = json.Unmarshal(bodyData, ret)

	if err != nil {
		return err
	}

	return nil
}
