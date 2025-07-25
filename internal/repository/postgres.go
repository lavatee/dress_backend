package repository

import (
	"fmt"

	"github.com/jmoiron/sqlx"
)

const (
	usersTable           = "users"
	productsTable        = "products"
	ordersTable          = "orders"
	productsMediaTable   = "media"
	orderedProductsTable = "ordered_products"
	likedProductsTable   = "liked_products"
	productsInCartTable  = "products_in_cart"
	sizesTable           = "sizes"
	collectionsTable     = "collections"
)

type PostgresConfig struct {
	Host     string
	Port     string
	Username string
	Password string
	DBName   string
	SSLMode  string
}

func NewPostgresDB(config PostgresConfig) (*sqlx.DB, error) {
	db, err := sqlx.Open("postgres", fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s", config.Host, config.Port, config.Username, config.Password, config.DBName, config.SSLMode))
	if err != nil {
		return nil, err
	}
	if err := db.Ping(); err != nil {
		return nil, err
	}
	return db, nil
}
