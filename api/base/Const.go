package api

type TradeSide int

const (
	BUY = 1 + iota
	SELL
	BUY_MARKET
	SELL_MARKET
)

func (ts TradeSide) String() string {
	switch ts {
	case 1:
		return "BUY"
	case 2:
		return "SELL"
	case 3:
		return "BUY_MARKET"
	case 4:
		return "SELL_MARKET"
	default:
		return "UNKNOWN"
	}
}

type TradeStatus int

func (ts TradeStatus) String() string {
	return tradeStatusSymbol[ts]
}

var tradeStatusSymbol = [...]string{"UNFINISH", "PART_FINISH", "FINISH", "CANCEL", "REJECT", "CANCEL_ING"}

const (
	ORDER_UNFINISH = iota
	ORDER_PART_FINISH
	ORDER_FINISH
	ORDER_CANCEL
	ORDER_REJECT
	ORDER_CANCEL_ING
	ORDER_EXPIRED
)

const (
	OPEN_BUY   = 1 + iota //开多
	OPEN_SELL             //开空
	CLOSE_BUY             //平多
	CLOSE_SELL            //平空
)

//k线周期
const (
	KLINE_PERIOD_1MIN = 1 + iota
	KLINE_PERIOD_5MIN
	KLINE_PERIOD_15MIN
	KLINE_PERIOD_30MIN
	KLINE_PERIOD_60MIN
	KLINE_PERIOD_4H
	KLINE_PERIOD_1DAY
	KLINE_PERIOD_1WEEK
)

var (
	THIS_WEEK_CONTRACT = "this_week" //周合约
	NEXT_WEEK_CONTRACT = "next_week" //次周合约
	QUARTER_CONTRACT   = "quarter"   //季度合约
)
