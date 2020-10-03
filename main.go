package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/bluebell/controller"

	"github.com/bluebell/pkg/snowflake"

	"github.com/bluebell/dao/redis"

	"github.com/bluebell/dao/mysql"

	"go.uber.org/zap"

	"github.com/bluebell/router"

	"github.com/bluebell/logger"

	"github.com/bluebell/settings"
)

/**
这是程序入口，我们将在这里做一系列初始化工作
*/

func main() {
	// 1.加载配置
	if err := settings.Init(); err != nil {
		log.Fatalf("settings.Init failed, err:%v\n", err)
		return
	}
	// 2.初始化Zap日志库
	if err := logger.Init(settings.Conf.LogConfig, settings.Conf.Mode); err != nil {
		log.Fatalf("init logger settings failed,err:%v\n", err)
		return
	}
	// 3.初始化MySQL连接()
	if err := mysql.Init(settings.Conf.MysqlConfig, settings.Conf.Mode); err != nil {
		zap.L().Error("mysql Init failed...", zap.Error(err))
		return
	}

	// 4.初始化Redis连接
	if err := redis.Init(settings.Conf.RedisConfig); err != nil {
		zap.L().Error("redis init failed...", zap.Error(err))
		return
	}
	// 5.初始化雪花算法的Node
	if err := snowflake.Init(settings.Conf.NodeId); err != nil {
		zap.L().Error("Init snowflake failed", zap.Error(err))
		return
	}
	//设置校验器
	if err := controller.InitTrans("zh"); err != nil {
		zap.L().Error("Init Trans failed", zap.Error(err))
		return
	}

	// 6.设置路由
	r := router.SetupRouter()

	// 7.设置优雅关机

	//优雅的关机，使用的http.Server内置的shutdown函数

	//7.1.将gin路由注册到http.Server中
	sev := http.Server{
		Addr:    settings.Conf.Port,
		Handler: r,
	}

	//7.2.开启一个goroutine启动服务
	go func() {
		err := sev.ListenAndServe()
		if err != nil && err != http.ErrServerClosed {
			log.Fatalf("Start Serve failed,err:%v\n", err)
		}
	}()

	//7.3.监听一个系统的中断信号，比如 kill kill -2
	//7.3.1首先，创建一个通道，用来接收信号，如果通道为空，程序将阻塞
	quit := make(chan os.Signal, 1)

	//7.3.2 signal.Notify会监听系统发出的 syscall.SIGINT 和 syscall.SIGTERM信号，并传入到quit中
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM) // 此处不会阻塞

	//7.3.3从通道中取值，取不到会阻塞
	<-quit

	//如果执行到这一行，说明监听到了进程退出的信号，执行退出
	zap.L().Info("Shutdown Server ...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// 5秒内优雅关闭服务（将未处理完的请求处理完再关闭服务），超过5秒就超时退出
	if err := sev.Shutdown(ctx); err != nil {
		zap.L().Error("Server Shutdown: ", zap.Error(err))
		return
	}
}
