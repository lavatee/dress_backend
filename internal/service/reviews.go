package service

import (
	"github.com/lavatee/dresscode_backend/internal/model"
	"github.com/lavatee/dresscode_backend/internal/repository"
)

type ReviewsService struct {
	repo *repository.Repository
}

func NewReviewsService(repo *repository.Repository) *ReviewsService {
	return &ReviewsService{repo: repo}
}

func (s *ReviewsService) CreateReview(review model.Review) (int, error) {
	return s.repo.Reviews.CreateReview(review)
}

func (s *ReviewsService) GetReviews(productId int) ([]model.Review, error) {
	return s.repo.Reviews.GetReviews(productId)
}

func (s *ReviewsService) DeleteReview(reviewId int) error {
	return s.repo.Reviews.DeleteReview(reviewId)
}

func (s *ReviewsService) UpdateReview(reviewId int, review model.Review) error {
	return s.repo.Reviews.UpdateReview(reviewId, review)
}

func (s *ReviewsService) GetProductRating(productId int) (*float64, error) {
	return s.repo.Reviews.GetProductRating(productId)
}
