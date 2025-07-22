package repository

import (
	"fmt"

	"github.com/jmoiron/sqlx"
	"github.com/lavatee/dresscode_backend/internal/model"
)

const (
	PhotoType = "photo"
	VideoType = "video"
)

type ProductsMediaPostgres struct {
	db *sqlx.DB
}

func NewProductsMediaPostgres(db *sqlx.DB) *ProductsMediaPostgres {
	return &ProductsMediaPostgres{
		db: db,
	}
}

func (r *ProductsMediaPostgres) CreateOneProductMedia(media model.ProductMedia, isProductMain bool) error {
	tx, err := r.db.Begin()
	query := fmt.Sprintf("INSERT INTO %s (type, url, product_id) VALUES ($1, $2, $3)", productsMediaTable)
	_, err = tx.Exec(query, media.Type, media.URL, media.ProductID)
	if err != nil {
		tx.Rollback()
		return err
	}
	if isProductMain && media.Type == PhotoType {
		query = fmt.Sprintf("UPDATE %s SET main_photo_url = $1 WHERE id = $2", productsTable)
		_, err = tx.Exec(query, media.URL, media.ProductID)
		if err != nil {
			tx.Rollback()
			return err
		}
	}
	if err := tx.Commit(); err != nil {
		return err
	}
	return nil
}

func (r *ProductsMediaPostgres) DeleteOneProductMedia(mediaId int) error {
	tx, err := r.db.Begin()
	if err != nil {
		return err
	}
	var mediaUrl string
	query := fmt.Sprintf("DELETE FROM %s WHERE id = $1 RETURNING url", productsMediaTable)
	row := tx.QueryRow(query, mediaId)
	if err := row.Scan(&mediaUrl); err != nil {
		tx.Rollback()
		return err
	}
	query = fmt.Sprintf("UPDATE %s SET main_photo_url = NULL WHERE main_photo_url = $1", productsTable)
	_, err = tx.Exec(query, mediaUrl)
	if err != nil {
		tx.Rollback()
		return err
	}
	if err := tx.Commit(); err != nil {
		return err
	}
	return err
}

func (r *ProductsMediaPostgres) AddProductMedia(productId int, media []model.ProductMedia) error {
	query := fmt.Sprintf("INSERT INTO %s (type, url, product_id) VALUES ", productsMediaTable)
	argsCounter := 0
	args := make([]interface{}, 0)
	for _, m := range media {
		query += fmt.Sprintf("($%d, $%d, $%d),", argsCounter+1, argsCounter+2, argsCounter+3)
		args = append(args, m.Type, m.URL, productId)
		argsCounter += 3
	}
	query = query[:len(query)-1]
	_, err := r.db.Exec(query, args...)
	return err
}

func (r *ProductsMediaPostgres) UpdateProductMedia(productId int, removedMedia []int, addedMedia []model.ProductMedia) error {
	query := fmt.Sprintf("DELETE FROM %s WHERE", productsMediaTable)
	argsCounter := 0
	args := make([]interface{}, 0)
	for _, mediaId := range removedMedia {
		query += fmt.Sprintf("id = $%d OR ", argsCounter+1)
		args = append(args, mediaId)
		argsCounter += 1
	}
	query = query[:len(query)-4]
	_, err := r.db.Exec(query, args...)
	if err != nil {
		return err
	}
	if len(addedMedia) > 0 {
		if err := r.AddProductMedia(productId, addedMedia); err != nil {
			return err
		}
	}
	return nil
}

func (r *ProductsMediaPostgres) GetProductMedia(productId int) ([]model.ProductMedia, error) {
	query := fmt.Sprintf("SELECT * FROM %s WHERE product_id = $1", productsMediaTable)
	var media []model.ProductMedia
	if err := r.db.Select(&media, query, productId); err != nil {
		return nil, err
	}
	return media, nil
}
