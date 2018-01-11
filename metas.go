package BtcQuant

type EventInterface interface {
}

//事件结构体
type Event struct {
	Name string      //事件名字
	Data interface{} //事件数据
}

//事件处理器
type Processor struct {
	ProcessorName string       //处理器名字
	EventName     string       //处理事件类型的名字
	EventHandler  func(*Event) //处理事件的函数
}
