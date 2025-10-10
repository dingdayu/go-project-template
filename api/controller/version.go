package controller

import (
	"net/http"

	"github.com/dingdayu/singbox-adapter/model/entity"

	"github.com/gin-gonic/gin"
)

// Version
// @Summary 获取接口版本
// @Produce  json
// @SuccessResponse 200 {object} api.Response
// @Failure 500 {object} api.Response
// @Router /api/v1/version [get]
func Version(c *gin.Context) {
	res := map[string]interface{}{}
	res["code"] = 200
	res["message"] = "success"
	res["data"] = map[string]string{
		"time":    entity.BuildTime,
		"version": entity.BuildVersion,
	}
	c.JSON(http.StatusOK, res)
}
