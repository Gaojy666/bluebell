package settings

import (
	"fmt"

	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
)

// 全局变量，用来保存程序的所有配置信息
var Conf = new(AppConfig)

type AppConfig struct {
	Name         string `mapstructure:"name"`
	Mode         string `mapstructure:"mode"`
	Version      string `mapstructure:"version"`
	Port         int    `mapstructure:"port"`
	StartTime    string `mapstructure:"start_time"`
	MachineID    int64  `mapstructure:"machine_id"`
	*LogConfig   `mapstructure:"log"`
	*MySQLConfig `mapstructure:"mysql"`
	*RedisConfig `mapstructure:"redis"`
}

type LogConfig struct {
	Level      string `mapstructure:"level"`
	Filename   string `mapstructure:"filename"`
	MaxSize    int    `mapstructure:"max_size"`
	MaxAge     int    `mapstructure:"max_age"`
	MaxBackups int    `mapstructure:"max_backups"`
}

type MySQLConfig struct {
	Host         string `mapstructure:"host"`
	User         string `mapstructure:"user"`
	Password     string `mapstructure:"password"`
	DbName       string `mapstructure:"dbname"`
	Port         int    `mapstructure:"port"`
	MaxOpenConn  int    `mapstructure:"max_open_conns""`
	MaxIdleConns int    `mapstructure:"max_idle_conns"`
}

type RedisConfig struct {
	Host     string `mapstructure:"host"`
	Password string `mapstructure:"password"`
	Port     int    `mapstructure:"port"`
	DB       int    `mapstructure:"db"`
	PoolSize int    `mapstructure:"pool_size"`
}

// 利用viper管理配置文件
func Init() (err error) {
	viper.SetConfigFile("./conf/config.yaml")

	//viper.SetConfigFile(filepath)

	err = viper.ReadInConfig() // 读取配置信息
	if err != nil {
		// 读取配置信息失败
		fmt.Printf(" viper.ReadInConfig() failed, err:%v\n", err)
		return
	}

	// 把读到的配置信息反序列化到Conf变量中
	if err := viper.Unmarshal(Conf); err != nil {
		fmt.Printf("viper.Unmarshal failed, err:#{err}\n")
	}
	// 监控配置文件变化
	viper.WatchConfig()
	viper.OnConfigChange(func(in fsnotify.Event) {
		// fsnotify监控文件变化的库
		fmt.Println("配置文件修改了...")
		if err := viper.Unmarshal(Conf); err != nil {
			fmt.Printf("vipe r.Unmarshal failed, err:%v\n", err)
		}
	})
	return

}
