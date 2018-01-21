package engine

//事件处理器
type Processor struct {
	ProcessorName string       //处理器名字
	EventName     string       //处理事件类型的名字
	EventHandler  func(*Event) (*Event,error)//处理事件的函数
}
