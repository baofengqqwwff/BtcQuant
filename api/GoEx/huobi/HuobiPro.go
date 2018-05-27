package huobi

import (
	. "BTCCandle/GoEx"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"github.com/baofengqqwwff/GoEx"
	"github.com/pkg/errors"
)

const (
	SYMBOLS_URL = "/v1/common/symbols"
)

type HuobiPro struct {
	*HuoBi_V2
}

func NewHuobiPro(client *http.Client, apikey, secretkey, accountId string) *HuobiPro {
	hbv2 := new(HuoBi_V2)
	hbv2.accountId = accountId
	hbv2.accessKey = apikey
	hbv2.secretKey = secretkey
	hbv2.httpClient = client
	hbv2.baseUrl = "https://api.huobipro.com"
	return &HuobiPro{hbv2}
}

func (hbpro *HuobiPro) GetExchangeName() string {
	return HUOBI_PRO
}

func (hbpro *HuobiPro) GetSymbols() ([]CurrencyPair, error) {
	symbolMapInterface, err := HttpGet2(hbpro.httpClient, hbpro.baseUrl+SYMBOLS_URL, nil)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	symbolDataInterface := symbolMapInterface["data"].([]interface{})
	symbols := []CurrencyPair{}
	for _, valueInterface := range symbolDataInterface {
		value := valueInterface.(map[string]interface{})
		currencyA := value["base-currency"].(string)
		currencyB := value["quote-currency"].(string)
		symbols = append(symbols, NewCurrencyPair2(strings.ToUpper(currencyA+"_"+currencyB)))
	}
	//fmt.Println(symbolMapInterface)
	return symbols, nil
}

func (hbpro *HuobiPro) GetKlineRecords(currency CurrencyPair, period, size, since int) ([]Kline, error) {
	v := url.Values{}
	v.Set("symbol", strings.ToLower(currency.ToSymbol("")))
	v.Set("size", strconv.Itoa(size))
	var _period string
	switch period {
	case goex.KLINE_PERIOD_1MIN:
		v.Set("period", "1min")
		_period = "1m"
	case goex.KLINE_PERIOD_5MIN:
		v.Set("period", "5min")
		_period = "5m"
	case goex.KLINE_PERIOD_15MIN:
		v.Set("period", "15min")
		_period = "15m"
	case goex.KLINE_PERIOD_30MIN:
		v.Set("period", "30min")
		_period = "30m"
	case goex.KLINE_PERIOD_60MIN:
		v.Set("period", "60min")
		_period = "60m"
	case goex.KLINE_PERIOD_1DAY:
		v.Set("period", "1day")
		_period = "1d"
	case goex.KLINE_PERIOD_1WEEK:
		v.Set("period", "1week")
		_period = "1w"
	default:
		return nil, errors.New("no this period")
	}
	if since != 0 {
		log.Println("please notice that since is not been used")
	}

	klineMapInterface, err := HttpGet2(hbpro.httpClient, hbpro.baseUrl+"/market/history/kline?"+v.Encode(), nil)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	klinesInterface := klineMapInterface["data"].([]interface{})
	var klineRecords []Kline

	for _, klineInterface := range klinesInterface {
		record := klineInterface.(map[string]interface{})
		klineRecord := Kline{}
		klineRecord.Close = record["close"].(float64)
		klineRecord.High = record["high"].(float64)
		klineRecord.Open = record["open"].(float64)
		klineRecord.Low = record["low"].(float64)
		klineRecord.Period = _period
		klineRecord.Symbol = currency.ToSymbol("_")
		klineRecord.Timestamp = int64(record["id"].(float64)) * 1000
		klineRecord.Vol = record["amount"].(float64)
		klineRecords = append(klineRecords, klineRecord)
	}

	return klineRecords, nil
}
