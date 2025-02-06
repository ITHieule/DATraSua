package controllers

import (
	"fmt"
	"net/http"
	"strconv"
	"web-api/internal/api/services"
	"web-api/internal/pkg/models/request"

	"github.com/gin-gonic/gin"
)

// CartController cấu trúc cho controller giỏ hàng
type CartController struct{}

// NewCartController tạo mới CartController
func NewCartController() *CartController {
	return &CartController{}
}

// GetCart trả về giỏ hàng của người dùng
func (ctrl *CartController) GetCart(c *gin.Context) {
	// Lấy userID từ URL parameter
	userIDStr := c.Param("userID")
	userID, err := strconv.Atoi(userIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	// Gọi service để lấy giỏ hàng
	cart, err := services.Cart.GetCart(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Error retrieving cart: %v", err)})
		return
	}

	// Trả về kết quả giỏ hàng
	c.JSON(http.StatusOK, gin.H{"cart": cart})
}

// AddToCart thêm một mặt hàng vào giỏ hàng
func (ctrl *CartController) AddToCart(c *gin.Context) {
	var requestParams request.Cartrequest

	// Parse dữ liệu từ body request
	if err := c.ShouldBindJSON(&requestParams); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input data"})
		return
	}

	// Lấy userID từ URL parameter
	userIDStr := c.Param("userID")
	userID, err := strconv.Atoi(userIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	// Gọi service để thêm item vào giỏ hàng
	cartItem, err := services.Cart.AddToCart(userID, &requestParams)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Error adding item to cart: %v", err)})
		return
	}

	// Trả về kết quả giỏ hàng sau khi thêm
	c.JSON(http.StatusOK, gin.H{"cart_item": cartItem})
}

// UpdateCart cập nhật thông tin giỏ hàng
func (ctrl *CartController) UpdateCart(c *gin.Context) {
	var requestParams request.Cartrequest

	// Parse dữ liệu từ body request
	if err := c.ShouldBindJSON(&requestParams); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input data"})
		return
	}

	// Lấy userID và cartItemID từ URL parameters
	userIDStr := c.Param("userID")
	userID, err := strconv.Atoi(userIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	cartItemIDStr := c.Param("cartItemID")
	cartItemID, err := strconv.Atoi(cartItemIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid cart item ID"})
		return
	}

	// Cập nhật ID của cartItem
	requestParams.ID = cartItemID

	// Gọi service để cập nhật giỏ hàng
	updatedCartItem, err := services.Cart.UpdateCart(userID, &requestParams)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Error updating cart: %v", err)})
		return
	}

	// Trả về kết quả giỏ hàng sau khi cập nhật
	c.JSON(http.StatusOK, gin.H{"updated_cart_item": updatedCartItem})
}

// RemoveFromCart xóa một mục khỏi giỏ hàng
func (ctrl *CartController) RemoveFromCart(c *gin.Context) {
	// Lấy userID và cartItemID từ URL parameters
	userIDStr := c.Param("userID")
	userID, err := strconv.Atoi(userIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	cartItemIDStr := c.Param("cartItemID")
	cartItemID, err := strconv.Atoi(cartItemIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid cart item ID"})
		return
	}

	// Gọi service để xóa mục khỏi giỏ hàng
	err = services.Cart.RemoveFromCart(userID, cartItemID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Error removing item from cart: %v", err)})
		return
	}

	// Trả về thông báo thành công
	c.JSON(http.StatusOK, gin.H{"message": "Item removed from cart successfully"})
}
