package repository

import (
	"database/sql"

	"github.com/jmoiron/sqlx"
	"github.com/lavatee/dresscode_backend/internal/model"
)

type Auth interface {
	CreateAdmin(name string, email string, password string) error
	CreateUser(user model.User) (int, error)
	SignIn(email, password_hash string) (int, error)
	NewAdmin(thisAdminId int, newAdminId int) error
	NewBuyer(thisAdminId int, newBuyerId int) error
	IsAdmin(userId int) bool
	GetUserRole(userId int) (string, error)
	GetUser(userId int) (model.User, error)
	RemoveBuyer(thisAdminId int, buyerId int) error
}

type Products interface {
	CreateProduct(product model.Product) (int, error)
	GetProduct(productId int, userId int) (model.Product, error)
	GetProducts(categoryId int, collection string, color string, sizes []string, minPrice int, maxPrice int, page int, userId int) ([]model.Product, error)
	DeleteProduct(productId int) error
	AddProductSizes(tx *sql.Tx, productId int, sizes []model.Size) error
	UpdateProductSizes(productId int, removedSizes []int, addedSizes []model.Size) error
	GetProductSizes(productId int) ([]model.Size, error)
	CreateCategory(category model.Category) (int, error)
	GetCategories() ([]model.Category, error)
	AddProductToCart(userId int, productId int, size string, amount int) error
	RemoveProductFromCart(userId int, productId int) error
	GetProductsInCart(userId int) ([]model.ProductInCart, error)
	AddProductToLiked(userId int, productId int) error
	RemoveProductFromLiked(userId int, productId int) error
	GetLikedProducts(userId int) ([]model.Product, error)
	ChangeProductSizesAmount(sizesMap map[int]int) error
	SearchProducts(userId int, userQuery string) ([]model.Product, error)
}

type Orders interface {
	CreateOrder(order model.Order) (model.Order, error)
	GetUserOrders(userId int, status string, orderType string) ([]model.Order, error)
	GetOrder(orderId int) (model.Order, error)
	GetDeliveryOrders(status string) ([]model.Order, error)
	GetPickupOrders(status string) ([]model.Order, error)
	SetOrderStatus(orderId int, status string) error
}

type ProductsMedia interface {
	CreateOneProductMedia(media model.ProductMedia, isProductMain bool) error
	DeleteOneProductMedia(mediaId int) error
	AddProductMedia(productId int, media []model.ProductMedia) error
	UpdateProductMedia(productId int, removedMedia []int, addedMedia []model.ProductMedia) error
	GetProductMedia(productId int) ([]model.ProductMedia, error)
}

type Reviews interface {
	CreateReview(review model.Review) (int, error)
	GetReviews(productId int) ([]model.Review, error)
	DeleteReview(reviewId int) error
	UpdateReview(reviewId int, review model.Review) error
	GetProductRating(productId int) (*float64, error)
}

type Repository struct {
	Auth
	Products
	Orders
	ProductsMedia
	Reviews
}

func NewRepository(db *sqlx.DB) *Repository {
	return &Repository{
		Auth:          NewAuthPostgres(db),
		Products:      NewProductsPostgres(db),
		Orders:        NewOrdersPostgres(db),
		ProductsMedia: NewProductsMediaPostgres(db),
		Reviews:       NewReviewsPostgres(db),
	}
}
