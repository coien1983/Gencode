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
	"wanhe_code/mygithub/Gencode/template"
)

func main() {
	//项目名称
	projectName := flag.String("p", "", "项目名称")
	controllerName := flag.String("c", "", "控制器名称")
	serviceName := flag.String("s", "", "service名称")
	repositoryName := flag.String("r", "", "repository名称")
	testName := flag.String("t", "", "测试名称")
	flag.Parse()

	if *projectName != "" {
		fmt.Println("项目名称:", *projectName)
		//项目初始化
		ProjectInit(*projectName)
		return
	} else {
		//fmt.Println("未提供项目名称")
	}

	if *controllerName != "" {
		fmt.Println("控制器名称:", *controllerName)
		//创建控制器
		ControllerInit(*controllerName)
		return
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

	if *testName != "" {
		fmt.Println("测试名称", *testName)
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

	// 要检测的目录路径
	dirPath := "controller"

	// 使用Stat函数获取目录信息
	_, err = os.Stat(dirPath)

	// 判断目录是否存在
	if os.IsNotExist(err) {
		//如果文件夹不存在，创建
		// 创建目录
		err = os.Mkdir(dirPath, 0755)
		if err != nil {
			fmt.Println("创建目录失败：", err)
			return
		}
	} else if err != nil {
		fmt.Println("发生错误：", err)
		return
	} else {

	}

	projectName := filepath.Base(dir)
	// 移除换行符
	controllerName = strcase.ToCamel(strings.ToLower(controllerName))

	//判断控制器文件是否存在
	isExist, err := IsFileExist(fmt.Sprintf("controller/%sController.go", controllerName))
	if err != nil {
		fmt.Println(err)
		return
	}

	if isExist {
		fmt.Println("目标目录下已经存在控制器")
		return
	}

	controllerContent := template.ControllerTemplate

	contentArr := []any{projectName, projectName, projectName, projectName, controllerName, controllerName, controllerName, controllerName, controllerName, controllerName, controllerName}
	controllerContent = fmt.Sprintf(controllerContent, contentArr...)

	BaseControllerName := controllerName
	controllerName = fmt.Sprintf("controller/%sController.go", controllerName)
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
		responseContent := template.ResponseTemplate
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

	isResExist, err := IsFileExist("services")
	if err != nil {
		fmt.Println(err)
		return
	}

	if !isResExist {
		err := os.Mkdir("services", 0755)
		if err != nil {
			fmt.Println("Failed to create directory:", err)
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

		serviceContent := template.ServiceTemplate

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
	content := template.BaseResponseTemplate
	err := os.WriteFile("responses/BaseResp.go", []byte(content), 0644)
	if err != nil {
		fmt.Println("BaseResp生成失败")
		return errors.New("BaseResp生成失败")
	}

	return nil
}

func InitBuildDockerFile(projectName, portName string) error {
	content := template.DockerTemplate
	content = fmt.Sprintf(content, projectName, portName, projectName)

	err := os.WriteFile("Dockerfile", []byte(content), 0644)
	if err != nil {
		fmt.Println("Dockerfile生成失败")
		return errors.New("Dockerfile生成失败")
	}

	return nil
}

func InitBuildSh(projectName string) error {
	content := template.BuildTemplate
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
	content := template.Pm2Template
	content = fmt.Sprintf(content, projectName, projectName)

	err := os.WriteFile("pm2.yaml", []byte(content), 0644)
	if err != nil {
		fmt.Println("pm2.yaml配置失败")
		return errors.New("pm2.yaml配置失败")
	}

	return nil
}

func InitRunSh(projectName string) error {
	content := template.RunTemplate

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
	content := template.ReadMeTemplate
	content = fmt.Sprintf(content, projectName)

	err := os.WriteFile("README.md", []byte(content), 0644)
	if err != nil {
		fmt.Println("配置主函数文件失败")
		return errors.New("配置主函数文件失败")
	}

	return nil
}

func InitMain(projectName string) error {
	content := template.MainTemplate
	content = fmt.Sprintf(content, projectName, projectName, projectName, projectName, projectName)

	err := os.WriteFile("main.go", []byte(content), 0644)
	if err != nil {
		fmt.Println("配置主函数文件失败")
		return errors.New("配置主函数文件失败")
	}

	return nil
}

func InitRecoverMiddleWare() error {
	content := template.RecoveryTemplate

	err := os.WriteFile("middlewares/recover_middleware.go", []byte(content), 0644)
	if err != nil {
		fmt.Println("配置跨域中间件失败")
		return errors.New("配置跨域中间件初始化失败")
	}

	return nil
}

func InitCorsMiddleWare() error {
	content := template.CorsTemplate

	err := os.WriteFile("middlewares/cors_middleware.go", []byte(content), 0644)
	if err != nil {
		fmt.Println("配置跨域中间件失败")
		return errors.New("配置跨域中间件初始化失败")
	}

	return nil
}

// InitRouterInit 路由配置初始化
func InitRouterInit(projectName string) error {
	content := template.RouteTemplate
	content = fmt.Sprintf(content, projectName, projectName, projectName)
	err := os.WriteFile("routers/router.go", []byte(content), 0644)
	if err != nil {
		fmt.Println("路由配置初始化失败")
		return errors.New("路由配置初始化失败")
	}

	return nil
}

func InitRedisInit() error {
	content := template.RedisTemplate

	err := os.WriteFile("sysinit/redis_init.go", []byte(content), 0644)
	if err != nil {
		fmt.Println("redis初始化失败")
		return errors.New("redis初始化失败")
	}

	return nil
}

// InitMysqlInit 数据库初始化
func InitMysqlInit() error {
	content := template.MysqlTemplate
	err := os.WriteFile("sysinit/mysql_init.go", []byte(content), 0644)
	if err != nil {
		fmt.Println("数据库初始化失败")
		return errors.New("数据库初始化失败")
	}

	return nil
}

// InitLogInit 配置日志初始化
func InitLogInit() error {
	content := template.LogTemplate

	err := os.WriteFile("sysinit/log_init.go", []byte(content), 0644)
	if err != nil {
		fmt.Println("配置日志初始化失败")
		return errors.New("配置日志初始化失败")
	}

	return nil
}

// InitSettingInit 配置初始化文件
func InitSettingInit() error {
	content := template.SettingTemplate

	err := os.WriteFile("sysinit/setting_init.go", []byte(content), 0644)
	if err != nil {
		fmt.Println("配置初始化文件失败")
		return errors.New("配置初始化文件失败")
	}

	return nil
}

// InitSettingFile 初始化配置模型
func InitSettingFile() error {
	content := template.SettingModelTemplate

	err := os.WriteFile("sysinit/setting_conf.go", []byte(content), 0644)
	if err != nil {
		fmt.Println("配置模型初始化失败")
		return errors.New("配置模型初始化失败")
	}

	return nil
}

func InitConfigFile(projectName, portName string) error {
	templateD := template.ConfigTemplate
	content := fmt.Sprintf(templateD, "test", projectName, portName, projectName)

	err := os.WriteFile("config/config.yaml", []byte(content), 0644)
	if err != nil {
		fmt.Println("配置文件初始化失败")
		return errors.New("配置文件初始化失败")
	}

	content = fmt.Sprintf(templateD, "dev", projectName, portName, projectName)
	err = os.WriteFile("config/config-dev.yaml", []byte(content), 0644)
	if err != nil {
		fmt.Println("配置文件初始化失败")
		return errors.New("配置文件初始化失败")
	}

	content = fmt.Sprintf(templateD, "online", projectName, portName, projectName)
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
		//"github.com/onsi/ginkgo/v2",
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
		"test",
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
	err = InitRouterInit(projectName)
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
	err = InitMain(projectName)
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
