package controllers

import (
	"net/http"
	"web-api/internal/api/services"
	"web-api/internal/pkg/models/response"

	"github.com/gin-gonic/gin"
)

type SizesController struct {
	*BaseController
}

var Sizes = &SizesController{}

func (c *SizesController) GetSizes(ctx *gin.Context) {
	result, err := services.OrderSize.SizesSevice()
	if err != nil {
		response.FailWithDetailed(ctx, http.StatusInternalServerError, nil, err.Error())
		return
	}
	response.OkWithData(ctx, result)
}
