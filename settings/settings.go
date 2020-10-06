package settings

/**
基于viper做的日志管理
*/

import (
	"fmt"

	"github.com/fsnotify/fsnotify"

	"github.com/spf13/viper"
)

var Conf = new(AppConfig) //提供一个全局的变量供外部调用

/**
配置项信息可以序列化到此结构体对象中
*/
type AppConfig struct {
	Name         string                 `mapstructure:"name"`           //应用名称
	Version      string                 `mapstructure:"version"`        //应用版本
	Mode         string                 `mapstructure:"mode"`           //当前运行模式 dev 开发 test 测试 release 线上
	NodeId       int64                  `mapstructure:"node_id"`        //雪花算法当前节点ID
	Port         string                 `mapstructure:"port"`           //启动端口
	Md5Secret    string                 `mapstructure:"md5Secret"`      //md5加密算法的秘钥
	TokenSecret  string                 `mapstructure:"tokenSecret"`    //token秘钥
	SingleSignOn bool                   `mapstructure:"single_sign_on"` //是否限制单点登录
	*MysqlConfig `mapstructure:"mysql"` //mysql配置
	*LogConfig   `mapstructure:"log"`   //log配置
	*RedisConfig `mapstructure:"redis"` //redis配置
}

type MysqlConfig struct {
	User         string `mapstructure:"user"`           //mysql用户名
	Password     string `mapstructure:"password"`       //mysql密码
	Host         string `mapstructure:"host"`           //mysql主机
	Port         string `mapstructure:"port"`           //mysql端口
	DbName       string `mapstructure:"dbname"`         //mysql数据库名
	MaxIdleConns int    `mapstructure:"max_idle_conns"` //空闲连接池中的连接数
	MaxOpenConns int    `mapstructure:"max_open_conns"` //最大打开的数据库连接数量
}

type LogConfig struct {
	Level      string `mapstructure:"level"`
	Filename   string `mapstructure:"filename"`
	MaxSize    int    `mapstructure:"max_size"`
	MaxBackups int    `mapstructure:"max_backups"`
	MaxAge     int    `mapstructure:"max_age"`
}

type RedisConfig struct {
	Addr     string `mapstructure:"addr"`
	Password string `mapstructure:"password"`
	DB       int    `mapstructure:"db"`
}

func Init() (err error) {
	//1.指定Viper需要加载的配置文件

	//方式1：指定配置文件名和配置文件位置（相对路径）
	viper.SetConfigName("config")
	viper.AddConfigPath("./conf")

	//方式2：直接指定配置文件路径
	//viper.SetConfigFile("./conf/config.yaml")

	//2.读取配置
	if err = viper.ReadInConfig(); err != nil {

		//判断是否是因为找不到配置文件的错误
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			fmt.Printf("config file not found,err:%v\n", err)
		} else {
			//其他错误
			fmt.Printf("viper read in config failed,err:%v\n", err)
		}
		return
	}
	// 3.反序列化
	if err = viper.Unmarshal(Conf); err != nil {
		fmt.Printf("viper unmarshal failed,err:%v\n", err)
		return
	}
	// 4.监听配置文件改变
	viper.WatchConfig()

	// 5.设置当配置文件改变的回调
	viper.OnConfigChange(func(e fsnotify.Event) {
		fmt.Println("注意：配置文件已变更...")
		if err = viper.Unmarshal(Conf); err != nil {
			fmt.Printf("viper unmarshal failed,err:%v\n", err)
			return
		}
	})
	return
}
