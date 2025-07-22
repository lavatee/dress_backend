package endpoint

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/lavatee/dresscode_backend/internal/service"
)

type Endpoint struct {
	services *service.Service
}

type ErrorResponse struct {
	Message string `json:"message"`
}

func NewEndpoint(services *service.Service) *Endpoint {
	return &Endpoint{
		services: services,
	}
}

func (e *Endpoint) InitRoutes() *gin.Engine {
	router := gin.New()
	router.Use(func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, GET, PUT, PATCH, DELETE, OPTIONS")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "X-Auth-Token, Content-Type, Origin, Authorization")
		c.Writer.Header().Set("Access-Control-Allow-Creditionals", "true")
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(http.StatusOK)
			return
		}
	})
	auth := router.Group("/auth")
	{
		auth.POST("/sign-up", e.SignUp)
		auth.POST("/sign-in", e.SignIn)
		auth.POST("/refresh", e.Refresh)
	}
	api := router.Group("/api", e.Middleware)
	{
		api.POST("/new-admin", e.NewAdmin)
		api.POST("/new-buyer", e.NewBuyer)
		api.GET("/my-id", e.GetUserIdByToken)
		api.GET("/my-role", e.GetUserRole)
		api.GET("/my-user", e.GetUser)
		api.DELETE("/remove-buyer", e.RemoveBuyer)
	}
	products := api.Group("/products")
	{
		products.POST("/", e.PostProduct)
		products.GET("/:id", e.GetProduct)
		products.GET("/", e.GetProducts)
		products.DELETE("/:id", e.DeleteProduct)
		products.PUT("/:id/sizes", e.UpdateProductSizes)
		products.GET("/:id/sizes", e.GetProductSizes)
		products.POST("/:id/cart", e.AddProductToCart)
		products.DELETE("/:id/cart", e.RemoveProductFromCart)
		products.GET("/cart", e.GetProductsInCart)
		products.POST("/:id/liked", e.AddProductToLiked)
		products.DELETE("/:id/liked", e.RemoveProductFromLiked)
		products.GET("/liked", e.GetLikedProducts)
		products.PUT("/:id/sizes/amount", e.ChangeProductSizesAmount)
		products.GET("/search", e.SearchProducts)
		products.GET("/categories", e.GetCategories)
		products.POST("/categories", e.CreateCategory)
	}
	reviews := api.Group("/reviews")
	{
		reviews.POST("/", e.CreateReview)
		reviews.GET("/:product_id", e.GetReviews)
		reviews.DELETE("/:id", e.DeleteReview)
		reviews.PUT("/:id", e.UpdateReview)
		reviews.GET("/:product_id/rating", e.GetProductRating)
	}
	return router
}
