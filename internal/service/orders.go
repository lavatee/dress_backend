package service

import "github.com/lavatee/dresscode_backend/internal/repository"

type OrdersService struct {
	repo *repository.Repository
}

func NewOrdersService(repo *repository.Repository) *OrdersService {
	return &OrdersService{
		repo: repo,
	}
}
