package controllers

import (
	"net/http"
	"web-api/internal/api/services"
	"web-api/internal/pkg/models/request"
	"web-api/internal/pkg/models/response"

	"github.com/gin-gonic/gin"
)

type BasesController struct {
	*BaseController
}

var Bases = &BasesController{}

func (c *BasesController) Getbases(ctx *gin.Context) {
	var requestParams request.User
	if err := c.ValidateReqParams(ctx, &requestParams); err != nil {
		response.FailWithDetailed(ctx, http.StatusBadRequest, nil, err.Error())
		return
	}
	result, err := services.Order.BasesSevice()
	if err != nil {
		response.FailWithDetailed(ctx, http.StatusInternalServerError, nil, err.Error())
	}
	response.OkWithData(ctx, result)
}
