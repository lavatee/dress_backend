package endpoint

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/lavatee/dresscode_backend/internal/model"
)

type UploadOneProductMediaInput struct {
	IsProductMain bool   `json:"is_product_main"`
	Type          string `json:"type"`
	ProductID     int    `json:"product_id"`
}

func (e *Endpoint) UploadOneProductMedia(c *gin.Context) {
	isProductMain := c.Query("is_product_main") == "true"
	mediaType := c.Query("type")
	if mediaType != "photo" && mediaType != "video" {
		c.AbortWithStatusJSON(http.StatusBadRequest, ErrorResponse{Message: "Invalid media type"})
		return
	}
	productId, err := strconv.Atoi(c.Query("product_id"))
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, ErrorResponse{Message: err.Error()})
		return
	}
	file, header, err := c.Request.FormFile("file")
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, ErrorResponse{Message: err.Error()})
		return
	}
	defer file.Close()
	fileName := header.Filename
	userId, err := e.GetUserId(c)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, ErrorResponse{Message: err.Error()})
		return
	}
	media := model.ProductMedia{
		Type: mediaType,
	}
	mediaId, url, err := e.services.ProductsMedia.UploadOneProductMedia(c, userId, productId, media, fileName, isProductMain, file)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, ErrorResponse{Message: err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"media_id": mediaId, "url": url})
}

func (e *Endpoint) GetProductMedia(c *gin.Context) {
	productId, err := strconv.Atoi(c.Param("product_id"))
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, ErrorResponse{Message: err.Error()})
		return
	}
	media, err := e.services.ProductsMedia.GetProductMedia(productId)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, ErrorResponse{Message: err.Error()})
		return
	}
	c.JSON(http.StatusOK, media)
}

func (e *Endpoint) DeleteOneProductMedia(c *gin.Context) {
	mediaId, err := strconv.Atoi(c.Param("media_id"))
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, ErrorResponse{Message: err.Error()})
		return
	}
	userId, err := e.GetUserId(c)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, ErrorResponse{Message: err.Error()})
		return
	}
	err = e.services.ProductsMedia.DeleteOneProductMedia(userId, mediaId)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, ErrorResponse{Message: err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Media deleted successfully"})
}
