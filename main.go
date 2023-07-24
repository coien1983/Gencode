package main

import (
	"bufio"
	"errors"
	"flag"
	"fmt"
	"github.com/iancoleman/strcase"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

func main() {
	//项目名称
	projectName := flag.String("p", "", "项目名称")
	controllerName := flag.String("c", "", "控制器名称")
	serviceName := flag.String("s", "", "service名称")
	repositoryName := flag.String("r", "", "repository名称")
	flag.Parse()

	if *projectName != "" {
		fmt.Println("项目名称:", *projectName)
		//项目初始化
		ProjectInit(*projectName)
	} else {
		//fmt.Println("未提供项目名称")
	}

	if *controllerName != "" {
		fmt.Println("控制器名称:", *controllerName)
		//创建控制器
		ControllerInit(*controllerName)
	} else {
		//fmt.Println("未提供controller名称")
	}

	if *serviceName != "" {
		fmt.Println("service名称:", *serviceName)
	} else {
		//fmt.Println("未提供service名称")
	}

	if *repositoryName != "" {
		fmt.Println("repository名称:", *repositoryName)
	} else {
		//fmt.Println("未提供repository名称")
	}
}

func IsFileExist(filePath string) (bool, error) {
	// 使用 os.Stat() 获取文件信息
	_, err := os.Stat(filePath)

	if err == nil {
		// 文件存在
		//fmt.Println("文件存在")
		return true, nil
	} else if os.IsNotExist(err) {
		// 文件不存在
		return false, nil
	} else {
		// 其他错误
		return false, errors.New("内部错误")
	}
}

func ControllerInit(controllerName string) {
	dir, err := os.Getwd()
	if err != nil {
		fmt.Println(err)
		return
	}

	projectName := filepath.Base(dir)
	// 移除换行符
	controllerName = strcase.ToCamel(strings.ToLower(controllerName))

	//判断控制器文件是否存在
	isExist, err := IsFileExist(fmt.Sprintf("controllers/%sController.go", controllerName))
	if err != nil {
		fmt.Println(err)
		return
	}

	if isExist {
		fmt.Println("目标目录下已经存在控制器")
		return
	}

	controllerContent := `
	"%s/requests"
	"%s/responses"
	"%s/services"
	"%s/sysinit"
	"fmt"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

type %sController struct {
}

func %sRegister(group *gin.RouterGroup) {
	%sController := &%sController{}
	group.Use()
	{
		//我的代理
		group.POST("/demo", %sController.DemoHandler)
	}
}

// DemoHandler
// @Tags Demo
// @Title Demo
// @Summary Demo
// @Description Demo
// @Accept json
// @Param object body requests.DemoReq true "demo"
// @Produce json
// @Success 000000 {object} responses.ResponseData{response=responses.DemoResp}
// @Failure 100001 {object} responses.ResponseData
// @router /v1/api/%s/demo [post]
func (c *%sController) DemoHandler(ctx *gin.Context) {
	var params requests.DemoReq
	err := ctx.ShouldBindJSON(&params)
	if err != nil {
		responses.ResponseError(ctx, "100001", err.Error())
		return
	}

	demoService := services.NewDemoService(sysinit.GetDB())

	data, err := demoService.Demo(ctx)
	if err != nil {
		responses.ResponseError(ctx, "100001", err.Error())
		return
	}

	responses.ResponseSuccess(ctx, data)
}
`
	contentArr := []any{projectName, projectName, projectName, projectName, controllerName, controllerName, controllerName, controllerName, controllerName, controllerName, controllerName}
	controllerContent = fmt.Sprintf(controllerContent, contentArr...)

	BaseControllerName := controllerName
	controllerName = fmt.Sprintf("%sController.go", controllerName)
	err = os.WriteFile(controllerName, []byte(controllerContent), 0644)
	if err != nil {
		fmt.Println("控制器生成失败")
		return
	}

	//生成response文件

	//判断控制器文件是否存在
	isExist, err = IsFileExist(fmt.Sprintf("responses/%sResp.go", controllerName))
	if err != nil {
		fmt.Println(err)
		return
	}

	if isExist {
		//fmt.Println("目标目录下已经存在响应文件")
		//return
	} else {
		responseContent := "package responses\n\ntype DemoResp struct {\n\tToken string `json:\"token\"` //登录token\n}"
		responseName := fmt.Sprintf("responses/%sResp.go", BaseControllerName)

		isResExist, err := IsFileExist("responses")
		if err != nil {
			fmt.Println(err)
			return
		}

		if !isResExist {
			err := os.Mkdir("responses", 0755)
			if err != nil {
				fmt.Println("Failed to create directory:", err)
				return
			}
		}

		err = os.WriteFile(responseName, []byte(responseContent), 0644)
		if err != nil {
			fmt.Println(err)
			fmt.Println("response生成失败")
			return
		}
	}

	//生成service文件
	isExist, err = IsFileExist(fmt.Sprintf("services/%sService.go", controllerName))
	if err != nil {
		fmt.Println(err)
		return
	}

	if isExist {
	} else {

		serviceContent := `package services

import (
	"%s/requests"
	"%s/responses"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"time"
)

type IDemoService interface {
	Demo(ctx *gin.Context, params requests.DemoReq) (*responses.DemoResp)
}

type DemoService struct {
	DB *gorm.DB
}

func NewDemoService(db *gorm.DB) IDemoService {
	return &DemoService{DB: db}
}
`

		serviceContent = fmt.Sprintf(serviceContent, projectName, projectName)
		serviceName := fmt.Sprintf("services/%sService.go", BaseControllerName)
		err = os.WriteFile(serviceName, []byte(serviceContent), 0644)
		if err != nil {
			fmt.Println(err)
			fmt.Println("service生成失败")
			return
		}
	}

	return
}

func InitBaseResponse() error {
	content := "package responses\n\nimport (\n\t\"github.com/gin-gonic/gin\"\n\t\"net/http\"\n)\n\ntype ResponseData struct {\n\tCode string      `json:\"code\" example:\"10000\"` //响应码\n\tMsg  string      `json:\"msg\" example:\"操作成功\"`   //响应信息\n\tData interface{} `json:\"data,omitempty\"`\n}\n\n// ResponseError 错误响应\nfunc ResponseError(c *gin.Context, code string, message string) {\n\tc.JSON(http.StatusOK, &ResponseData{\n\t\tCode: code,\n\t\tMsg:  message,\n\t\tData: nil,\n\t})\n}\n\nfunc ResponseSuccess(c *gin.Context, data interface{}) {\n\tc.JSON(http.StatusOK, &ResponseData{\n\t\tCode: \"000000\",\n\t\tMsg:  \"操作成功\",\n\t\tData: data,\n\t})\n}\n\n"
	err := os.WriteFile("responses/base_response.go", []byte(content), 0644)
	if err != nil {
		fmt.Println("base_response生成失败")
		return errors.New("base_response生成失败")
	}

	return nil
}

func InitBuildDockerFile(projectName, portName string) error {
	content := "FROM golang:latest\n\nWORKDIR /app\n\nCOPY . .\n\nRUN go build -o %s main.go\n\nEXPOSE %s\n\nCMD [\"./%s\"]"
	content = fmt.Sprintf(content, projectName, portName, projectName)

	err := os.WriteFile("Dockerfile", []byte(content), 0644)
	if err != nil {
		fmt.Println("Dockerfile生成失败")
		return errors.New("Dockerfile生成失败")
	}

	return nil
}

func InitBuildSh(projectName string) error {
	content := `rm -rf ./app
rm -rf %s.tar.gz
mkdir app
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags "-s -w" -o ./%s ./main.go
chmod +x ./%s
cp %s ./app/
cp -R ./conf ./app/
cp -R ./docs ./app/
cp run.sh ./app/
cp pm2.yml ./app/
tar -zcvf %s.tar.gz ./app
`
	content = fmt.Sprintf(content, projectName, projectName, projectName, projectName, projectName)

	err := os.WriteFile("build.sh", []byte(content), 0644)
	if err != nil {
		fmt.Println("build.sh配置失败")
		return errors.New("build.sh配置失败")
	}

	err = os.Chmod("build.sh", 0755)
	if err != nil {
		fmt.Println("build.sh配置失败：", err)
		return errors.New("build.sh配置失败")
	}

	return nil
}

func InitPm2(projectName string) error {
	content := "apps:\n  - name: %s\n    instances: 1\n    exec_mode: fork\n    interpreter: \"./%s\"\n    interpreter_args: \"\"\n    script: \"./pm2.yml\""

	content = fmt.Sprintf(content, projectName, projectName)

	err := os.WriteFile("pm2.yaml", []byte(content), 0644)
	if err != nil {
		fmt.Println("pm2.yaml配置失败")
		return errors.New("pm2.yaml配置失败")
	}

	return nil
}

func InitRunSh(projectName string) error {
	content := "#ps -ef|grep %s|grep -v grep|awk '{print $2}'|xargs kill\npid=`ps -ef|grep %s|grep -v grep|awk '{print $2}'`\nif [ -n \"$pid\" ]; then\n    echo \"process id is: $pid\"\n    kill $pid\n    echo 'process killed'\nelse\n    echo \"%s process id is not found!\"\nfi\nexport GIN_MODE=release\necho 'starting...'\n#nohup ./%s &\npm2 start pm2.yml\necho 'started!'"

	content = fmt.Sprintf(content, projectName, projectName, projectName, projectName)

	err := os.WriteFile("run.sh", []byte(content), 0644)
	if err != nil {
		fmt.Println("run.sh配置失败")
		return errors.New("run.sh配置失败")
	}

	err = os.Chmod("run.sh", 0755)
	if err != nil {
		fmt.Println("run.sh配置失败：", err)
		return errors.New("run.sh配置失败")
	}

	return nil
}

func InitReadMe(projectName string) error {
	content := `# %s

config: stores configuration functions
common: stores helper functions
docs: stores Swagger documentation
codes: stores error codes
enums: stores constant codes
hooks: stores asynchronous hooks
middleware: stores middleware
controller: stores controllers
model: stores database models
pkg: stores third-party packages
repository: stores database interaction layer
requests: stores request models
responses: stores response models
routers: stores routers
sysinit: stores system initialization
services: stores service layer
jobs: stores timer jobs
`

	content = fmt.Sprintf(content, projectName)

	err := os.WriteFile("README.md", []byte(content), 0644)
	if err != nil {
		fmt.Println("配置主函数文件失败")
		return errors.New("配置主函数文件失败")
	}

	return nil
}

func InitMain() error {
	content := `package main

import (
	"$project_name/docs"
	"$project_name/routers"
	"$project_name/sysinit"
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
	if err := sysinit.LogInit(sysinit.Conf.LogConfig); err != nil {
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

		docs.SwaggerInfo.Title = "agent_shop Api"
		docs.SwaggerInfo.Description = "This is JiHe Api"
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

	err := os.WriteFile("main.go", []byte(content), 0644)
	if err != nil {
		fmt.Println("配置主函数文件失败")
		return errors.New("配置主函数文件失败")
	}

	return nil
}

func InitRecoverMiddleWare() error {
	content := `package middlewares

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"net"
	"net/http"
	"net/http/httputil"
	"os"
	"runtime/debug"
	"strings"
)

func GinRecovery() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				var brokenPipe bool
				if ne, ok := err.(*net.OpError); ok {
					if se, ok := ne.Err.(*os.SyscallError); ok {
						if strings.Contains(strings.ToLower(se.Error()), "broken pipe") || strings.Contains(strings.ToLower(se.Error()), "connection reset by peer") {
							brokenPipe = true
						}
					}
				}

				httpRequest, _ := httputil.DumpRequest(c.Request, false)
				if brokenPipe {
					zap.L().Error(c.Request.URL.Path,
						zap.Any("error", err),
						zap.String("request", string(httpRequest)),
						zap.String("stack", string(debug.Stack())),
					)
				} else {
					zap.L().Error("[Recovery from panic]",
						zap.Any("error", err),
						zap.String("request", string(httpRequest)),
					)
				}

				c.AbortWithStatus(http.StatusInternalServerError)
			}
		}()

		c.Next()
	}
}
`

	err := os.WriteFile("middlewares/recover_middleware.go", []byte(content), 0644)
	if err != nil {
		fmt.Println("配置跨域中间件失败")
		return errors.New("配置跨域中间件初始化失败")
	}

	return nil
}

func InitCorsMiddleWare() error {
	content := `package middlewares

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func CORS() gin.HandlerFunc {

	return func(context *gin.Context) {

		method := context.Request.Method

		context.Header("Access-Control-Allow-Origin", "*")
		context.Header("Access-Control-Allow-Headers", "Content-Type,AccessToken,X-CSRF-Token, Authorization, Token, x-token")
		context.Header("Access-Control-Allow-Methods", "POST, GET, OPTIONS, DELETE, PATCH, PUT")
		context.Header("Access-Control-Expose-Headers", "Content-Length, Access-Control-Allow-Origin, Access-Control-Allow-Headers, Content-Type")
		context.Header("Access-Control-Allow-Credentials", "true")

		if method == "OPTIONS" {
			context.AbortWithStatus(http.StatusNoContent)
		}

		context.Next()
	}
}
`

	err := os.WriteFile("middlewares/cors_middleware.go", []byte(content), 0644)
	if err != nil {
		fmt.Println("配置跨域中间件失败")
		return errors.New("配置跨域中间件初始化失败")
	}

	return nil
}

// InitRouterInit 路由配置初始化
func InitRouterInit() error {
	content := `package router

import (
	api "$project_name/controllers"
	"$project_name/middlewares"
	"$project_name/sysinit"
	"github.com/gin-gonic/gin"
)

func RouterInit() *gin.Engine {

	ro := gin.New()
	if sysinit.Conf.Mode == "release" {
		ro.Use(middlewares.GinRecovery(), middlewares.CORS())
	} else {
		ro.Use(gin.Logger(), middlewares.GinRecovery(), middlewares.CORS())
	}

	rootPath := ro.Group("/v1")
	{
		AuthPath := rootPath.Group("/api")
		{
			api.DemoRegister(AuthPath)
		}
	}

	return ro
}
`

	err := os.WriteFile("routers/router.go", []byte(content), 0644)
	if err != nil {
		fmt.Println("路由配置初始化失败")
		return errors.New("路由配置初始化失败")
	}

	return nil
}

func InitRedisInit() error {
	content := `package sysinit

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/go-redis/redis"
	"time"
)

var rdb *redis.Client

func RedisInit(setting *AppConfig) (err error) {

	redisConf := setting.RedisC

	rdb = redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", redisConf.Host, redisConf.Port),
		Password: redisConf.Password,
		DB:       redisConf.Db,
		PoolSize: redisConf.PoolSize,
	})

	_, err = rdb.Ping().Result()
	return

}

func GetRedis() *redis.Client {
	return rdb
}

func RedisClose() {
	_ = rdb.Close()
}

//计算Redis的剩余时间
func RedisTtl(key string) (time.Duration, error) {
	prefix := Conf.RedisC.Prefix
	keyR := fmt.Sprintf("%s%s", prefix, key)
	return rdb.TTL(keyR).Result()
}

func RedisIsExists(key string) (int64, error) {
	prefix := Conf.RedisC.Prefix
	keyR := fmt.Sprintf("%s%s", prefix, key)
	return rdb.Exists(keyR).Result()
}

// RedisReadString 读取字符串
func RedisReadString(key string) (string, error) {
	prefix := Conf.RedisC.Prefix
	keyR := fmt.Sprintf("%s%s", prefix, key)

	return rdb.Get(keyR).Result()
}

// RedisWriteString 写入字符串
func RedisWriteString(key string, value interface{}, expiredTime int64) error {
	prefix := Conf.RedisC.Prefix
	keyR := fmt.Sprintf("%s%s", prefix, key)

	newTime := time.Duration(expiredTime) * time.Second
	return rdb.Set(keyR, value, newTime).Err()
}

// RedisReadStruct 读取结构体
func RedisReadStruct(key string, obj interface{}) error {

	if data, err := RedisReadString(key); err == nil {
		return json.Unmarshal([]byte(data), obj)
	} else {
		return err
	}
}

// RedisWriteStruct 写入结构体
func RedisWriteStruct(key string, obj interface{}, expiredTime int64) error {
	data, err := json.Marshal(obj)
	if err == nil {
		return RedisWriteString(key, string(data), expiredTime)
	} else {
		return err
	}
}

// RedisDelete 删除键
func RedisDelete(key string) error {
	prefix := Conf.RedisC.Prefix
	keyR := fmt.Sprintf("%s%s", prefix, key)
	return rdb.Del(keyR).Err()
}

// RedisPopQueue 从队列中获取数据
func RedisPopQueue(key string) (number string, err error) {
	prefix := Conf.RedisC.Prefix
	keyR := fmt.Sprintf("%s%s", prefix, key)
	val, err := rdb.RPop(keyR).Result()
	if err != nil {
		return "", err
	} else {
		return val, nil
	}
}

// RedisPushQueue 推送数据入队列
func RedisPushQueue(key string, value interface{}) (err error) {
	prefix := Conf.RedisC.Prefix
	keyR := fmt.Sprintf("%s%s", prefix, key)
	n, err := rdb.LPush(keyR, value).Result()
	if err != nil {
		return err
	}

	if n < 1 {
		return errors.New("入列失败")
	}

	return nil
}

func RedisIncr(key string) (err error) {
	prefix := Conf.RedisC.Prefix
	keyR := fmt.Sprintf("%s%s", prefix, key)
	n, err := rdb.Incr(keyR).Result()
	if err != nil {
		return err
	}

	if n < 1 {
		return errors.New("操作失败")
	}

	return nil
}

func RedisDecr(key string) (err error) {
	prefix := Conf.RedisC.Prefix
	keyR := fmt.Sprintf("%s%s", prefix, key)
	n, err := rdb.Decr(keyR).Result()
	if err != nil {
		return err
	}

	if n < 1 {
		return errors.New("操作失败")
	}

	return nil
}

func RedisHSet(key string, field, data string, dayTime time.Time) error {
	prefix := Conf.RedisC.Prefix
	keyR := fmt.Sprintf("%s%s", prefix, key)

	err := rdb.HSet(keyR, field, data).Val()
	if !err {
		return errors.New("操作失败")
	}

	_, _ = rdb.ExpireAt(keyR, dayTime).Result()

	return nil
}

func RedisHGet(key, data string) (string, error) {
	prefix := Conf.RedisC.Prefix
	keyR := fmt.Sprintf("%s%s", prefix, key)
	data, err := rdb.HGet(keyR, data).Result()
	if err != nil {
		return "", err
	}

	return data, nil
}

func RedisHIncrBy(key string, data string, expiredTime int64) error {
	prefix := Conf.RedisC.Prefix
	keyR := fmt.Sprintf("%s%s", prefix, key)
	n, err := rdb.HIncrBy(keyR, data, expiredTime).Result()
	if err != nil {
		return err
	}

	if n < 1 {
		return errors.New("操作失败")
	}

	return nil
}

func RedisExpiredKey(key string, t time.Time) {
	prefix := Conf.RedisC.Prefix
	keyR := fmt.Sprintf("%s%s", prefix, key)

	_, _ = rdb.ExpireAt(keyR, t).Result()
}
`

	err := os.WriteFile("sysinit/redis_init.go", []byte(content), 0644)
	if err != nil {
		fmt.Println("redis初始化失败")
		return errors.New("redis初始化失败")
	}

	return nil
}

// InitMysqlInit 数据库初始化
func InitMysqlInit() error {
	content := `package sysinit

import (
	"fmt"
	"go.uber.org/zap"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
	"log"
	"os"
	"time"
)

var db *gorm.DB

func MysqlInit(setting *AppConfig) (err error) {

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		setting.MysqlC.User,
		setting.MysqlC.Password,
		setting.MysqlC.Host,
		setting.MysqlC.Port,
		setting.MysqlC.Dbname,
	)

	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer（日志输出的目标，前缀和日志包含的内容——译者注）
		logger.Config{
			SlowThreshold:             time.Second, // 慢 SQL 阈值
			LogLevel:                  logger.Info, // 日志级别
			IgnoreRecordNotFoundError: true,        // 忽略ErrRecordNotFound（记录未找到）错误
			Colorful:                  true,        // 禁用彩色打印
		},
	)

	if setting.Mode == "dev" {
		db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{
			NamingStrategy: schema.NamingStrategy{
				SingularTable: true,
			},
			Logger: newLogger,
		})
	} else {
		db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{
			NamingStrategy: schema.NamingStrategy{
				SingularTable: true,
			},
			Logger: newLogger,
		})
	}

	if err != nil {
		if setting.Mode == "dev" {
			fmt.Println(err)
		}
		zap.L().Error("connect DB failed", zap.Error(err))
		return
	}

	sqlDB, err := db.DB()
	if err != nil {
		if setting.Mode == "dev" {
			fmt.Println(err)
		}
		zap.L().Error("db.DB() failed", zap.Error(err))
	}

	// SetMaxIdle 设置空闲连接池中连接的最大数量
	sqlDB.SetMaxIdleConns(setting.MysqlC.MaxIdle)

	// SetMaxOpen 设置打开数据库的最大数量
	sqlDB.SetMaxOpenConns(setting.MysqlC.MaxOpen)

	return
}

func GetDB() *gorm.DB {
	return db
}

func CloseDB() {
	sqlDB, err := db.DB()
	if err != nil {
		zap.L().Error("db.DB() failed", zap.Error(err))
	}
	_ = sqlDB.Close()
}
`

	err := os.WriteFile("sysinit/mysql_init.go", []byte(content), 0644)
	if err != nil {
		fmt.Println("数据库初始化失败")
		return errors.New("数据库初始化失败")
	}

	return nil
}

// InitLogInit 配置日志初始化
func InitLogInit() error {
	content := `package sysinit

import (
	"github.com/natefinch/lumberjack"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func LogInit(logger *LogC) (err error) {
	writerSyncer := getLogWriter(
		logger.Filename,
		logger.MaxSize,
		logger.MaxBackups,
		logger.MaxAge,
	)
	encoder := getEncoder()

	var l = new(zapcore.Level)
	err = l.UnmarshalText([]byte(logger.Level))
	if err != nil {
		return
	}

	core := zapcore.NewCore(encoder, writerSyncer, l)
	lg := zap.New(core, zap.AddCaller())
	//替换zap库中全局的logger
	zap.ReplaceGlobals(lg)
	return
}

func getEncoder() zapcore.Encoder {
	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	encoderConfig.TimeKey = "time"
	encoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder
	encoderConfig.EncodeDuration = zapcore.SecondsDurationEncoder
	encoderConfig.EncodeCaller = zapcore.ShortCallerEncoder
	return zapcore.NewJSONEncoder(encoderConfig)
}

func getLogWriter(filename string, maxSize, maxBackup, maxAge int) zapcore.WriteSyncer {
	lumberJackLogger := &lumberjack.Logger{
		Filename:   filename,
		MaxSize:    maxSize,
		MaxAge:     maxAge,
		MaxBackups: maxBackup,
	}

	return zapcore.AddSync(lumberJackLogger)
}
`

	err := os.WriteFile("sysinit/log_init.go", []byte(content), 0644)
	if err != nil {
		fmt.Println("配置日志初始化失败")
		return errors.New("配置日志初始化失败")
	}

	return nil
}

// InitSettingInit 配置初始化文件
func InitSettingInit() error {
	content := `package sysinit

import (
	"fmt"
	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
)

// Conf 全局变量，用来保存程序的所有配置
var Conf = new(AppConfig)

func SettingInit(fileName *string) (err error) {
	viper.SetConfigFile(*fileName)
	viper.AddConfigPath(".")
	err = viper.ReadInConfig()
	if err != nil {
		fmt.Printf("viper.ReadInConfig() failed,err:%v\n", err)
		return
	}

	if err = viper.Unmarshal(Conf); err != nil {
		fmt.Printf("viper.Unmarshal failed,err:%v\n", err)
		return
	}

	viper.WatchConfig()
	viper.OnConfigChange(func(in fsnotify.Event) {
		fmt.Println("配置文件修改了...")
		if err = viper.Unmarshal(Conf); err != nil {
			fmt.Printf("viper.Unmarshal failed,err:%v\n", err)
		}
	})

	return
}
`

	err := os.WriteFile("sysinit/setting_init.go", []byte(content), 0644)
	if err != nil {
		fmt.Println("配置初始化文件失败")
		return errors.New("配置初始化文件失败")
	}

	return nil
}

// InitSettingFile 初始化配置模型
func InitSettingFile() error {
	content := "package sysinit\n\ntype AppConfig struct {\n\tName    string `mapstructure:\"name\"`\n\tMode    string `mapstructure:\"mode\"`\n\tVersion string `mapstructure:\"version\"`\n\tPort    int    `mapstructure:\"port\"`\n\t*LogC\n\t*MysqlC\n\t*RedisC\n\t*WechatC\n\t*QiNiuC\n}\n\ntype LogC struct {\n\tLevel      string `mapstructure:\"level\"`\n\tFilename   string `mapstructure:\"filename\"`\n\tMaxSize    int    `mapstructure:\"max_size\"`\n\tMaxAge     int    `mapstructure:\"max_age\"`\n\tMaxBackups int    `mapstructure:\"max_backups\"`\n}\n\ntype MysqlC struct {\n\tHost         string `mapstructure:\"host\"`\n\tUser         string `mapstructure:\"user\"`\n\tPassword     string `mapstructure:\"password\"`\n\tDbname       string `mapstructure:\"dbname\"`\n\tPort         int    `mapstructure:\"port\"`\n\tMaxOpen int    `mapstructure:\"max_open\"`\n\tMaxIdle int    `mapstructure:\"max_idle\"`\n}\n\ntype RedisC struct {\n\tHost     string `mapstructure:\"host\"`\n\tPassword string `mapstructure:\"password\"`\n\tPort     int    `mapstructure:\"port\"`\n\tDb       int    `mapstructure:\"db\"`\n\tPoolSize int    `mapstructure:\"pool_size\"`\n\tPrefix   string `mapstructure:\"prefix\"`\n}\n\ntype WechatC struct {\n\tAppID     string `mapstructure:\"app_id\"`\n\tSecret    string `mapstructure:\"secret\"`\n\tMchId     string `mapstructure:\"mch_id\"`\n\tAppKey    string `mapstructure:\"app_key\"`\n\tNotifyUrl string `mapstructure:\"notify_url\"`\n\tPrefix    string `mapstructure:\"prefix\"`\n}\n\ntype QiNiuC struct {\n\tSecretKey string `mapstructure:\"secret_key\"`\n\tAccessKey string `mapstructure:\"access_key\"`\n\tPicDomain string `mapstructure:\"pic_domain\"`\n\tBucket    string `mapstructure:\"bucket\"`\n}"

	err := os.WriteFile("sysinit/setting_conf.go", []byte(content), 0644)
	if err != nil {
		fmt.Println("配置模型初始化失败")
		return errors.New("配置模型初始化失败")
	}

	return nil
}

func InitConfigFile(projectName, portName string) error {
	content := `name: "%s"
mode: "test"
version: "0.0.1"
port: %s

# 日志配置
LogC:
  level: "info"
  filename: "$project_name.log"
  max_size: 200
  max_age: 30
  max_backups: 7

# 微信配置
WechatC:
  app_id: ""
  secret: ""
  mch_id: ""
  app_key: ""
  notify_url: ""
  prefix: ""

# 数据库配置
MysqlC:
  host: ""
  port:
  user: ""
  password: ""
  dbname: ""
  max_open: 200
  max_idle: 50

# redis配置
RedisC:
  host: 127.0.0.1
  port: 6379
  password:
  db: 0
  pool_size: 100
  prefix: ""

# 七牛配置
QiNiuC:
  secret_key: ""
  access_key: ""
  pic_domain: ""
  bucket: ""
`
	content = fmt.Sprintf(content, projectName, portName)

	err := os.WriteFile("config/config.yaml", []byte(content), 0644)
	if err != nil {
		fmt.Println("配置文件初始化失败")
		return errors.New("配置文件初始化失败")
	}

	err = os.WriteFile("config/config-dev.yaml", []byte(content), 0644)
	if err != nil {
		fmt.Println("配置文件初始化失败")
		return errors.New("配置文件初始化失败")
	}

	err = os.WriteFile("config/config-online.yaml", []byte(content), 0644)
	if err != nil {
		fmt.Println("配置文件初始化失败")
		return errors.New("配置文件初始化失败")
	}

	return nil
}

// InitDir 初始化目录文件
func InitDir(projectName string) error {

	err := os.Mkdir(projectName, 0755)
	if err != nil {
		fmt.Println("Failed to create directory:", err)
		return errors.New("Failed to create directory")
	}

	err = os.Chdir(projectName)
	if err != nil {
		fmt.Println("Failed to change directory:", err)
		return errors.New("Failed to change directory")
	}

	// 初始化Go模块
	cmd := exec.Command("go", "mod", "init", projectName)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Run()

	// 引入依赖项
	dependencies := []string{
		"github.com/gin-gonic/gin",
		"github.com/fsnotify/fsnotify",
		"github.com/spf13/viper",
		"gorm.io/driver/mysql",
		"gorm.io/gorm",
		"github.com/go-redis/redis",
		"github.com/natefinch/lumberjack",
		"go.uber.org/zap",
		"go.uber.org/zap/zapcore",
	}

	for _, dependency := range dependencies {
		cmd2 := exec.Command("go", "get", "-u", dependency)
		cmd2.Stdout = os.Stdout
		cmd2.Stderr = os.Stderr
		cmd2.Run()
	}

	// 引入gin-swagger支持
	swaggerDependencies := []string{
		"github.com/swaggo/gin-swagger",
		"github.com/swaggo/files",
	}

	for _, dependency := range swaggerDependencies {
		cmd3 := exec.Command("go", "get", "-u", dependency)
		cmd3.Stdout = os.Stdout
		cmd3.Stderr = os.Stderr
		cmd3.Run()
	}

	// 创建基本目录结构
	directories := []string{
		"config",
		"common",
		"docs",
		"errcodes",
		"enums",
		"hooks",
		"middlewares",
		"controller",
		"model",
		"pkg",
		"repository",
		"requests",
		"responses",
		"routers",
		"sysinit",
		"services",
		"jobs",
	}

	for _, dir := range directories {
		os.Mkdir(dir, 0755)
	}

	return nil
}

func ProjectInit(projectDe string) {
	//初始化项目设置
	projectD := strings.Split(projectDe, ":")
	if len(projectD) != 2 {
		fmt.Println("创建项目工程参数有误")
		return
	}

	projectName := projectD[0]
	portName := projectD[1]

	err := InitDir(projectName)
	if err != nil {
		return
	}

	//初始化配置文件
	err = InitConfigFile(projectName, portName)
	if err != nil {
		return
	}

	//生成配置模型
	err = InitSettingFile()
	if err != nil {
		return
	}

	//生成配置初始化文件
	err = InitSettingInit()
	if err != nil {
		return
	}

	//配置日志初始化
	err = InitLogInit()
	if err != nil {
		return
	}

	//mysql初始化
	err = InitMysqlInit()
	if err != nil {
		return
	}

	//redis初始化
	err = InitRedisInit()
	if err != nil {
		return
	}

	//路由初始化
	err = InitRouterInit()
	if err != nil {
		return
	}

	//跨域中间件初始化
	err = InitCorsMiddleWare()
	if err != nil {
		return
	}

	//生成recover中间件
	err = InitRecoverMiddleWare()
	if err != nil {
		return
	}

	//基础response
	err = InitBaseResponse()
	if err != nil {
		return
	}

	//生成主函数
	err = InitMain()
	if err != nil {
		return
	}

	//生成readme文件
	err = InitReadMe(projectName)
	if err != nil {
		return
	}

	//生成run.sh脚本
	err = InitRunSh(projectName)
	if err != nil {
		return
	}

	//配置pm2文件
	err = InitPm2(projectName)
	if err != nil {
		return
	}

	//配置打包脚本
	err = InitBuildSh(projectName)
	if err != nil {
		return
	}

	err = InitBuildDockerFile(projectName, portName)
	if err != nil {
		return
	}
}

// InitProjectSetting 初始化项目设置
func InitProjectSetting() (string, string, string, error) {
	//项目名称
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("请输入项目名称: ")
	projectName, err := reader.ReadString('\n')
	if err != nil {
		fmt.Println("无法读取输入:", err)
		return "", "", "", errors.New("无法读取输入")
	}
	// 移除换行符
	projectName = projectName[:len(projectName)-1]

	//作者名称
	fmt.Print("请输入项目作者: ")
	authorName, err := reader.ReadString('\n')
	if err != nil {
		fmt.Println("无法读取输入:", err)
		return "", "", "", errors.New("无法读取输入")
	}
	authorName = authorName[:len(authorName)-1]

	//运行端口
	fmt.Print("请输入项目运行端口: ")
	portName, err := reader.ReadString('\n')
	if err != nil {
		fmt.Println("无法读取输入:", err)
		return "", "", "", errors.New("无法读取输入")
	}
	portName = portName[:len(portName)-1]

	return projectName, authorName, portName, nil
}
