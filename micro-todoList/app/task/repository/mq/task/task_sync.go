package task

import (
	"context"
	"encoding/json"

	"github.com/CocaineCong/micro-todoList/app/task/repository/mq"
	"github.com/CocaineCong/micro-todoList/app/task/service"
	"github.com/CocaineCong/micro-todoList/consts"
	"github.com/CocaineCong/micro-todoList/idl/pb"
	log "github.com/CocaineCong/micro-todoList/pkg/logger"
)

type SyncTask struct {
}

// RunTaskCreate MQ到mysql
func (s *SyncTask) RunTaskCreate(ctx context.Context) error {
	// 先获得消息队列的名字
	rabbitMqQueue := consts.RabbitMqTaskQueue
	// msgs 是存放消息的 chan
	msgs, err := mq.ConsumeMessage(ctx, rabbitMqQueue)
	if err != nil {
		return err
	}

	// 阻塞一个 chan ，使下面的线程一直执行
	var forever chan struct{}
	// 这里开一个线程，就是让 MQ 里的数据能够一直存到数据库中
	go func() {
		// 开始循环取出信息
		for d := range msgs {
			log.LogrusObj.Infof("Received run Task: %s", d.Body)

			// 落库
			reqRabbitMQ := new(pb.TaskRequest)
			// 反序列化到 reqRabbitMQ
			err = json.Unmarshal(d.Body, reqRabbitMQ)
			if err != nil {
				log.LogrusObj.Infof("Received run Task: %s", err)
			}
			// 存到 MySQL
			err = service.TaskMQ2MySQL(ctx, reqRabbitMQ)
			if err != nil {
				log.LogrusObj.Infof("Received run Task: %s", err)
			}
			d.Ack(false)
		}
	}()
	log.LogrusObj.Infoln(err)
	<-forever

	return nil
}
