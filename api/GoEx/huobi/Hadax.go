package huobi

import (
	"net/http"
	"log"

	. "BTCCandle/GoEx"
)


type Hadax struct {
	*HuobiPro
}

func NewHadax(client *http.Client, apikey, secretkey, accountId string) *Hadax {
	hbv2 := new(HuoBi_V2)
	hbv2.accountId = accountId
	hbv2.accessKey = apikey
	hbv2.secretKey = secretkey
	hbv2.httpClient = client
	hbv2.baseUrl = "https://api.hadax.com"
	return &Hadax{&HuobiPro{hbv2}}
}

func (hadax *Hadax) GetExchangeName() string {
	return "hadax.com"
}

func (hadax *Hadax) GetSymbols() ([]CurrencyPair, error) {
	symbolMapInterface, err := HttpGet2(hadax.httpClient, hadax.baseUrl+"/v1/hadax/common/symbols", nil)
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
		symbols = append(symbols, NewCurrencyPair2(currencyA+"_"+currencyB))
	}
	//fmt.Println(symbolMapInterface)
	return symbols, nil
}
