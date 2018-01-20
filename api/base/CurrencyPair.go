package api

import (
	"strings"
)

var UNKNOWN = Currency{"UNKNOWN", ""}
var UNKNOWN_PAIR = CurrencyPair{UNKNOWN, UNKNOWN}

type Currency struct {
	Symbol string
	Desc   string
}

type CurrencyPair struct {
	CurrencyA Currency
	CurrencyB Currency
}

func (c *Currency) String() string {
	return c.Symbol
}

func NewCurrency(symbol, desc string) Currency {
	return Currency{symbol, desc}
}

func NewCurrencyPair(currencyA Currency, currencyB Currency) CurrencyPair {
	return CurrencyPair{currencyA, currencyB}
}

func NewCurrencyPair2(currencyPairSymbol string) CurrencyPair {
	currencys := strings.Split(currencyPairSymbol, "_")
	if len(currencys) == 2 {
		return CurrencyPair{NewCurrency(currencys[0], ""),
			NewCurrency(currencys[1], "")}
	}
	return UNKNOWN_PAIR
}

func (pair *CurrencyPair) ToSymbol(joinChar string, reverse bool) string {
	if reverse {
		return strings.Join([]string{pair.CurrencyB.Symbol, pair.CurrencyA.Symbol}, joinChar)
	}
	return strings.Join([]string{pair.CurrencyA.Symbol, pair.CurrencyB.Symbol}, joinChar)
}
