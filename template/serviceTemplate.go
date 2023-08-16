package template

var ServiceTemplate = `package services

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

func (d *DemoService) Demo(ctx *gin.Context, params requests.DemoReq) (*responses.DemoResp) {
	//TODO you can write your service code here
	panic("do something")
}

func NewDemoService(db *gorm.DB) IDemoService {
	return &DemoService{DB: db}
}
`
