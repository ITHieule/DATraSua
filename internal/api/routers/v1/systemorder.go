package router_v1

import (
	"web-api/internal/api/controllers"

	"github.com/gin-gonic/gin"
)

func RegisterOrderSystemRouter(router *gin.RouterGroup) {

	//router.GET("/getbook", controllers.Book.GetBook)
	//router.GET("/oderstatsbook", controllers.Book.Oderstat)
	//router.POST("/Addbook", controllers.Book.AddBook)
	//router.DELETE("/Deletebook", controllers.Book.DeleteBook)
	//router.PUT("/Updatebook", controllers.Book.UpdateBook)
	//router.POST("/Oderbook", controllers.Book.OderBook)
	//router.POST("/Searchbook", controllers.Book.SearchBook)
	router.POST("/register", controllers.Order.Register)
	router.POST("/Login", controllers.Order.Login)
}
