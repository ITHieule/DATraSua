package controllers

import (
	"net/http"
	"web-api/internal/api/services"
	"web-api/internal/pkg/models/request"
	"web-api/internal/pkg/models/response"

	"github.com/gin-gonic/gin"
)

type StatisticalController struct {
	*BaseController
}

var Statistical = &StatisticalController{}

func (c *StatisticalController) GetStatistical(ctx *gin.Context) {
	var requestParams request.Statisticalrequest

	// Validate request params
	if err := c.ValidateReqParams(ctx, &requestParams); err != nil {
		response.FailWithDetailed(ctx, http.StatusBadRequest, nil, err.Error())
		return
	}
	// Gọi service với tham số year và month
	result, err := services.Statis.GetStatisticalSevice(&requestParams)
	if err != nil {
		response.FailWithDetailed(ctx, http.StatusInternalServerError, nil, err.Error())
		return
	}

	// Trả về kết quả thành công
	response.OkWithData(ctx, result)
}
