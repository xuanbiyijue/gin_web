package script

import (
	"context"

	"github.com/CocaineCong/micro-todoList/app/task/repository/mq/task"
	log "github.com/CocaineCong/micro-todoList/pkg/logger"
)

// TaskCreateSync 执行MQ到MySQL的任务
func TaskCreateSync(ctx context.Context) {
	tSync := new(task.SyncTask)
	// 开一个线程，从MQ中拿到数据存到数据库
	err := tSync.RunTaskCreate(ctx)
	if err != nil {
		log.LogrusObj.Infof("RunTaskCreate:%s", err)
	}
}
