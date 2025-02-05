package controllers

import (
	"net/http"
	"web-api/internal/api/services"
	"web-api/internal/pkg/models/response"

	"github.com/gin-gonic/gin"
)

type BaseseController struct {
	*BaseController
}

var Basese = &BaseseController{}

func (c *BaseseController) Getbasese(ctx *gin.Context) {
	result, err := services.Order.BasesSevice()
	if err != nil {
		response.FailWithDetailed(ctx, http.StatusInternalServerError, nil, err.Error())
		return
	}
	response.OkWithData(ctx, result)
}
