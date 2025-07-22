package service

import (
	"errors"

	"github.com/lavatee/dresscode_backend/internal/model"
	"github.com/lavatee/dresscode_backend/internal/repository"
)

var (
	ErrUserNotAdmin = errors.New("user is not admin")
	ErrUserNotBuyer = errors.New("user is not buyer")
)

type ProductsService struct {
	repo *repository.Repository
}

func NewProductsService(repo *repository.Repository) *ProductsService {
	return &ProductsService{
		repo: repo,
	}
}

func (s *ProductsService) CreateProduct(userId int, product model.Product) (int, error) {
	if !s.repo.Auth.IsAdmin(userId) {
		return 0, ErrUserNotAdmin
	}
	return s.repo.Products.CreateProduct(product)
}

func (s *ProductsService) GetProduct(productId int, userId int) (model.Product, error) {
	product, err := s.repo.Products.GetProduct(productId, userId)
	if err != nil {
		return model.Product{}, err
	}
	sizes, err := s.repo.Products.GetProductSizes(productId)
	if err != nil {
		return model.Product{}, err
	}
	product.Sizes = sizes
	media, err := s.repo.ProductsMedia.GetProductMedia(productId)
	if err != nil {
		return model.Product{}, err
	}
	product.Media = media
	rating, err := s.repo.Reviews.GetProductRating(productId)
	if err != nil {
		return model.Product{}, err
	}
	product.Rating = rating
	return product, nil
}

func (s *ProductsService) GetProducts(categoryId int, collection string, color string, sizes []string, minPrice int, maxPrice int, page int, userId int) ([]model.Product, error) {
	return s.repo.Products.GetProducts(categoryId, collection, color, sizes, minPrice, maxPrice, page, userId)
}

func (s *ProductsService) DeleteProduct(userId int, productId int) error {
	if !s.repo.Auth.IsAdmin(userId) {
		return ErrUserNotAdmin
	}
	return s.repo.Products.DeleteProduct(productId)
}

func (s *ProductsService) AddProductToCart(userId int, productId int, size string, amount int) error {
	return s.repo.Products.AddProductToCart(userId, productId, size, amount)
}

func (s *ProductsService) RemoveProductFromCart(userId int, productId int) error {
	return s.repo.Products.RemoveProductFromCart(userId, productId)
}

func (s *ProductsService) GetProductsInCart(userId int) ([]model.ProductInCart, error) {
	return s.repo.Products.GetProductsInCart(userId)
}

func (s *ProductsService) AddProductToLiked(userId int, productId int) error {
	return s.repo.Products.AddProductToLiked(userId, productId)
}

func (s *ProductsService) RemoveProductFromLiked(userId int, productId int) error {
	return s.repo.Products.RemoveProductFromLiked(userId, productId)
}

func (s *ProductsService) GetLikedProducts(userId int) ([]model.Product, error) {
	return s.repo.Products.GetLikedProducts(userId)
}

func (s *ProductsService) ChangeProductSizesAmount(userId int, sizesMap map[int]int) error {
	if !s.repo.Auth.IsAdmin(userId) {
		return ErrUserNotAdmin
	}
	return s.repo.Products.ChangeProductSizesAmount(sizesMap)
}

func (s *ProductsService) UpdateProductSizes(userId int, productId int, removedSizes []int, addedSizes []model.Size) error {
	if !s.repo.Auth.IsAdmin(userId) {
		return ErrUserNotAdmin
	}
	return s.repo.Products.UpdateProductSizes(productId, removedSizes, addedSizes)
}

func (s *ProductsService) GetProductSizes(productId int) ([]model.Size, error) {
	return s.repo.Products.GetProductSizes(productId)
}

func (s *ProductsService) CreateCategory(userId int, category model.Category) (int, error) {
	if !s.repo.Auth.IsAdmin(userId) {
		return 0, ErrUserNotAdmin
	}
	return s.repo.Products.CreateCategory(category)
}

func (s *ProductsService) GetCategories() ([]model.Category, error) {
	return s.repo.Products.GetCategories()
}

func (s *ProductsService) SearchProducts(userId int, userQuery string) ([]model.Product, error) {
	return s.repo.Products.SearchProducts(userId, userQuery)
}
