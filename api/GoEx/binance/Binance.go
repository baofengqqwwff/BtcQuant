package binance

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"time"

	. "BTCCandle/GoEx"
	"math"
	"strings"
)

const (
	EXCHANGE_NAME = "binance.com"

	API_BASE_URL = "https://api.binance.com/"
	API_V1       = API_BASE_URL + "api/v1/"
	API_V3       = API_BASE_URL + "api/v3/"

	TICKER_URI             = "ticker/24hr?symbol=%s"
	TICKERS_URI            = "ticker/allBookTickers"
	DEPTH_URI              = "depth?symbol=%s&limit=%d"
	ACCOUNT_URI            = "account?"
	ORDER_URI              = "order?"
	UNFINISHED_ORDERS_INFO = "openOrders?"
	ALL_ORDERS             = "allOrders?"
	KLINNE_URI             = "klines?"
	TIME_URI               = "time"
	EXCHANGEINFO_URI       = "exchangeInfo"
)

type Binance struct {
	accessKey,
	secretKey string
	httpClient  *http.Client
	timeAdjust  int64
	SymbolsInfo map[string]map[string]int //记录开仓的最小小数点
}

func (bn *Binance)GetSymbols()([]CurrencyPair,error){
	exchangeInfoUri := API_V1 + EXCHANGEINFO_URI
	bodyDataMap, err := HttpGet(bn.httpClient, exchangeInfoUri)
	if err != nil {
		return nil,err
	}
	symbols := []CurrencyPair{}
	symbolsInfo, _ := bodyDataMap["symbols"].([]interface{})
	for _, infoInterface := range symbolsInfo {
		info, _ := infoInterface.(interface{})
		symbolInfo, _ := info.(map[string]interface{})
		symbols = append(symbols,NewCurrencyPair2(symbolInfo["baseAsset"].(string)+"_"+symbolInfo["quoteAsset"].(string)))
	}
	return symbols,nil
}

func (bn *Binance) getExchangeInfo() {
	exchangeInfoUri := API_V1 + EXCHANGEINFO_URI
	bodyDataMap, err := HttpGet(bn.httpClient, exchangeInfoUri)
	if err != nil {
		log.Println("binance初始化失败，重试中")
		time.Sleep(2*time.Second)
		bn.getExchangeInfo()
		return
	}
	symbolsInfo, _ := bodyDataMap["symbols"].([]interface{})
	for _, infoInterface := range symbolsInfo {
		info, _ := infoInterface.(interface{})
		symbolInfo, _ := info.(map[string]interface{})
		symbolName, _ := symbolInfo["symbol"].(string)
		symbolFilterInfo, _ := symbolInfo["filters"].([]interface{})
		//获取价格精度
		pricePrecisionInfo := symbolFilterInfo[0].(map[string]interface{})
		pricePrecisionString, _ := pricePrecisionInfo["minPrice"].(string)
		pricePrecisionFloat, _ := strconv.ParseFloat(pricePrecisionString, 64)
		pricePrecision := int(-math.Log10(pricePrecisionFloat))

		qtyPrecisonInfo := symbolFilterInfo[1].(map[string]interface{})
		qtyPrecisionString, _ := qtyPrecisonInfo["minQty"].(string)
		qtyPrecisionFloat, _ := strconv.ParseFloat(qtyPrecisionString, 64)
		qtyPrecision := int(-math.Log10(qtyPrecisionFloat))

		symbolPrecisonMap := map[string]int{}
		symbolPrecisonMap["pricePrecision"] = pricePrecision
		symbolPrecisonMap["qtyPrecision"] = qtyPrecision
		bn.SymbolsInfo[symbolName] = symbolPrecisonMap
	}
}

func (bn *Binance) SetTimeAdjust(timeAdjust int64) {
	bn.timeAdjust = timeAdjust
}
func (bn *Binance) syncTime() {
	bn.timeAdjust = int64(time.Nanosecond)*time.Now().UnixNano()/int64(time.Millisecond) - bn.GetTime()
}
func (bn *Binance) buildParamsSigned(postForm *url.Values) error {
	postForm.Set("recvWindow", "5000")
	tonce := strconv.FormatInt(int64(time.Nanosecond)*time.Now().UnixNano()/int64(time.Millisecond)-bn.timeAdjust, 10)
	postForm.Set("timestamp", tonce)
	//fmt.Println(bn.GetTime())
	//fmt.Println(strconv.FormatInt(int64(time.Nanosecond) * time.Now().UnixNano() / int64(time.Millisecond), 10))
	//postForm.Set("timestamp", strconv.FormatInt(bn.GetTime(),10))
	payload := postForm.Encode()
	sign, _ := GetParamHmacSHA256Sign(bn.secretKey, payload)
	postForm.Set("signature", sign)
	return nil
}

func New(client *http.Client, api_key, secret_key string) *Binance {
	binance := &Binance{api_key, secret_key, client, 0, map[string]map[string]int{}}
	binance.getExchangeInfo()
	return binance
}

func (bn *Binance) GetExchangeName() string {
	return EXCHANGE_NAME
}
func (bn *Binance) GetTime() int64 {
	timeUri := API_V1 + TIME_URI
	bodyDataMap, _ := HttpGet(bn.httpClient, timeUri)
	time := int64(bodyDataMap["serverTime"].(float64))
	return time

}

func (bn *Binance) GetKlineRecords(currency CurrencyPair, period, size, since int) ([]Kline, error) {
	klineUri := API_V1 + KLINNE_URI
	params := url.Values{}
	params.Set("symbol", currency.ToSymbol(""))
	if size > 0 {
		params.Set("limit", strconv.Itoa(size))
	}
	if since > 0 {
		params.Set("startTime", strconv.Itoa(since))
	}
	var _period string
	switch period {
	case KLINE_PERIOD_1MIN:
		{
			params.Set("interval", "1m")
			_period = "1m"
		}
	case KLINE_PERIOD_5MIN:
		{
			params.Set("interval", "5m")
			_period = "5m"
		}
	case KLINE_PERIOD_15MIN:
		{
			params.Set("interval", "15m")
			_period = "15m"
		}
	case KLINE_PERIOD_30MIN:
		{
			params.Set("interval", "30m")
			_period = "30m"
		}
	case KLINE_PERIOD_60MIN:
		{
			params.Set("interval", "1h")
			_period = "1h"
		}
	case KLINE_PERIOD_4H:
		{
			params.Set("interval", "4h")
			_period = "4h"
		}
	case KLINE_PERIOD_1DAY:
		{
			params.Set("interval", "1d")
			_period = "1d"
		}
	default:
		return nil, errors.New("do not have this period")
	}
	path := klineUri + params.Encode()
	respList, err := HttpGet3(bn.httpClient, path, map[string]string{"X-MBX-APIKEY": bn.accessKey})
	if err != nil {
		log.Println(err)
		if strings.Contains(err.Error(), "-1021") {
			log.Println("同步服务器时间")
			bn.syncTime()
		}
		return nil, err
	}
	klineList := []Kline{}
	for _, resp := range respList {
		kline := Kline{}
		respBodyList, _ := resp.([]interface{})
		kline.Symbol = currency.ToSymbol("_")
		kline.Timestamp = int64(respBodyList[0].(float64))
		kline.Period = _period
		kline.Open, _ = strconv.ParseFloat(respBodyList[1].(string), 64)
		kline.High, _ = strconv.ParseFloat(respBodyList[2].(string), 64)
		kline.Low, _ = strconv.ParseFloat(respBodyList[3].(string), 64)
		kline.Close, _ = strconv.ParseFloat(respBodyList[4].(string), 64)
		kline.Vol, _ = strconv.ParseFloat(respBodyList[5].(string), 64)
		klineList = append(klineList, kline)
	}
	return klineList, nil
}

func (bn *Binance) GetAllBookTickers() ([]*Ticker, error) {
	tickerUri := API_V1 + TICKERS_URI
	//bodyDataMapList, err := HttpGet3(bn.httpClient, tickerUri,nil)
	respData, err := NewHttpRequest(bn.httpClient, "GET", tickerUri, "", nil)
	if err!=nil{
		return nil,err
	}
	var bodyDataMapList []interface{}
	err = json.Unmarshal(respData, &bodyDataMapList)
	if err != nil {
		log.Println("GetTicker error:", err)
		if strings.Contains(err.Error(), "-1021") {
			log.Println("同步服务器时间")
			bn.syncTime()
		}
		return nil, err
	}
	var tickers []*Ticker
	for _, bodyDataMap := range bodyDataMapList {
		var ticker Ticker
		tickerMap, _ := bodyDataMap.(map[string]interface{})
		ticker.Symbol, _ = tickerMap["symbol"].(string)
		ticker.Buy, _ = strconv.ParseFloat(tickerMap["bidPrice"].(string), 10)
		ticker.Sell, _ = strconv.ParseFloat(tickerMap["askPrice"].(string), 10)
		tickers = append(tickers, &ticker)
	}

	return tickers, nil
}

func (bn *Binance) GetTicker(currency CurrencyPair) (*Ticker, error) {
	tickerUri := API_V1 + fmt.Sprintf(TICKER_URI, currency.ToSymbol(""))
	bodyDataMap, err := HttpGet(bn.httpClient, tickerUri)

	if err != nil {
		log.Println("GetTicker error:", err)
		if strings.Contains(err.Error(), "-1021") {
			log.Println("同步服务器时间")
			bn.syncTime()
		}
		return nil, err
	}
	var tickerMap map[string]interface{} = bodyDataMap
	var ticker Ticker

	ticker.Date = uint64(tickerMap["closeTime"].(float64))
	ticker.Last, _ = strconv.ParseFloat(tickerMap["lastPrice"].(string), 10)
	ticker.Buy, _ = strconv.ParseFloat(tickerMap["bidPrice"].(string), 10)
	ticker.Sell, _ = strconv.ParseFloat(tickerMap["askPrice"].(string), 10)
	ticker.Low, _ = strconv.ParseFloat(tickerMap["lowPrice"].(string), 10)
	ticker.High, _ = strconv.ParseFloat(tickerMap["highPrice"].(string), 10)
	ticker.Vol, _ = strconv.ParseFloat(tickerMap["volume"].(string), 10)
	ticker.Symbol, _ = tickerMap["symbol"].(string)
	return &ticker, nil
}

func (bn *Binance) GetDepth(size int, currencyPair CurrencyPair) (*Depth, error) {
	if size > 100 {
		size = 100
	} else if size < 5 {
		size = 5
	}

	apiUrl := fmt.Sprintf(API_V1+DEPTH_URI, currencyPair.ToSymbol(""), size)
	resp, err := HttpGet(bn.httpClient, apiUrl)
	if err != nil {
		log.Println("GetDepth error:", err)
		if strings.Contains(err.Error(), "-1021") {
			log.Println("同步服务器时间")
			bn.syncTime()
		}
		return nil, err
	}

	if _, isok := resp["code"]; isok {
		return nil, errors.New(resp["msg"].(string))
	}

	bids := resp["bids"].([]interface{})
	asks := resp["asks"].([]interface{})

	//log.Println(bids)
	//log.Println(asks)

	depth := new(Depth)

	for _, bid := range bids {
		_bid := bid.([]interface{})
		amount := ToFloat64(_bid[1])
		price := ToFloat64(_bid[0])
		dr := DepthRecord{Amount: amount, Price: price}
		depth.BidList = append(depth.BidList, dr)
	}

	for _, ask := range asks {
		_ask := ask.([]interface{})
		amount := ToFloat64(_ask[1])
		price := ToFloat64(_ask[0])
		dr := DepthRecord{Amount: amount, Price: price}
		depth.AskList = append(depth.AskList, dr)
	}

	return depth, nil
}

func (bn *Binance) placeOrder(amount, price string, pair CurrencyPair, orderType, orderSide string) (*Order, error) {
	path := API_V3 + ORDER_URI
	params := url.Values{}
	params.Set("symbol", pair.ToSymbol(""))
	params.Set("side", orderSide)
	params.Set("type", orderType)

	params.Set("quantity", amount)
	params.Set("type", "LIMIT")
	params.Set("timeInForce", "GTC")

	switch orderType {
	case "LIMIT":
		params.Set("price", price)
	}

	bn.buildParamsSigned(&params)

	resp, err := HttpPostForm2(bn.httpClient, path, params,
		map[string]string{"X-MBX-APIKEY": bn.accessKey})
	//log.Println("resp:", string(resp), "err:", err)
	if err != nil {
		return nil, err
	}

	respmap := make(map[string]interface{})
	err = json.Unmarshal(resp, &respmap)
	if err != nil {
		log.Println(string(resp))
		return nil, err
	}
	if _, isok := respmap["code"]; isok == true {
		return nil, errors.New(respmap["msg"].(string))
	}

	orderId, isok := respmap["orderId"].(string)
	if !isok {
		return nil, errors.New(string(resp))
	}

	side := BUY
	if orderSide == "SELL" {
		side = SELL
	}
	return &Order{
		Currency:   pair,
		OrderID:    ToInt(orderId),
		Price:      ToFloat64(price),
		Amount:     ToFloat64(amount),
		DealAmount: 0,
		AvgPrice:   0,
		Side:       TradeSide(side),
		Status:     ORDER_UNFINISH,
		OrderTime:  int(time.Now().Unix())}, nil
}

func (bn *Binance) GetAccount() (*Account, error) {
	params := url.Values{}
	bn.buildParamsSigned(&params)
	path := API_V3 + ACCOUNT_URI + params.Encode()
	respmap, err := HttpGet2(bn.httpClient, path, map[string]string{"X-MBX-APIKEY": bn.accessKey})
	if err != nil {
		log.Println(err)
		if strings.Contains(err.Error(), "-1021") {
			log.Println("同步服务器时间")
			bn.syncTime()
		}
		return nil, err
	}
	//log.Println("respmap:", respmap)
	if _, isok := respmap["code"]; isok == true {
		return nil, errors.New(respmap["msg"].(string))
	}
	acc := Account{}
	acc.Exchange = bn.GetExchangeName()
	acc.SubAccounts = make(map[Currency]SubAccount)

	balances := respmap["balances"].([]interface{})
	for _, v := range balances {
		//log.Println(v)
		vv := v.(map[string]interface{})
		currency := NewCurrency(vv["asset"].(string), "")
		acc.SubAccounts[currency] = SubAccount{
			Currency:     currency,
			Amount:       ToFloat64(vv["free"]),
			ForzenAmount: ToFloat64(vv["locked"]),
		}
	}

	return &acc, nil
}

func (bn *Binance) LimitBuy(amount, price string, currencyPair CurrencyPair) (*Order, error) {
	return bn.placeOrder(amount, price, currencyPair, "LIMIT", "BUY")
}

func (bn *Binance) LimitSell(amount, price string, currencyPair CurrencyPair) (*Order, error) {
	return bn.placeOrder(amount, price, currencyPair, "LIMIT", "SELL")
}

func (bn *Binance) MarketBuy(amount, price string, currencyPair CurrencyPair) (*Order, error) {
	return bn.placeOrder(amount, price, currencyPair, "MARKET", "BUY")
}

func (bn *Binance) MarketSell(amount, price string, currencyPair CurrencyPair) (*Order, error) {
	return bn.placeOrder(amount, price, currencyPair, "MARKET", "SELL")
}

func (bn *Binance) CancelOrder(orderId string, currencyPair CurrencyPair) (bool, error) {
	path := API_V3 + ORDER_URI
	params := url.Values{}
	params.Set("symbol", currencyPair.ToSymbol(""))
	params.Set("orderId", orderId)

	bn.buildParamsSigned(&params)

	resp, err := HttpDeleteForm(bn.httpClient, path, params, map[string]string{"X-MBX-APIKEY": bn.accessKey})

	//log.Println("resp:", string(resp), "err:", err)
	if err != nil {
		return false, err
	}

	respmap := make(map[string]interface{})
	err = json.Unmarshal(resp, &respmap)
	if err != nil {
		log.Println(string(resp))
		return false, err
	}

	orderIdCanceled, isok := respmap["orderId"].(string)
	if !isok {
		return false, errors.New(string(resp))
	}
	if orderIdCanceled != orderId {
		return false, errors.New("orderId doesn't match")
	}

	return true, nil
}

func (bn *Binance) GetOneOrder(orderId string, currencyPair CurrencyPair) (*Order, error) {
	params := url.Values{}
	params.Set("symbol", currencyPair.ToSymbol(""))
	params.Set("orderId", orderId)

	bn.buildParamsSigned(&params)
	path := API_V3 + ORDER_URI + params.Encode()

	respmap, err := HttpGet2(bn.httpClient, path, map[string]string{"X-MBX-APIKEY": bn.accessKey})

	if err != nil {
		if strings.Contains(err.Error(), "-1021") {
			log.Println("同步服务器时间")
			bn.syncTime()
		}
		return nil, err
	}
	status := respmap["status"].(string)

	ord := Order{}
	ord.Currency = currencyPair
	ord.OrderID = ToInt(orderId)

	if status == "FILLED" {
		ord.Status = ORDER_FINISH
	} else {
		ord.Status = ORDER_UNFINISH
	}
	ord.Amount = ToFloat64(respmap["origQty"].(string))
	ord.Price = ToFloat64(respmap["price"].(string))

	return &ord, nil
}

func (bn *Binance) GetUnfinishOrders(currencyPair CurrencyPair) ([]Order, error) {
	params := url.Values{}
	params.Set("symbol", currencyPair.ToSymbol(""))

	bn.buildParamsSigned(&params)
	path := API_V3 + UNFINISHED_ORDERS_INFO + params.Encode()

	respmap, err := HttpGet3(bn.httpClient, path, map[string]string{"X-MBX-APIKEY": bn.accessKey})
	//log.Println("respmap", respmap, "err", err)
	if err != nil {
		if strings.Contains(err.Error(), "-1021") {
			log.Println("同步服务器时间")
			bn.syncTime()
		}
		return nil, err
	}

	orders := make([]Order, 0)
	for _, v := range respmap {
		ord := v.(map[string]interface{})
		side := ord["type"].(string)
		orderSide := SELL
		if side == "BUY" {
			orderSide = BUY
		}

		orders = append(orders, Order{
			OrderID:   ToInt(ord["orderId"]),
			Currency:  currencyPair,
			Price:     ToFloat64(ord["price"]),
			Amount:    ToFloat64(ord["origQty"]),
			Side:      TradeSide(orderSide),
			Status:    ORDER_UNFINISH,
			OrderTime: ToInt(ord["time"])})
	}
	return orders, nil
}

func (bn *Binance) GetAllOrders(currencyPair CurrencyPair) ([]Order, error) {
	params := url.Values{}
	params.Set("symbol", currencyPair.ToSymbol(""))

	bn.buildParamsSigned(&params)
	path := API_V3 + ALL_ORDERS + params.Encode()

	respmap, err := HttpGet3(bn.httpClient, path, map[string]string{"X-MBX-APIKEY": bn.accessKey})
	//log.Println("respmap", respmap, "err", err)
	if err != nil {
		if strings.Contains(err.Error(), "-1021") {
			log.Println("同步服务器时间")
			bn.syncTime()
		}
		return nil, err
	}

	orders := make([]Order, 0)
	for _, v := range respmap {
		ord := v.(map[string]interface{})
		side := ord["type"].(string)
		orderSide := SELL
		if side == "BUY" {
			orderSide = BUY
		}
		ordStatus := ord["status"].(string)
		var status TradeStatus
		switch ordStatus {
		case "NEW":
			status = ORDER_UNFINISH
		case "PARTIALLY_FILLED":
			status = ORDER_PART_FINISH
		case "FILLED":
			status = ORDER_FINISH
		case "CANCELED":
			status = ORDER_CANCEL
		case "REJECTED":
			status = ORDER_REJECT
		case "PENDING_CANCEL":
			status = ORDER_CANCEL_ING
		case "EXPIRED":
			status = ORDER_EXPIRED
		}

		orders = append(orders, Order{
			OrderID:    ToInt(ord["orderId"]),
			Currency:   currencyPair,
			Price:      ToFloat64(ord["price"]),
			Amount:     ToFloat64(ord["origQty"]),
			DealAmount: ToFloat64(ord["executedQty"]),
			Side:       TradeSide(orderSide),
			Status:     status,
			OrderTime:  ToInt(ord["time"])})
	}
	return orders, nil
}
