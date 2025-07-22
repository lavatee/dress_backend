package repository

import (
	"database/sql"
	"fmt"

	"github.com/jmoiron/sqlx"
	"github.com/lavatee/dresscode_backend/internal/model"
)

const (
	deliveryType              = "delivery"
	pickupType                = "pickup"
	pendingStatus             = "pending"
	createdStatus             = "created"
	deliveredToShopStatus     = "delivered_to_shop"
	issuedStatus              = "issued"
	sentToCustomerStatus      = "sent_to_customer"
	deliveredToCustomerStatus = "delivered_to_customer"
)

type OrdersPostgres struct {
	db *sqlx.DB
}

func NewOrdersPostgres(db *sqlx.DB) *OrdersPostgres {
	return &OrdersPostgres{
		db: db,
	}
}

func (r *OrdersPostgres) CreateOrder(order model.Order) (model.Order, error) {
	tx, err := r.db.Begin()
	if err != nil {
		return model.Order{}, err
	}
	query := fmt.Sprintf("INSERT INTO %s (user_id, status, type, shop_point, payment_id, order_price, delivery_address, delivery_index, delivery_price) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9) RETURNING id", ordersTable)
	row := tx.QueryRow(query, order.UserID, order.Status, order.Type, order.ShopPoint, order.PaymentID, order.OrderPrice, order.DeliveryAddress, order.DeliveryIndex, order.DeliveryPrice)
	if err := row.Scan(&order.ID); err != nil {
		tx.Rollback()
		return model.Order{}, err
	}
	if err := r.CreateOrderedProducts(tx, order.OrderedProducts); err != nil {
		tx.Rollback()
		return model.Order{}, err
	}
	if err := tx.Commit(); err != nil {
		tx.Rollback()
		return model.Order{}, err
	}
	return order, nil
}

func (r *OrdersPostgres) CreateOrderedProducts(tx *sql.Tx, orderedProducts []model.OrderedProduct) error {
	query := fmt.Sprintf("INSERT INTO %s (order_id, product_id, size, amount, price, product_name) VALUES ", orderedProductsTable)
	argsCounter := 0
	args := make([]interface{}, 0)
	sizes := make([]SizeInfo, 0)
	sizesMap := make(map[string]bool)
	for _, orderedProduct := range orderedProducts {
		query += fmt.Sprintf("($%d, $%d, $%d, $%d, $%d, $%d),", argsCounter+1, argsCounter+2, argsCounter+3, argsCounter+4, argsCounter+5, argsCounter+6)
		args = append(args, orderedProduct.OrderID, orderedProduct.ProductID, orderedProduct.Size, orderedProduct.Amount, orderedProduct.Price, orderedProduct.ProductName)
		argsCounter += 6
		if _, exists := sizesMap[fmt.Sprintf("%d_%s", orderedProduct.ProductID, orderedProduct.Size)]; exists {
			return fmt.Errorf("2 sizes with the same product_id and size_name")
		}
		sizes = append(sizes, SizeInfo{ProductID: orderedProduct.ProductID, SizeName: orderedProduct.Size, Amount: orderedProduct.Amount})
		sizesMap[fmt.Sprintf("%d_%s", orderedProduct.ProductID, orderedProduct.Size)] = true
	}
	query = query[:len(query)-1]
	_, err := tx.Exec(query, args...)
	if err != nil {
		return err
	}
	if err := r.DecreaseProductSizeAmount(tx, sizes); err != nil {
		return err
	}
	return nil
}

type SizeInfo struct {
	ProductID int
	SizeName  string
	Amount    int
}

func (r *OrdersPostgres) DecreaseProductSizeAmount(tx *sql.Tx, sizes []SizeInfo) error {
	query := fmt.Sprintf("UPDATE %s SET amount = amount - CASE", sizesTable)
	whereConditions := ""
	for _, size := range sizes {
		query += fmt.Sprintf(" WHEN product_id = %d AND name = '%s' AND amount >= %d THEN %d",
			size.ProductID, size.SizeName, size.Amount, size.Amount)
		if whereConditions != "" {
			whereConditions += " OR "
		}
		whereConditions += fmt.Sprintf("(product_id = %d AND name = '%s')",
			size.ProductID, size.SizeName)
	}
	query += " ELSE 0 END WHERE " + whereConditions

	result, err := tx.Exec(query)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected != int64(len(sizes)) {
		return fmt.Errorf("not enough size on the warehouse or size not found")
	}

	return nil
}

func (r *OrdersPostgres) GetUserOrders(userId int, status string, orderType string) ([]model.Order, error) {
	var orders []model.Order
	query := fmt.Sprintf("SELECT * FROM %s WHERE user_id = $1", ordersTable)
	args := []interface{}{userId}
	if status != "" {
		query += " AND status = $2"
		args = append(args, status)
	}
	if orderType != "" {
		query += " AND type = $3"
		args = append(args, orderType)
	}
	if err := r.db.Select(&orders, query, args...); err != nil {
		return nil, err
	}
	return orders, nil
}

func (r *OrdersPostgres) GetOrder(orderId int) (model.Order, error) {
	var order model.Order
	query := fmt.Sprintf("SELECT o.*, u.email as user_email, u.name as user_name FROM %s o LEFT JOIN %s u ON o.user_id = u.id WHERE o.id = $1", ordersTable, usersTable)
	if err := r.db.Get(&order, query, orderId); err != nil {
		return model.Order{}, err
	}
	query = fmt.Sprintf("SELECT * FROM %s WHERE order_id = $1", orderedProductsTable)
	if err := r.db.Select(&order.OrderedProducts, query, orderId); err != nil {
		return model.Order{}, err
	}
	return order, nil
}

func (r *OrdersPostgres) GetDeliveryOrders(status string) ([]model.Order, error) {
	var orders []model.Order
	query := fmt.Sprintf("SELECT o.*, u.email as user_email, u.name as user_name FROM %s o LEFT JOIN %s u ON o.user_id = u.id WHERE status = $1 AND type = $2", ordersTable, usersTable)
	if err := r.db.Select(&orders, query, status, deliveryType); err != nil {
		return nil, err
	}
	return orders, nil
}

func (r *OrdersPostgres) GetPickupOrders(status string) ([]model.Order, error) {
	var orders []model.Order
	query := fmt.Sprintf("SELECT o.*, u.email as user_email, u.name as user_name FROM %s o LEFT JOIN %s u ON o.user_id = u.id WHERE status = $1 AND type = $2", ordersTable, usersTable)
	if err := r.db.Select(&orders, query, status, pickupType); err != nil {
		return nil, err
	}
	return orders, nil
}

func (r *OrdersPostgres) SetOrderStatus(orderId int, status string) error {
	query := fmt.Sprintf("UPDATE %s SET status = $1 WHERE id = $2", ordersTable)
	_, err := r.db.Exec(query, status, orderId)
	if err != nil {
		return err
	}
	return nil
}

func (r *OrdersPostgres) SetDeliveryId(orderId int, deliveryId string) error {
	query := fmt.Sprintf("UPDATE %s SET delivery_id = $1 WHERE id = $2", ordersTable)
	_, err := r.db.Exec(query, deliveryId, orderId)
	if err != nil {
		return err
	}
	return nil
}

func (r *OrdersPostgres) SetPickupId(orderId int, pickupId string) error {
	query := fmt.Sprintf("UPDATE %s SET pickup_id = $1 WHERE id = $2", ordersTable)
	_, err := r.db.Exec(query, pickupId, orderId)
	if err != nil {
		return err
	}
	return nil
}
