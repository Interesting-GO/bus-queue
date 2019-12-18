/**
 * @Author: DollarKillerX
 * @Description: bus 消息总线
 * @Github: https://github.com/dollarkillerx
 * @Date: Create in 下午9:51 2019/12/18
 */
package bus_queue

import "sync"

// 事件
type DataEvent struct {
	Data  interface{} // 数据
	Topic string      // 话题
}

// 传播
type DataChannel chan DataEvent
type DataChannelSlice []DataChannel

// 事件总线
// 存储有关订阅者感兴趣的特定主题信息
type EventBus struct {
	subscribers map[string]DataChannelSlice
	rm          sync.RWMutex
}

// 注意: 这个系统整体的设计 和原来用channel传递数据的思路有点不一样
// 这里是 订阅者创建chan 注入到整个消息池中
// 发布者遍历这个消息池 进行数据发送 !!!

// 订阅主题
func (e *EventBus) Subscribe(topic string, ch DataChannel) {
	e.rm.Lock()
	defer e.rm.Unlock()
	if prev, found := e.subscribers[topic]; found {
		// 如果存在就添加
		e.subscribers[topic] = append(prev, ch)
	} else {
		// 反之初始化她
		e.subscribers[topic] = append([]DataChannel{}, ch)
	}
}

// 发布主题
func (e *EventBus) Publish(topic string, data interface{}) {
	e.rm.RLock()
	defer e.rm.RUnlock()
	if chans, found := e.subscribers[topic]; found {
		// 这样做是因为切片引用相同的数组，即使它们是按值传递的
		// 因此我们正在使用我们的元素创建一个新切片，从而能正确地保持锁定
		chanels := append(DataChannelSlice{}, chans...)
		go func(data DataEvent, dataChannelSlices DataChannelSlice) {
			for _, ch := range dataChannelSlices {
				ch <- data
			}
		}(DataEvent{data, topic}, chanels)
	}
} // 请注意，我们在发布方法中使用了 Goroutine 来避免阻塞发布者
