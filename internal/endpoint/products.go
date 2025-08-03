package endpoint

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/lavatee/dresscode_backend/internal/model"
)

type CreateProductInput struct {
	Name         string       `json:"name" binding:"required"`
	Description  string       `json:"description"`
	Price        int          `json:"price" binding:"required"`
	CollectionID int          `json:"collection_id"`
	Category     string       `json:"category" binding:"required"`
	Color        string       `json:"color" binding:"required"`
	Sizes        []model.Size `json:"sizes"`
}

func (e *Endpoint) PostProduct(c *gin.Context) {
	var input CreateProductInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, ErrorResponse{Message: err.Error()})
		return
	}
	userId, err := e.GetUserId(c)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, ErrorResponse{Message: err.Error()})
		return
	}
	productId, err := e.services.Products.CreateProduct(userId, model.Product{
		Name:         input.Name,
		Description:  input.Description,
		Price:        input.Price,
		CollectionID: input.CollectionID,
		Category:     input.Category,
		Color:        input.Color,
		Sizes:        input.Sizes,
	})
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, ErrorResponse{Message: err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"product_id": productId})
}

func (e *Endpoint) GetProduct(c *gin.Context) {
	productId, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, ErrorResponse{Message: err.Error()})
		return
	}
	userId, err := e.GetUserId(c)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, ErrorResponse{Message: err.Error()})
		return
	}
	product, err := e.services.Products.GetProduct(productId, userId)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, ErrorResponse{Message: err.Error()})
		return
	}
	c.JSON(http.StatusOK, map[string]interface{}{
		"product": product,
	})
}

func (e *Endpoint) GetProducts(c *gin.Context) {
	collectionId, err := strconv.Atoi(c.Query("collection_id"))
	if err != nil {
		collectionId = 0
	}
	category := c.Query("category")
	colorsString := c.Query("colors")
	colors := []string{}
	if colorsString != "" {
		colors = strings.Split(colorsString, ",")
	}
	sizesString := c.Query("sizes")
	sizes := []string{}
	if sizesString != "" {
		sizes = strings.Split(sizesString, ",")
	}
	minPrice, err := strconv.Atoi(c.Query("min_price"))
	if err != nil {
		minPrice = 0
	}
	maxPrice, err := strconv.Atoi(c.Query("max_price"))
	if err != nil {
		maxPrice = 0
	}
	page, err := strconv.Atoi(c.Query("page"))
	if err != nil {
		page = 1
	}

	// Пытаемся получить userId, но не возвращаем ошибку если токен отсутствует
	userId, _ := e.GetUserId(c)

	products, err := e.services.Products.GetProducts(collectionId, category, colors, sizes, minPrice, maxPrice, page, userId)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, ErrorResponse{Message: err.Error()})
		return
	}
	c.JSON(http.StatusOK, map[string]interface{}{
		"products": products,
	})
}

func (e *Endpoint) DeleteProduct(c *gin.Context) {
	productId, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, ErrorResponse{Message: err.Error()})
		return
	}
	userId, err := e.GetUserId(c)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, ErrorResponse{Message: err.Error()})
		return
	}
	err = e.services.Products.DeleteProduct(userId, productId)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, ErrorResponse{Message: err.Error()})
		return
	}
	c.JSON(http.StatusOK, map[string]interface{}{
		"message": "Product deleted successfully",
	})
}

type UpdateProductSizesInput struct {
	RemovedSizes []int        `json:"removed_sizes"`
	AddedSizes   []model.Size `json:"added_sizes"`
}

func (e *Endpoint) UpdateProductSizes(c *gin.Context) {
	var input UpdateProductSizesInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, ErrorResponse{Message: err.Error()})
		return
	}
	productId, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, ErrorResponse{Message: err.Error()})
		return
	}
	userId, err := e.GetUserId(c)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, ErrorResponse{Message: err.Error()})
		return
	}
	if err = e.services.Products.UpdateProductSizes(userId, productId, input.RemovedSizes, input.AddedSizes); err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, ErrorResponse{Message: err.Error()})
		return
	}
	c.JSON(http.StatusOK, map[string]interface{}{
		"message": "Product sizes updated successfully",
	})
}

func (e *Endpoint) GetProductSizes(c *gin.Context) {
	productId, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, ErrorResponse{Message: err.Error()})
		return
	}
	sizes, err := e.services.Products.GetProductSizes(productId)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, ErrorResponse{Message: err.Error()})
		return
	}
	c.JSON(http.StatusOK, map[string]interface{}{
		"sizes": sizes,
	})
}

func (e *Endpoint) CreateCollection(c *gin.Context) {
	var input model.Collection
	if err := c.ShouldBindJSON(&input); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, ErrorResponse{Message: err.Error()})
		return
	}
	userId, err := e.GetUserId(c)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, ErrorResponse{Message: err.Error()})
		return
	}
	collectionId, err := e.services.Products.CreateCollection(userId, input)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, ErrorResponse{Message: err.Error()})
		return
	}
	c.JSON(http.StatusOK, map[string]interface{}{
		"collection_id": collectionId,
	})
}

func (e *Endpoint) GetCollections(c *gin.Context) {
	collections, err := e.services.Products.GetCollections()
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, ErrorResponse{Message: err.Error()})
		return
	}
	c.JSON(http.StatusOK, map[string]interface{}{
		"collections": collections,
	})
}

type AddProductToCartInput struct {
	Size   string `json:"size" binding:"required"`
	Amount int    `json:"amount" binding:"required"`
}

func (e *Endpoint) AddProductToCart(c *gin.Context) {
	var input AddProductToCartInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, ErrorResponse{Message: err.Error()})
		return
	}
	productId, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, ErrorResponse{Message: err.Error()})
		return
	}
	userId, err := e.GetUserId(c)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, ErrorResponse{Message: err.Error()})
		return
	}
	err = e.services.Products.AddProductToCart(userId, productId, input.Size, input.Amount)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, ErrorResponse{Message: err.Error()})
		return
	}
	c.JSON(http.StatusOK, map[string]interface{}{
		"message": "Product added to cart successfully",
	})
}

func (e *Endpoint) RemoveProductFromCart(c *gin.Context) {
	productId, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, ErrorResponse{Message: err.Error()})
		return
	}
	userId, err := e.GetUserId(c)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, ErrorResponse{Message: err.Error()})
		return
	}
	err = e.services.Products.RemoveProductFromCart(userId, productId)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, ErrorResponse{Message: err.Error()})
		return
	}
	c.JSON(http.StatusOK, map[string]interface{}{
		"message": "Product removed from cart successfully",
	})
}

func (e *Endpoint) GetProductsInCart(c *gin.Context) {
	userId, err := e.GetUserId(c)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, ErrorResponse{Message: err.Error()})
		return
	}
	productsInCart, err := e.services.Products.GetProductsInCart(userId)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, ErrorResponse{Message: err.Error()})
		return
	}
	c.JSON(http.StatusOK, map[string]interface{}{
		"products_in_cart": productsInCart,
	})
}

func (e *Endpoint) AddProductToLiked(c *gin.Context) {
	productId, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, ErrorResponse{Message: err.Error()})
		return
	}
	userId, err := e.GetUserId(c)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, ErrorResponse{Message: err.Error()})
		return
	}
	err = e.services.Products.AddProductToLiked(userId, productId)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, ErrorResponse{Message: err.Error()})
		return
	}
	c.JSON(http.StatusOK, map[string]interface{}{
		"message": "Product added to liked successfully",
	})
}

func (e *Endpoint) RemoveProductFromLiked(c *gin.Context) {
	productId, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, ErrorResponse{Message: err.Error()})
		return
	}
	userId, err := e.GetUserId(c)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, ErrorResponse{Message: err.Error()})
		return
	}
	err = e.services.Products.RemoveProductFromLiked(userId, productId)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, ErrorResponse{Message: err.Error()})
		return
	}
	c.JSON(http.StatusOK, map[string]interface{}{
		"message": "Product removed from liked successfully",
	})
}

func (e *Endpoint) GetLikedProducts(c *gin.Context) {
	userId, err := e.GetUserId(c)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, ErrorResponse{Message: err.Error()})
		return
	}
	likedProducts, err := e.services.Products.GetLikedProducts(userId)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, ErrorResponse{Message: err.Error()})
		return
	}
	c.JSON(http.StatusOK, map[string]interface{}{
		"liked_products": likedProducts,
	})
}

type UpdatedSize struct {
	SizeID int `json:"size_id" binding:"required"`
	Amount int `json:"amount" binding:"required"`
}

type ChangeProductSizesAmountInput struct {
	Sizes []UpdatedSize `json:"sizes"`
}

func (e *Endpoint) ChangeProductSizesAmount(c *gin.Context) {
	var input ChangeProductSizesAmountInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, ErrorResponse{Message: err.Error()})
		return
	}
	userId, err := e.GetUserId(c)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, ErrorResponse{Message: err.Error()})
		return
	}
	sizesMap := make(map[int]int)
	for _, size := range input.Sizes {
		sizesMap[size.SizeID] = size.Amount
	}
	if err = e.services.Products.ChangeProductSizesAmount(userId, sizesMap); err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, ErrorResponse{Message: err.Error()})
		return
	}
	c.JSON(http.StatusOK, map[string]interface{}{
		"message": "Product sizes amount changed successfully",
	})
}

func (e *Endpoint) SearchProducts(c *gin.Context) {
	query := c.Query("query")
	userId, err := e.GetUserId(c)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, ErrorResponse{Message: err.Error()})
		return
	}
	products, err := e.services.Products.SearchProducts(userId, query)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, ErrorResponse{Message: err.Error()})
		return
	}
	c.JSON(http.StatusOK, map[string]interface{}{
		"products": products,
	})
}
