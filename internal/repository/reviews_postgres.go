package repository

import (
	"fmt"

	"github.com/jmoiron/sqlx"
	"github.com/lavatee/dresscode_backend/internal/model"
)

const reviewsTable = "reviews"

type ReviewsPostgres struct {
	db *sqlx.DB
}

func NewReviewsPostgres(db *sqlx.DB) *ReviewsPostgres {
	return &ReviewsPostgres{db: db}
}

func (r *ReviewsPostgres) CreateReview(review model.Review) (int, error) {
	query := fmt.Sprintf("INSERT INTO %s (product_id, user_id, rating, comment) VALUES ($1, $2, $3, $4) RETURNING id", reviewsTable)
	var id int
	err := r.db.QueryRow(query, review.ProductID, review.UserID, review.Rating, review.Comment).Scan(&id)
	if err != nil {
		return 0, err
	}
	return id, nil
}

func (r *ReviewsPostgres) GetReviews(productId int) ([]model.Review, error) {
	query := fmt.Sprintf("SELECT r.*, u.name AS user_name FROM %s r JOIN users u ON r.user_id = u.id WHERE r.product_id = $1 ORDER BY r.created_at DESC", reviewsTable)
	var reviews []model.Review
	err := r.db.Select(&reviews, query, productId)
	if err != nil {
		return nil, err
	}
	return reviews, nil
}

func (r *ReviewsPostgres) DeleteReview(reviewId int) error {
	query := fmt.Sprintf("DELETE FROM %s WHERE id = $1", reviewsTable)
	_, err := r.db.Exec(query, reviewId)
	if err != nil {
		return err
	}
	return nil
}

func (r *ReviewsPostgres) UpdateReview(reviewId int, review model.Review) error {
	query := fmt.Sprintf("UPDATE %s SET rating = $1, comment = $2 WHERE id = $3 AND user_id = $4", reviewsTable)
	_, err := r.db.Exec(query, review.Rating, review.Comment, reviewId, review.UserID)
	if err != nil {
		return err
	}
	return nil
}

func (r *ReviewsPostgres) GetProductRating(productId int) (*float64, error) {
	query := fmt.Sprintf("SELECT AVG(rating) FROM %s WHERE product_id = $1", reviewsTable)
	var rating *float64
	err := r.db.QueryRow(query, productId).Scan(&rating)
	if err != nil {
		return nil, err
	}
	return rating, nil
}
