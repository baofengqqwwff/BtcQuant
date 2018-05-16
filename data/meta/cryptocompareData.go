package meta

type PriceMulti map[string]CoinPrices

type CoinPrices map[string]float64

type CryptocompareKlines []OHLCV

type OHLCV struct {
	Time int
	Close float64
	High float64
	Low float64
	Open float64
	Volumefrom float64
	Volumeto float64
}