package config

import (
	"fmt"

	"gopkg.in/ini.v1"
)

// 或者考虑使用结构体来保存配置信息
var (
	Db         string
	DbHost     string
	DbPort     string
	DbUser     string
	DbPassWord string
	DbName     string
	Charset    string

	RabbitMQ         string
	RabbitMQUser     string
	RabbitMQPassWord string
	RabbitMQHost     string
	RabbitMQPort     string

	EtcdHost string
	EtcdPort string

	UserServiceAddress string
	TaskServiceAddress string
)

// Init 初始化配置
func Init() {
	// 加载配置文件
	file, err := ini.Load("./config/config.ini")
	if err != nil {
		fmt.Println("配置文件读取错误，请检查文件路径:", err)
	}
	// 分别加载 MySQL、ETCD、Rabbit MQ、微服务配置
	LoadMysqlData(file)
	LoadEtcd(file)
	LoadRabbitMQ(file)
	LoadServer(file)
}

// LoadMysqlData 加载 MySQL 配置
func LoadMysqlData(file *ini.File) {
	Db = file.Section("mysql").Key("Db").String()
	DbHost = file.Section("mysql").Key("DbHost").String()
	DbPort = file.Section("mysql").Key("DbPort").String()
	DbUser = file.Section("mysql").Key("DbUser").String()
	DbPassWord = file.Section("mysql").Key("DbPassWord").String()
	DbName = file.Section("mysql").Key("DbName").String()
	Charset = file.Section("mysql").Key("Charset").String()
}

// LoadRabbitMQ 加载 RabbitMQ 配置
func LoadRabbitMQ(file *ini.File) {
	RabbitMQ = file.Section("rabbitmq").Key("RabbitMQ").String()
	RabbitMQUser = file.Section("rabbitmq").Key("RabbitMQUser").String()
	RabbitMQPassWord = file.Section("rabbitmq").Key("RabbitMQPassWord").String()
	RabbitMQHost = file.Section("rabbitmq").Key("RabbitMQHost").String()
	RabbitMQPort = file.Section("rabbitmq").Key("RabbitMQPort").String()
}

// LoadEtcd 加载 Etcd 配置
func LoadEtcd(file *ini.File) {
	EtcdHost = file.Section("etcd").Key("EtcdHost").String()
	EtcdPort = file.Section("etcd").Key("EtcdPort").String()
}

// LoadServer 加载服务器配置，这是针对用户微服务和任务微服务的
func LoadServer(file *ini.File) {
	UserServiceAddress = file.Section("server").Key("UserServiceAddress").String()
	TaskServiceAddress = file.Section("server").Key("TaskServiceAddress").String()
}
