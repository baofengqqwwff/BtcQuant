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

}

func NewStrategy(argsMap map[string]interface{}) *Strategy {
	strategy := &Strategy{}
	strategy.init()
	//判断参数是否存在
	if _, ok := argsMap["name"]; ok {
		//存在
		name := argsMap["name"].(string)
		strategy.name = name
	}

	//tick处理函数
	if _, ok := argsMap["onTickFunc"]; ok {
		//存在
		onTickFunc := argsMap["onTickFunc"].(func(event *engine.Event) (*engine.Event, error))
		strategy.onTickFunc = onTickFunc
		strategy.engine.RegisterProcessor(&engine.Processor{ProcessorName: "tickProcessor", EventName: "newTick", EventHandler: strategy.onTickFunc})

	}

	//bar处理函数
	if _, ok := argsMap["onBarFunc"]; ok {
		//存在
		onBarFunc := argsMap["onBarFunc"].(func(event *engine.Event) (*engine.Event, error))
		strategy.onBarFunc = onBarFunc
		strategy.engine.RegisterProcessor(&engine.Processor{ProcessorName: "barProcessor", EventName: "newBar", EventHandler: strategy.onBarFunc})
	}

	//bar处理函数
	if _, ok := argsMap["apis"]; ok {
		//存在
		apisList := argsMap["apis"].([]string)
		for _, apiName := range apisList {
			switch strings.ToUpper(apiName) {
			case "BINANCE":
				{
					api.RigisterBinance(strategy.engine)
				}
			}
		}


	}
	strategy.engine.PutEvent(&engine.Event{Name: "logEvent", Data: "初始化成功"})
	return strategy
}
