package mq

import (
	"strings"

	"github.com/streadway/amqp"

	"github.com/CocaineCong/micro-todoList/config"
)

var RabbitMq *amqp.Connection

// InitRabbitMQ 初始化MQ
func InitRabbitMQ() {
	connString := strings.Join([]string{config.RabbitMQ, "://", config.RabbitMQUser, ":", config.RabbitMQPassWord, "@", config.RabbitMQHost, ":", config.RabbitMQPort, "/"}, "")
	conn, err := amqp.Dial(connString)
	if err != nil {
		panic(err)
	}
	RabbitMq = conn
}
