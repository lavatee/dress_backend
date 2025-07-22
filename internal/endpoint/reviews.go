package endpoint

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/lavatee/dresscode_backend/internal/model"
)

type CreateReviewInput struct {
	ProductID int    `json:"product_id" binding:"required"`
	Rating    int    `json:"rating" binding:"required" validate:"required,min=1,max=5"`
	Comment   string `json:"comment" binding:"required"`
}

func (e *Endpoint) CreateReview(c *gin.Context) {
	var input CreateReviewInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, ErrorResponse{Message: err.Error()})
		return
	}
	userId, err := e.GetUserId(c)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, ErrorResponse{Message: err.Error()})
		return
	}
	review := model.Review{
		ProductID: input.ProductID,
		UserID:    userId,
		Rating:    input.Rating,
		Comment:   input.Comment,
	}
	id, err := e.services.Reviews.CreateReview(review)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, ErrorResponse{Message: err.Error()})
		return
	}
	c.JSON(http.StatusOK, map[string]interface{}{"id": id})
}

func (e *Endpoint) GetReviews(c *gin.Context) {
	productId, err := strconv.Atoi(c.Param("product_id"))
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, ErrorResponse{Message: err.Error()})
		return
	}
	reviews, err := e.services.Reviews.GetReviews(productId)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, ErrorResponse{Message: err.Error()})
		return
	}
	c.JSON(http.StatusOK, map[string]interface{}{"reviews": reviews})
}

func (e *Endpoint) DeleteReview(c *gin.Context) {
	reviewId, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, ErrorResponse{Message: err.Error()})
		return
	}
	err = e.services.Reviews.DeleteReview(reviewId)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, ErrorResponse{Message: err.Error()})
		return
	}
	c.JSON(http.StatusOK, map[string]interface{}{"message": "Review deleted successfully"})
}

type UpdateReviewInput struct {
	Rating  int    `json:"rating" binding:"required" validate:"required,min=1,max=5"`
	Comment string `json:"comment" binding:"required"`
}

func (e *Endpoint) UpdateReview(c *gin.Context) {
	reviewId, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, ErrorResponse{Message: err.Error()})
		return
	}
	var input UpdateReviewInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, ErrorResponse{Message: err.Error()})
		return
	}
	userId, err := e.GetUserId(c)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, ErrorResponse{Message: err.Error()})
		return
	}
	review := model.Review{
		ID:      reviewId,
		Rating:  input.Rating,
		Comment: input.Comment,
		UserID:  userId,
	}
	err = e.services.Reviews.UpdateReview(reviewId, review)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, ErrorResponse{Message: err.Error()})
		return
	}
	c.JSON(http.StatusOK, map[string]interface{}{"message": "Review updated successfully"})
}

func (e *Endpoint) GetProductRating(c *gin.Context) {
	productId, err := strconv.Atoi(c.Param("product_id"))
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, ErrorResponse{Message: err.Error()})
		return
	}
	rating, err := e.services.Reviews.GetProductRating(productId)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, ErrorResponse{Message: err.Error()})
		return
	}
	c.JSON(http.StatusOK, map[string]interface{}{"rating": rating})
}
