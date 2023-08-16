package template

var ResponseTemplate = "package responses\n\ntype DemoResp struct {\n\tToken string `json:\"token\"` //登录token\n}"

var BaseResponseTemplate = "package responses\n\nimport (\n\t\"github.com/gin-gonic/gin\"\n\t\"net/http\"\n)\n\ntype ResponseData struct {\n\tCode string      `json:\"code\" example:\"10000\"` //响应码\n\tMsg  string      `json:\"msg\" example:\"操作成功\"`   //响应信息\n\tData interface{} `json:\"data,omitempty\"`\n}\n\n// ResponseError 错误响应\nfunc ResponseError(c *gin.Context, code string, message string) {\n\tc.JSON(http.StatusOK, &ResponseData{\n\t\tCode: code,\n\t\tMsg:  message,\n\t\tData: nil,\n\t})\n}\n\nfunc ResponseSuccess(c *gin.Context, data interface{}) {\n\tc.JSON(http.StatusOK, &ResponseData{\n\t\tCode: \"000000\",\n\t\tMsg:  \"操作成功\",\n\t\tData: data,\n\t})\n}\n\n"
