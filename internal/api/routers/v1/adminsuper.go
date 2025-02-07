package router_v1

import (
	"web-api/internal/api/controllers"

	"github.com/gin-gonic/gin"
)

// FILE - ADMIN ROUTER
func RegisterAdminsSystemRouter(router *gin.RouterGroup) {

	//api - Login
	router.POST("/loginadmin", controllers.AdminSuper.Loginadmin)

	router.GET("/GetUsers", controllers.AdminSuper.GetUsers)

	router.PUT("/Updateadmin", controllers.AdminSuper.UpdateAdmidsuper)

	router.DELETE("/Deleteadmin", controllers.AdminSuper.DeleteAdmidsuper)

}
