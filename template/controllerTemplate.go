package template

// ControllerTemplate 控制器模版
var ControllerTemplate = `
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
		//demo
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
