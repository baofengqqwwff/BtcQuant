package main

import (
	"github.com/baofengqqwwff/BtcQuant/engine"
	. "github.com/baofengqqwwff/BtcQuant/strategy"
)

func onBar(bar *engine.Event) (*engine.Event, error) {
	return nil, nil
}
func onTick(tick *engine.Event) (*engine.Event, error) {
	return nil, nil
}

func main() {
	//api暂时只实现了BINANCE
	apis :=[]string{"Binance"}
	stratgyConfigMap := map[string]interface{}{"name":"demo","onTickFunc":onTick,"onBarFunc":onBar,"apis":apis}
	strategy := NewStrategy(stratgyConfigMap)

}
