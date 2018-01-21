package strategy

import (
	"github.com/baofengqqwwff/BtcQuant/engine"
	"strings"
	"github.com/baofengqqwwff/BtcQuant/api"
)

type Strategy struct {
	name   string
	engine *engine.Engine
	onBar  func(bar *engine.Event)
	onTick func(tick *engine.Event)
}

func (strategy *Strategy) init() {
	strategy.engine = engine.NewEngine()
	strategy.engine.RegisterProcessor(&engine.Processor{ProcessorName: "barProcessor", EventName: "newBar", EventHandler: strategy.onBar})
	strategy.engine.RegisterProcessor(&engine.Processor{ProcessorName: "tickProcessor", EventName: "newTick", EventHandler: strategy.onTick})
}

func NewStrategy(argsMap map[string]interface{}) *Strategy {
	strategy := &Strategy{}
	//判断参数是否存在
	if _, ok := argsMap["name"]; ok {
		//存在
		name := argsMap["name"].(string)
		strategy.name = name
	}

	//tick处理函数
	if _, ok := argsMap["onTick"]; ok {
		//存在
		onTickFunc := argsMap["onTick"].(func(event *engine.Event))
		strategy.onTick = onTickFunc
	}

	//bar处理函数
	if _, ok := argsMap["onBar"]; ok {
		//存在
		onBarFunc := argsMap["onBar"].(func(event *engine.Event))
		strategy.onBar = onBarFunc
	}

	//bar处理函数
	if _, ok := argsMap["apis"]; ok {
		//存在
		apisList := argsMap["onBar"].([]string)
		for _, apiName := range apisList {
			switch strings.ToUpper(apiName) {
			case "BINANCE":
				{
					api.RigisterBinance(strategy.engine)
				}
			}
		}
		strategy.init()

	}
	return strategy
}
