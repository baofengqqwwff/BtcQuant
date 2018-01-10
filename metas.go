package BtcQuant

//事件结构体
type Event struct {
	Name string      //事件名字
	Data interface{} //事件数据
}

//事件处理器
type Processor struct {

}