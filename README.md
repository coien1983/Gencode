# 项目名称

这个项目主要是为了快速生成一个gin框架开发的项目结构

## 功能特点

项目结构包括config,common,controller,models,docs,enums,jobs,hooks,middlewares,pkg,request,responses
routers,services,sysinit,main函数,Dockfile
同时包括单机部署打包脚本,build.sh,支持pm2守护运行的run.sh和pm2.yml,

## 快速开始

git clone git@github.com:coien1983/Gencode.git 到Go语言的src目录下

进行项目目录，执行go install，此时会在golang的bin目录下有一个Gencode可执行文件

## 使用示例

这条命令可以在当前目录下生成一个名为demo的工程项目，执行目录为8080
Gencode -p Demo:8080

这条目录，可以在当前目录的controller下生成一个DemoController.go的文件，并生成与之配套的service文件
Gencode -c Demo


