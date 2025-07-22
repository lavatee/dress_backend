package service

import (
	"github.com/dgrijalva/jwt-go"
	"github.com/lavatee/dresscode_backend/internal/model"
	"github.com/lavatee/dresscode_backend/internal/repository"
)

type Auth interface {
	CreateAdmin(name string, email string, password string) error
	SignUp(user model.User) (int, error)
	SignIn(email, password string) (string, string, error)
	Refresh(refreshToken string) (string, string, error)
	ParseToken(token string) (jwt.MapClaims, error)
	NewAdmin(thisAdminId, newAdminId int) error
	NewBuyer(thisAdminId, newBuyerId int) error
	RemoveBuyer(thisAdminId, buyerId int) error
	GetUserRole(userId int) (string, error)
	GetUser(userId int) (model.User, error)
}

type Products interface {
	CreateProduct(userId int, product model.Product) (int, error)
	GetProduct(productId int, userId int) (model.Product, error)
	GetProducts(categoryId int, collection string, color string, sizes []string, minPrice int, maxPrice int, page int, userId int) ([]model.Product, error)
	DeleteProduct(userId int, productId int) error
	UpdateProductSizes(userId int, productId int, removedSizes []int, addedSizes []model.Size) error
	GetProductSizes(productId int) ([]model.Size, error)
	CreateCategory(userId int, category model.Category) (int, error)
	GetCategories() ([]model.Category, error)
	AddProductToCart(userId int, productId int, size string, amount int) error
	RemoveProductFromCart(userId int, productId int) error
	GetProductsInCart(userId int) ([]model.ProductInCart, error)
	AddProductToLiked(userId int, productId int) error
	RemoveProductFromLiked(userId int, productId int) error
	GetLikedProducts(userId int) ([]model.Product, error)
	ChangeProductSizesAmount(userId int, sizesMap map[int]int) error
	SearchProducts(userId int, userQuery string) ([]model.Product, error)
}

type Orders interface {
}

type ProductsMedia interface {
}

type Reviews interface {
	CreateReview(review model.Review) (int, error)
	GetReviews(productId int) ([]model.Review, error)
	DeleteReview(reviewId int) error
	UpdateReview(reviewId int, review model.Review) error
	GetProductRating(productId int) (*float64, error)
}

type Service struct {
	Auth
	Products
	Orders
	ProductsMedia
	Reviews
}

func NewService(repo *repository.Repository) *Service {
	return &Service{
		Auth:          NewAuthService(repo),
		Products:      NewProductsService(repo),
		Orders:        NewOrdersService(repo),
		ProductsMedia: NewProductsMediaService(repo),
		Reviews:       NewReviewsService(repo),
	}
}
