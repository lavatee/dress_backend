package service

import (
	"context"
	"mime/multipart"

	"github.com/dgrijalva/jwt-go"
	"github.com/lavatee/dresscode_backend/internal/model"
	"github.com/lavatee/dresscode_backend/internal/repository"
	"github.com/minio/minio-go/v7"
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
	GetProducts(collectionId int, category string, colors []string, sizes []string, minPrice int, maxPrice int, page int, userId int) ([]model.Product, error)
	DeleteProduct(userId int, productId int) error
	UpdateProductSizes(userId int, productId int, removedSizes []int, addedSizes []model.Size) error
	GetProductSizes(productId int) ([]model.Size, error)
	CreateCollection(userId int, collection model.Collection) (int, error)
	GetCollections() ([]model.Collection, error)
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
	UploadOneProductMedia(ctx context.Context, userId int, productID int, media model.ProductMedia, fileName string, isProductMain bool, file multipart.File) (int, string, error)
	DeleteOneProductMedia(userId int, mediaId int) error
	GetProductMedia(productID int) ([]model.ProductMedia, error)
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

func NewService(repo *repository.Repository, s3 *minio.Client, bucket string) *Service {
	return &Service{
		Auth:          NewAuthService(repo),
		Products:      NewProductsService(repo),
		Orders:        NewOrdersService(repo),
		ProductsMedia: NewProductsMediaService(repo, s3, bucket),
		Reviews:       NewReviewsService(repo),
	}
}
