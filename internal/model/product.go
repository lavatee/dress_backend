package model

import "time"

type Product struct {
	ID             int            `json:"id" db:"id"`
	Name           string         `json:"name" db:"name"`
	Description    string         `json:"description" db:"description"`
	Price          int            `json:"price" db:"price"`
	CollectionID   int            `json:"collection_id" db:"collection_id"`
	Category       string         `json:"category" db:"category"`
	Color          string         `json:"color" db:"color"`
	MainPhotoURL   *string        `json:"main_photo_url" db:"main_photo_url"`
	CollectionName string         `json:"collection_name" db:"collection_name"`
	IsLiked        bool           `json:"is_liked" db:"is_liked"`
	Media          []ProductMedia `json:"media"`
	Sizes          []Size         `json:"sizes"`
	Rating         *float64       `json:"rating"`
}

type ProductMedia struct {
	ID        int    `json:"id" db:"id"`
	Type      string `json:"type" db:"type"`
	URL       string `json:"url" db:"url"`
	ProductID int    `json:"product_id" db:"product_id"`
}

type Size struct {
	ID        int    `json:"id" db:"id"`
	ProductID int    `json:"product_id" db:"product_id"`
	Name      string `json:"name" db:"name"`
	Amount    int    `json:"amount" db:"amount"`
}

type LikedProduct struct {
	ID        int `json:"id" db:"id"`
	ProductID int `json:"product_id" db:"product_id"`
	UserID    int `json:"user_id" db:"user_id"`
}

type ProductInCart struct {
	ID           int     `json:"id" db:"id"`
	ProductID    int     `json:"product_id" db:"product_id"`
	ProductName  string  `json:"product_name" db:"product_name"`
	UserID       int     `json:"user_id" db:"user_id"`
	Size         string  `json:"size" db:"size"`
	Amount       int     `json:"amount" db:"amount"`
	Exists       bool    `json:"exists" db:"existence"`
	MainPhotoURL *string `json:"main_photo_url" db:"main_photo_url"`
}

type Category struct {
	ID   int    `json:"id" db:"id"`
	Name string `json:"name" db:"name"`
}

type Review struct {
	ID        int       `json:"id" db:"id"`
	ProductID int       `json:"product_id" db:"product_id"`
	UserID    int       `json:"user_id" db:"user_id"`
	Rating    int       `json:"rating" db:"rating"`
	Comment   string    `json:"comment" db:"comment"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UserName  string    `json:"user_name" db:"user_name"`
}

type Collection struct {
	ID   int    `json:"id" db:"id"`
	Name string `json:"name" db:"name"`
}
