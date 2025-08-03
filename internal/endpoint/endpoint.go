package endpoint

import (
	"time"

	"github.com/gin-contrib/cors"
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

	config := cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "HEAD", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Length", "Content-Type", "Authorization"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}

	router.Use(cors.New(config))
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
		// Сначала идут статические маршруты
		products.GET("/collections", e.GetCollections)
		products.POST("/collections", e.CreateCollection)
		products.GET("/cart", e.GetProductsInCart)
		products.GET("/liked", e.GetLikedProducts)
		products.GET("/search", e.SearchProducts)

		// Затем маршруты с параметрами
		products.POST("/:id/cart", e.AddProductToCart)
		products.DELETE("/:id/cart", e.RemoveProductFromCart)
		products.POST("/:id/liked", e.AddProductToLiked)
		products.DELETE("/:id/liked", e.RemoveProductFromLiked)
		products.PUT("/:id/sizes", e.UpdateProductSizes)
		products.GET("/:id/sizes", e.GetProductSizes)
		products.PUT("/:id/sizes/amount", e.ChangeProductSizesAmount)
		products.GET("/:id", e.GetProduct)
		products.DELETE("/:id", e.DeleteProduct)

		// В конце общие маршруты
		products.POST("/", e.PostProduct)
		products.GET("/", e.GetProducts)
	}
	productsMedia := api.Group("/products-media")
	{
		productsMedia.POST("/", e.UploadOneProductMedia)
		productsMedia.GET("/:product_id", e.GetProductMedia)
		productsMedia.DELETE("/:media_id", e.DeleteOneProductMedia)
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
