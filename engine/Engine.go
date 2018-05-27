package engine

import (
	"container/list"
	"github.com/pkg/errors"
	"time"
)

type Engine struct {
	Event        *list.List    //事件
	Processor    *list.List    //事件处理器
	TimeInterval time.Duration //时间间隔，默认一秒
}

func NewEngine() *Engine {
	engine := &Engine{list.New(), list.New(), time.Second}
	engine.RegisterProcessor(&Processor{ProcessorName: "logProcessor", EventName: "logEvent", EventHandler: LogHandler})
	engine.startPutTimeEvent()
	engine.engineStart()
	return engine
}

//启动引擎
func (engine *Engine) engineStart() {
	engine.PutEvent(&Event{Name: "logEvent", Data: "引擎启动"})
	go func() {
		for {
			eventElem := engine.Event.Front()
			if eventElem != nil {
				event := eventElem.Value.(*Event)
				engine.Event.Remove(eventElem)
				go func() {
					engine.processEvent(event)
				}()
			}
		}
	}()

}

//设置定时器间隔
func (engine *Engine) SetTimeInterval(timeInterval time.Duration) {
	engine.TimeInterval = timeInterval
}

//推送事件
func (engine *Engine) PutEvent(event *Event) {
	engine.Event.PushBack(event)
}

//处理事件
func (engine *Engine) processEvent(event *Event) {

	for selectProcessorElem := engine.Processor.Front(); selectProcessorElem != nil; selectProcessorElem = selectProcessorElem.Next() {
		selectProcessor, _ := selectProcessorElem.Value.(*Processor)
		if selectProcessor.EventName == event.Name {
			//日志记录每个请求,推送所有除了日志的事件信息
			if event.Name != "logEvent" {
				engine.PutEvent(&Event{Name: "logEvent", Data: *event})
			}

			event, err := selectProcessor.EventHandler(event)
			//错误处理
			if err != nil {
				engine.PutEvent(&Event{Name: "logEvent", Data: err})
			}
			//如果返回数据再次推送
			if err == nil && event != nil {
				engine.PutEvent(event)
			}
		}
	}
}

//注册处理器
func (engine *Engine) RegisterProcessor(processor *Processor) error {

	for selectProcessorElem := engine.Processor.Front(); selectProcessorElem != nil; selectProcessorElem = selectProcessorElem.Next() {
		selectProcessor, _ := selectProcessorElem.Value.(*Processor)
		if selectProcessor.ProcessorName == processor.ProcessorName {
			engine.PutEvent(&Event{Name: "logEvent", Data: processor.ProcessorName + "注册失败"})
			return errors.New("ProcessorName has been registered")
		}
	}
	engine.Processor.PushBack(processor)
	engine.PutEvent(&Event{Name: "logEvent", Data: "注册了处理器:" + processor.ProcessorName})
	return nil
}

//注销处理器
func (engine *Engine) UnRegisterProcessor(processorName string) (bool, error) {
	for selectProcessorElem := engine.Processor.Front(); selectProcessorElem != nil; selectProcessorElem = selectProcessorElem.Next() {
		selectProcessor, _ := selectProcessorElem.Value.(*Processor)
		if selectProcessor.ProcessorName == processorName {
			engine.Processor.Remove(selectProcessorElem)
			return true, nil
		}
	}
	return false, errors.New("no this processor")
}

//定时推送定时事件
func (engine *Engine) startPutTimeEvent() {
	go func() {
		for {
			time.Sleep(engine.TimeInterval)
			engine.PutEvent(&Event{Name: "logEvent", Data: "定时器推送事件"})
			engine.PutEvent(&Event{Name: "timer", Data: time.Now()})
		}
	}()
}
