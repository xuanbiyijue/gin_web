/*
管理配置的第三方库:    viper
配置信息文件:         config.yaml
 */

package config

import (
	"fmt"
	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
)


// AppConfig 主要配置
type AppConfig struct {
	Name       string `mapstructure:"name"`
	Mode       string `mapstructure:"mode"`
	Port       int    `mapstructure:"port"`
	StartTime  string `mapstructure:"start_time"`
	MachineID  int64  `mapstructure:"machine_id"`

	*LogConfig   `mapstructure:"log"`
	*MySQLConfig `mapstructure:"mysql"`
	*RedisConfig `mapstructure:"redis"`
}

// LogConfig 日志配置
type LogConfig struct {
	Level      string `mapstructure:"level"`
	Filename   string `mapstructure:"filename"`
	MaxSize    int    `mapstructure:"max_size"`
	MaxAge     int    `mapstructure:"max_age"`
	MaxBackups int    `mapstructure:"max_backups"`
}

// MySQLConfig mysql配置
type MySQLConfig struct {
	Host         string `mapstructure:"host"`
	User         string `mapstructure:"user"`
	Password     string `mapstructure:"password"`
	DbName       string `mapstructure:"db_name"`
	Port         int    `mapstructure:"port"`
	MaxOpenConns int    `mapstructure:"max_open_conns"`
	MaxIdleConns int    `mapstructure:"max_idle_conns"`
}

// RedisConfig redis配置
type RedisConfig struct {
	Host     string `mapstructure:"host"`
	Password string `mapstructure:"password"`
	Port     int    `mapstructure:"port"`
	DB       int    `mapstructure:"db"`
	PoolSize int    `mapstructure:"pool_size"`
}


// 实例化配置信息
var Config = new(AppConfig)


// Init 初始化配置
func Init(file_path string) (err error) {
	// 获取配置文件
	viper.SetConfigFile(file_path)  // 指定配置文件路径
	err = viper.ReadInConfig()      // 查找并读取配置文件
	if err != nil {
		panic(fmt.Errorf("Fatal error config file: %s \n", err))
		return
	}

	// 把读取到的信息反序列化到 Config 变量中
	err = viper.Unmarshal(Config)
	if err != nil {
		panic(fmt.Errorf("viper.Unmarshal failed, err: %v", err))
		return
	}

	// 监视配置信息，以防发生变化
	viper.WatchConfig()
	viper.OnConfigChange(func(e fsnotify.Event) {
		fmt.Println("config had been updated!")
		// 再次反序列化
		if err := viper.Unmarshal(Config); err != nil {
			panic(fmt.Errorf("viper.Unmarshal failed, err: %v", err))
			return
		}
	})
	return
}

