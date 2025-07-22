package model

type Order struct {
	ID              int              `json:"id" db:"id"`
	Type            string           `json:"type" db:"type"`
	Status          string           `json:"status" db:"status"`
	ShopPoint       string           `json:"shop_point" db:"shop_point"`
	UserID          int              `json:"user_id" db:"user_id"`
	PaymentID       string           `json:"payment_id" db:"payment_id"`
	OrderPrice      int              `json:"order_price" db:"order_price"`
	DeliveryID      string           `json:"delivery_id" db:"delivery_id"`
	PickupID        string           `json:"pickup_id" db:"pickup_id"`
	DeliveryPrice   int              `json:"delivery_price" db:"delivery_price"`
	DeliveryAddress string           `json:"delivery_address" db:"delivery_address"`
	DeliveryIndex   int              `json:"delivery_index" db:"delivery_index"`
	UserEmail       string           `json:"user_email" db:"user_email"`
	OrderedProducts []OrderedProduct `json:"ordered_products" db:"ordered_products"`
}

type OrderedProduct struct {
	ID          int    `json:"id" db:"id"`
	OrderID     int    `json:"order_id" db:"order_id"`
	ProductID   int    `json:"product_id" db:"product_id"`
	Size        string `json:"size" db:"size"`
	Amount      int    `json:"amount" db:"amount"`
	Price       int    `json:"price" db:"price"`
	ProductName string `json:"product_name" db:"product_name"`
}
