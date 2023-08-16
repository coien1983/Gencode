package template

var MainTemplate = `package main

import (
	"%s/docs"
	"%s/routers"
	"%s/sysinit"
	"context"
	"flag"
	"fmt"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"go.uber.org/zap"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

var configFile = flag.String("f", "conf/config.yaml", "the config file")

var spec []byte

func main() {
	flag.Parse()

	// 1。加载配置
	if err := sysinit.SettingInit(configFile); err != nil {
		fmt.Printf("init settings failed,err:%v\n", err)
		return
	}

	// 2。初始化日志
	if err := sysinit.LogInit(sysinit.Conf.LogC); err != nil {
		fmt.Printf("init logger failed,err:%v\n", err)
		return
	}

	defer zap.L().Sync() //将缓冲区的日志追加到日志中
	zap.L().Debug("logger init success ...")
	// 3。初始化mysql连接
	if err := sysinit.MysqlInit(sysinit.Conf); err != nil {
		fmt.Printf("init mysql failed,err:%v\n", err)
		return
	}

	//5.初始化 redis
	if err := sysinit.RedisInit(sysinit.Conf); err != nil {
		fmt.Printf("init redis failed,err:%v\n", err)
		return
	}

	//6。初始化gin框架内部校验器
	if err := sysinit.InitTrans("zh"); err != nil {
		fmt.Printf("init validator trans failed,err :%v\n", err)
		return
	}

	// 7。注册路由
	r := router.RouterInit()
	if sysinit.Conf.Mode == "test" || sysinit.Conf.Mode == "dev" {

		docs.SwaggerInfo.Title = "%s Api"
		docs.SwaggerInfo.Description = "This is %s Api"
		docs.SwaggerInfo.Version = "1.0"
		docs.SwaggerInfo.Host = fmt.Sprintf("%s:%d", "localhost", 7050)
		docs.SwaggerInfo.BasePath = "/v1"
		r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	}

	// 8。启动服务（优雅关机）
	srv := &http.Server{
		Addr:    fmt.Sprintf(":%d", sysinit.Conf.Port),
		Handler: r,
	}

	go func() {
		//开启一个goroutine启动服务
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			zap.L().Fatal("listen:", zap.Error(err))
		}
	}()

	// 等待终端信号来优雅的关闭服务器，为关闭服务器 设置一个5秒的超时
	quit := make(chan os.Signal, 1) //创建一个收信号的通道
	//kill 默认会发送 syscall.SIGTERM 信号
	//kill -2 发送   syscall.SIGINT  信号，我们常用的ctrl+c 就是触发系统的SIGINT信号
	//kill -9 发送   syscall.SIGKILL 信号，但是不能被捕获，所以不需要添加它
	//signal.Notify把收到的syscall.SIGINT或者syscall.SIGTERM信号转发给quit
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit //阻塞在此，当接收到上述梁总信号时才会往下执行
	//创建一个5米秒超时的context
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	// 5秒内优雅的关闭服务(将未处理完的请求处理完再关闭服务，超过5秒就超时退出)
	if err := srv.Shutdown(ctx); err != nil {
		zap.L().Fatal("Server Shutdown:", zap.Error(err))
	}

	zap.L().Info("Server existing")
}
`
