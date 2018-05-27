package main

import (
	"net/http"
	"io/ioutil"
	"github.com/baofengqqwwff/BtcQuant/data/meta"
	"net/url"
	"encoding/json"
	"fmt"
)

func GetPricemulti(arg map[string]string)(meta.PriceMulti,error){
	v := url.Values{}
	for key,value := range arg{
		v.Set(key,value)
	}
	res, err := http.Get("https://min-api.cryptocompare.com/data/pricemulti?"+v.Encode())
	if err!=nil{
		return nil,err
	}
	defer res.Body.Close()

	resbody, _ := ioutil.ReadAll(res.Body)

	var msg meta.PriceMulti
	err = json.Unmarshal(resbody, &msg)
	if err!=nil{
		return nil,err
	}
	return msg,nil
}

func main() {
	//res, err := http.Get("https://min-api.cryptocompare.com/data/pricemulti?fsyms=BTC,ETH&tsyms=USD,EUR")
	//res, err := http.Get("https://min-api.cryptocompare.com/data/price?fsym=BTC&tsyms=USD,JPY,EUR")
	//res, err := http.Get("https://min-api.cryptocompare.com/data/histominute?fsym=BTC&tsym=GBP&limit=10")

	fmt.Println(GetPricemulti(map[string]string{"fsyms":"BTC","tsyms":"USD,EUR"}))
	//v := url.Values{}
	//v.Set("fsym", "BTC")
	//v.Set("tsyms", "USD,JPY,EUR")
	//body := v.Encode() //把form数据编下码
	//res, _ := http.Get("https://min-api.cryptocompare.com/data/price?"+body)
	//
	//defer res.Body.Close()
	//
	//resbody, err := ioutil.ReadAll(res.Body)
	//if err != nil {
	//	panic(err)
	//}
	//
	//var msg meta.OHLCV
	//fmt.Println(msg,resbody)
	//var responbody map[string]interface{}
	//json.Unmarshal(resbody, &responbody)
	//msgInterfaces := responbody["Data"].([]interface{})
	//msg = msgInterfaces[0].(meta.OHLCV)
	//fmt.Println(msg)
}

//res, err := http.Get("https://min-api.cryptocompare.com/data/pricemulti?fsyms=BTC,ETH&tsyms=USD,EUR")
//res, err := http.Get("https://min-api.cryptocompare.com/data/price?fsym=BTC&tsyms=USD,JPY,EUR")
//res, err := http.Get("https://min-api.cryptocompare.com/data/histominute?fsym=BTC&tsym=GBP&limit=10")
//
//v := url.Values{}
//v.Set("fsym", "BTC")
//v.Set("tsyms", "USD,JPY,EUR")
//body := v.Encode() //把form数据编下码
//res, _ := http.Get("https://min-api.cryptocompare.com/data/price?"+body)
//
//defer res.Body.Close()
//
//resbody, err := ioutil.ReadAll(res.Body)
//if err != nil {
//panic(err)
//}
//
//var msg meta.OHLCV
//fmt.Println(msg,resbody)
////var responbody map[string]interface{}
////json.Unmarshal(resbody, &responbody)
////msgInterfaces := responbody["Data"].([]interface{})
////msg = msgInterfaces[0].(meta.OHLCV)
////fmt.Println(msg)
