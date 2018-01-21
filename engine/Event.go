package engine

/*事件类型
  newBar  新的bar
  newTick 新的tick
  apiEvent 运用到api的事件，包括下单，取消单，获取tick等等
  logEvent 日志服务
*/
//事件结构体
type Event struct {
	Name string      //事件名字
	Data interface{} //事件数据
}
