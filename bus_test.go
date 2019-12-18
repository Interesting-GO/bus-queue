/**
 * @Author: DollarKillerX
 * @Description: 关于队列的测试
 * @Github: https://github.com/dollarkillerx
 * @Date: Create in 下午9:53 2019/12/18
 */
package bus_queue

import (
	"fmt"
	"math/rand"
	"testing"
	"time"
)

// 单例
var eb = &EventBus{
	subscribers: map[string]DataChannelSlice{},
}

// 测试定时发送数据
func pubTo(topic string, data string) {
	rand.Seed(time.Now().UnixNano())
	for {
		eb.Publish(topic, data)
		time.Sleep(time.Duration(rand.Intn(1000)) * time.Millisecond)
	}
}

func printDataEvent(ch string, data DataEvent) {
	fmt.Printf("Channel:%s;Topic:%s;DataEvent:%v\n", ch, data.Topic, data.Data)
}

func TestMan(t *testing.T) {
	ch1 := make(chan DataEvent)
	ch2 := make(chan DataEvent)
	ch3 := make(chan DataEvent)

	eb.Subscribe("topic1", ch1)
	eb.Subscribe("topic2", ch2)
	eb.Subscribe("topic3", ch3)

	go pubTo("topic1", "Hi topic 1")
	go pubTo("topic2", "Hello topic 2")

	for {
		select {
		case d := <-ch1:
			go printDataEvent("ch1", d)
		case d := <-ch2:
			go printDataEvent("ch2", d)
		case d := <-ch3:
			go printDataEvent("ch3", d)
		}
	}
}
