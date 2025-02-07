package controllers

import (
	"net/http"
	"web-api/internal/api/services"
	"web-api/internal/pkg/models/request"
	"web-api/internal/pkg/models/response"

	"github.com/gin-gonic/gin"
)

type AdminSuperController struct {
	*BaseController
}

var AdminSuper = &AdminSuperController{}

func (c *AdminSuperController) GetUsers(ctx *gin.Context) {
	result, err := services.AdminSuper.GetUsersSevice()
	if err != nil {
		response.FailWithDetailed(ctx, http.StatusInternalServerError, nil, err.Error())
		return
	}
	response.OkWithData(ctx, result)
}

// func (c *AdminSuperController) UpdateAdmidsuper(ctx *gin.Context) {
// 	var requestParams request.User

// 	if err := c.ValidateReqParams(ctx, &requestParams); err != nil {
// 		response.FailWithDetailed(ctx, http.StatusBadRequest, nil, err.Error())
// 		return
// 	}
// 	result, err := services.User.UpdateUserSevice(&requestParams)
// 	if err != nil {
// 		response.FailWithDetailed(ctx, http.StatusInternalServerError, nil, err.Error())
// 		return
// 	}
// 	response.OkWithData(ctx, result)
// }

func (c *AdminSuperController) Loginadmin(ctx *gin.Context) {
	var requestParams request.AdminSuper
	if err := c.ValidateReqParams(ctx, &requestParams); err != nil {
		response.FailWithDetailed(ctx, http.StatusBadRequest, nil, err.Error())
		return
	}

	token, err := services.AdminSuper.LoginAdminSuperService(&requestParams)
	if err != nil {
		response.FailWithDetailed(ctx, http.StatusUnauthorized, nil, err.Error())
		return
	}

	response.OkWithData(ctx, gin.H{"token": token})

}
