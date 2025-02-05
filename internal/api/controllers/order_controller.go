package controllers

import (
	"net/http"
	"web-api/internal/api/services"
	"web-api/internal/pkg/models/request"
	"web-api/internal/pkg/models/response"

	"github.com/gin-gonic/gin"
)

type OrderController struct {
	*BaseController
}

var Order = &OrderController{}

func (c *OrderController) Register(ctx *gin.Context) {
	var requestParams request.User
	if err := c.ValidateReqParams(ctx, &requestParams); err != nil {
		response.FailWithDetailed(ctx, http.StatusBadRequest, nil, err.Error())
		return
	}
	result, err := services.User.Register(&requestParams)
	if err != nil {
		response.FailWithDetailed(ctx, http.StatusInternalServerError, nil, err.Error())
	}
	response.OkWithData(ctx, result)
}

func (c *OrderController) Login(ctx *gin.Context) {
	var requestParams request.User
	if err := c.ValidateReqParams(ctx, &requestParams); err != nil {
		response.FailWithDetailed(ctx, http.StatusBadRequest, nil, err.Error())
		return
	}
	token, err := services.User.Login(&requestParams)
	if err != nil {
		response.FailWithDetailed(ctx, http.StatusInternalServerError, nil, err.Error())
	}
	response.OkWithData(ctx, token)

}
