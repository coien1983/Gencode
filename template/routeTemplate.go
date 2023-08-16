package template

var RouteTemplate = `package router

import (
	api "%s/controllers"
	"%s/middlewares"
	"%s/sysinit"
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
