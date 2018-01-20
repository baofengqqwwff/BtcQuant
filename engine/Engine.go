package engine

import (
	"container/list"
	"github.com/pkg/errors"
	"log"
	"time"
)

type Engine struct {
	Event        *list.List    //事件
	Processor    *list.List    //事件处理器
	TimeInterval time.Duration //时间间隔，默认一秒
}

func NewEngine() *Engine {
	engine := &Engine{list.New(), list.New(), time.Second}
	engine.engineStart()
	engine.startPutTimeEvent()
	return engine
}

//启动引擎
func (engine *Engine) engineStart() {
	log.Println("引擎启动")
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
			selectProcessor.EventHandler(event)
		}
	}
}

//注册处理器
func (engine *Engine) RegisterProcessor(processor *Processor) {
	log.Println("注册了处理器")
	engine.Processor.PushBack(processor)
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
			log.Println("定时器推送事件")
			engine.PutEvent(&Event{Name: "Timer", Data: time.Now()})
		}
	}()
}
