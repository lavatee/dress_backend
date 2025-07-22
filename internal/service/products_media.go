package service

import "github.com/lavatee/dresscode_backend/internal/repository"

type ProductsMediaService struct {
	repo *repository.Repository
}

func NewProductsMediaService(repo *repository.Repository) *ProductsMediaService {
	return &ProductsMediaService{
		repo: repo,
	}
}
