package wrappers

import (
	"context"
	"strconv"

	"github.com/afex/hystrix-go/hystrix"
	"go-micro.dev/v4/client"

	"github.com/CocaineCong/micro-todoList/idl/pb"
)

// NewTask 创建一个任务模型
func NewTask(id uint64, name string) *pb.TaskModel {
	return &pb.TaskModel{
		Id:         id,
		Title:      name,
		Content:    "响应超时",
		StartTime:  1000,
		EndTime:    1000,
		Status:     0,
		CreateTime: 1000,
		UpdateTime: 1000,
	}
}

// DefaultTasks 降级函数，如果发生降级，则不再进行rpc通信，而是生成10条仿照任务样式本地生成的数据
func DefaultTasks(resp interface{}) {
	models := make([]*pb.TaskModel, 0)
	// 生成 10 条假数据
	var i uint64
	for i = 0; i < 10; i++ {
		models = append(models, NewTask(i, "降级备忘录"+strconv.Itoa(20+int(i))))
	}
	// 验证响应数据是不是任务清单类型的数据，并将生成的假数据放进响应数据里
	result := resp.(*pb.TaskListResponse)
	result.TaskList = models
}

type TaskWrapper struct {
	client.Client
}

func (wrapper *TaskWrapper) Call(ctx context.Context, req client.Request, rsp interface{}, opts ...client.CallOption) error {
	cmdName := req.Service() + "." + req.Endpoint()
	// 配置信息
	config := hystrix.CommandConfig{
		Timeout:                3000,
		RequestVolumeThreshold: 20,   // 熔断器请求阈值，默认20，意思是有20个请求才能进行错误百分比计算
		ErrorPercentThreshold:  50,   // 错误百分比，当错误超过百分比时，直接进行降级处理，直至熔断器再次 开启，默认50%
		SleepWindow:            5000, // 过多长时间，熔断器再次检测是否开启，单位毫秒ms（默认5秒）
	}
	hystrix.ConfigureCommand(cmdName, config)
	// hystrix.Do() 如果程序正常则执行第一个函数，否则执行第二个函数
	return hystrix.Do(cmdName, func() error {
		return wrapper.Client.Call(ctx, req, rsp)
	}, func(err error) error {
		DefaultTasks(rsp)
		return err
	})
}

// NewProductWrapper 初始化Wrapper
func NewTaskWrapper(c client.Client) client.Client {
	return &TaskWrapper{c}
}
