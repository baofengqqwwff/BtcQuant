package BtcQuant

type Engine struct {
	Event     []interface{} //事件
	Processor []interface{} //事件处理器
}

func (engine *Engine) processEvent() {

}
func (engine *Engine) RegisterProcessor(processor *interface{}){

}
func (engine *Engine) UnRegisterProcessor(processor *interface{}){

}