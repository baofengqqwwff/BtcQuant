package strategy

import (
	"github.com/baofengqqwwff/BtcQuant/engine"
	"strings"
	"github.com/baofengqqwwff/BtcQuant/api"
)

type Strategy struct {
	name       string
	engine     *engine.Engine
	onBarFunc  func(bar *engine.Event) (*engine.Event, error)
	onTickFunc func(tick *engine.Event) (*engine.Event, error)
}

func (strategy *Strategy) init() {
	strategy.engine = engine.NewEngine()
	strategy.engine.RegisterProcessor(&engine.Processor{ProcessorName: "barProcessor", EventName: "newBar", EventHandler: strategy.onBarFunc})
	strategy.engine.RegisterProcessor(&engine.Processor{ProcessorName: "tickProcessor", EventName: "newTick", EventHandler: strategy.onTickFunc})
	strategy.engine.PutEvent(&engine.Event{Name: "logEvent", Data: "初始化成功"})
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
	if _, ok := argsMap["onTickFunc"]; ok {
		//存在
		onTickFunc := argsMap["onTick"].(func(event *engine.Event) (*engine.Event, error))
		strategy.onTickFunc = onTickFunc
	}

	//bar处理函数
	if _, ok := argsMap["onBarFunc"]; ok {
		//存在
		onBarFunc := argsMap["onBar"].(func(event *engine.Event) (*engine.Event, error))
		strategy.onBarFunc = onBarFunc
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
