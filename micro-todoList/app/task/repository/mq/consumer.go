package mq

import (
	"context"

	"github.com/streadway/amqp"
)

// ConsumeMessage MQ到mysql，传递消息
func ConsumeMessage(ctx context.Context, queueName string) (msgs <-chan amqp.Delivery, err error) {
	// 打开 MQ 的 Channel
	ch, err := RabbitMq.Channel()
	if err != nil {
		return
	}
	// 创建一个消息队列
	q, _ := ch.QueueDeclare(queueName, true, false, false, false, nil)
	// 设置一次取几条数据
	err = ch.Qos(1, 0, false)
	// 开始传递消息
	return ch.Consume(q.Name, "", false, false, false, false, nil)
}
