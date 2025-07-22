package repository

import (
	"database/sql"
	"fmt"
	"strings"

	"github.com/jmoiron/sqlx"
	"github.com/lavatee/dresscode_backend/internal/model"
	"github.com/sirupsen/logrus"
)

const (
	adminRole    = "admin"
	buyerRole    = "buyer"
	customerRole = "customer"
	pageLimit    = 21
)

type ProductsPostgres struct {
	db *sqlx.DB
}

func NewProductsPostgres(db *sqlx.DB) *ProductsPostgres {
	return &ProductsPostgres{
		db: db,
	}
}

func (r *ProductsPostgres) CreateProduct(product model.Product) (int, error) {
	var productId int
	tx, err := r.db.Begin()
	if err != nil {
		return 0, err
	}
	query := fmt.Sprintf("INSERT INTO %s (name, description, price, category_id, collection, color) VALUES ($1, $2, $3, $4, $5, $6) RETURNING id", productsTable)
	row := tx.QueryRow(query, product.Name, product.Description, product.Price, product.CategoryID, product.Collection, product.Color)

	if err := row.Scan(&productId); err != nil {
		tx.Rollback()
		return 0, err
	}
	if err := r.AddProductSizes(tx, productId, product.Sizes); err != nil {
		tx.Rollback()
		return 0, err
	}
	if err := tx.Commit(); err != nil {
		return 0, err
	}

	return productId, nil
}

func (r *ProductsPostgres) GetProduct(productId int, userId int) (model.Product, error) {
	query := fmt.Sprintf("SELECT p.*, c.name as category_name, EXISTS(SELECT 1 FROM %s lp WHERE lp.product_id = p.id AND lp.user_id = $2) as is_liked FROM %s p JOIN %s c ON p.category_id = c.id WHERE p.id = $1", likedProductsTable, productsTable, categoriesTable)
	var product model.Product
	if err := r.db.Get(&product, query, productId, userId); err != nil {
		return model.Product{}, err
	}
	return product, nil
}

func (r *ProductsPostgres) GetProducts(categoryId int, collection string, color string, sizes []string, minPrice int, maxPrice int, page int, userId int) ([]model.Product, error) {
	query := fmt.Sprintf(`
		SELECT DISTINCT 
			p.*,
			c.name as category_name,
			EXISTS(SELECT 1 FROM %s lp WHERE lp.product_id = p.id AND lp.user_id = $1) as is_liked 
		FROM %s p
		LEFT JOIN %s c ON p.category_id = c.id`,
		likedProductsTable, productsTable, categoriesTable)

	args := []interface{}{userId}
	where := ""
	argIdx := 2

	if collection != "" {
		where += fmt.Sprintf("p.collection = $%d AND ", argIdx)
		args = append(args, collection)
		argIdx++
	}

	if color != "" {
		where += fmt.Sprintf("p.color = $%d AND ", argIdx)
		args = append(args, color)
		argIdx++
	}

	if categoryId != 0 {
		where += fmt.Sprintf("p.category_id = $%d AND ", argIdx)
		args = append(args, categoryId)
		argIdx++
	}

	if minPrice > 0 {
		where += fmt.Sprintf("p.price >= $%d AND ", argIdx)
		args = append(args, minPrice)
		argIdx++
	}

	if maxPrice > 0 {
		where += fmt.Sprintf("p.price <= $%d AND ", argIdx)
		args = append(args, maxPrice)
		argIdx++
	}

	if len(sizes) > 0 {
		query += fmt.Sprintf(" LEFT JOIN %s s ON s.product_id = p.id", sizesTable)
		where += "("
		for i, size := range sizes {
			if i > 0 {
				where += " OR "
			}
			where += fmt.Sprintf("(s.name = $%d AND s.amount > 0)", argIdx)
			args = append(args, size)
			argIdx++
		}
		where += ") AND "
	}

	if len(where) > 0 {
		where = where[:len(where)-5]
		query += " WHERE " + where
	}

	offset := (page - 1) * pageLimit
	query += fmt.Sprintf(" ORDER BY p.id LIMIT %d OFFSET %d", pageLimit, offset)
	logrus.Infof("query: %s", query)
	var products []model.Product
	if err := r.db.Select(&products, query, args...); err != nil {
		return nil, fmt.Errorf("error getting products: %w", err)
	}

	return products, nil
}

func (r *ProductsPostgres) DeleteProduct(productId int) error {
	query := fmt.Sprintf("DELETE FROM %s WHERE product_id = $1", productsInCartTable)
	_, err := r.db.Exec(query, productId)
	if err != nil {
		return err
	}
	query = fmt.Sprintf("DELETE FROM %s WHERE product_id = $1", likedProductsTable)
	_, err = r.db.Exec(query, productId)
	if err != nil {
		return err
	}
	query = fmt.Sprintf("DELETE FROM %s WHERE product_id = $1", productsMediaTable)
	_, err = r.db.Exec(query, productId)
	if err != nil {
		return err
	}
	query = fmt.Sprintf("DELETE FROM %s WHERE product_id = $1", sizesTable)
	_, err = r.db.Exec(query, productId)
	if err != nil {
		return err
	}
	query = fmt.Sprintf("DELETE FROM %s WHERE id = $1", productsTable)
	_, err = r.db.Exec(query, productId)
	return err
}

func (r *ProductsPostgres) AddProductSizes(tx *sql.Tx, productId int, sizes []model.Size) error {
	query := fmt.Sprintf("INSERT INTO %s (product_id, name, amount) VALUES ", sizesTable)
	argsCounter := 0
	args := make([]interface{}, 0)
	for _, s := range sizes {
		query += fmt.Sprintf("($%d, $%d, $%d),", argsCounter+1, argsCounter+2, argsCounter+3)
		args = append(args, productId, s.Name, s.Amount)
		argsCounter += 3
	}
	query = query[:len(query)-1]
	_, err := tx.Exec(query, args...)
	return err
}

func (r *ProductsPostgres) AddProductToCart(userId int, productId int, size string, amount int) error {
	query := fmt.Sprintf("INSERT INTO %s (user_id, product_id, size, amount) VALUES ($1, $2, $3, $4)", productsInCartTable)
	_, err := r.db.Exec(query, userId, productId, size, amount)
	return err
}

func (r *ProductsPostgres) RemoveProductFromCart(userId int, productId int) error {
	query := fmt.Sprintf("DELETE FROM %s WHERE user_id = $1 AND product_id = $2", productsInCartTable)
	_, err := r.db.Exec(query, userId, productId)
	return err
}

func (r *ProductsPostgres) GetProductsInCart(userId int) ([]model.ProductInCart, error) {
	query := fmt.Sprintf("SELECT c.*, p.name as product_name, p.main_photo_url FROM %s c JOIN %s p ON c.product_id = p.id WHERE c.user_id = $1", productsInCartTable, productsTable)
	rows, err := r.db.Query(query, userId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var productsInCart []model.ProductInCart
	if err := r.db.Select(&productsInCart, query, userId); err != nil {
		return nil, err
	}
	return productsInCart, nil
}

func (r *ProductsPostgres) AddProductToLiked(userId int, productId int) error {
	query := fmt.Sprintf("INSERT INTO %s (user_id, product_id) VALUES ($1, $2)", likedProductsTable)
	_, err := r.db.Exec(query, userId, productId)
	return err
}

func (r *ProductsPostgres) RemoveProductFromLiked(userId int, productId int) error {
	query := fmt.Sprintf("DELETE FROM %s WHERE user_id = $1 AND product_id = $2", likedProductsTable)
	_, err := r.db.Exec(query, userId, productId)
	return err
}

func (r *ProductsPostgres) GetLikedProducts(userId int) ([]model.Product, error) {
	query := fmt.Sprintf("SELECT p.*, c.name as category_name, EXISTS(SELECT 1 FROM %s lp WHERE lp.product_id = p.id AND lp.user_id = $1) as is_liked FROM %s lp JOIN %s p ON lp.product_id = p.id JOIN %s c ON p.category_id = c.id WHERE lp.user_id = $1", likedProductsTable, likedProductsTable, productsTable, categoriesTable)
	var products []model.Product
	if err := r.db.Select(&products, query, userId); err != nil {
		return nil, err
	}
	return products, nil
}

func (r *ProductsPostgres) ChangeProductSizesAmount(sizesMap map[int]int) error {
	if len(sizesMap) == 0 {
		return nil
	}
	query := fmt.Sprintf("UPDATE %s SET amount = CASE", sizesTable)
	args := make([]interface{}, 0, len(sizesMap)*2)
	ids := make([]interface{}, 0, len(sizesMap))
	argIdx := 1
	for sizeId, amount := range sizesMap {
		query += fmt.Sprintf(" WHEN id = $%d THEN ($%d)::bigint", argIdx, argIdx+1)
		args = append(args, sizeId, amount)
		ids = append(ids, sizeId)
		argIdx += 2
	}
	query += " END WHERE id IN ("
	idStartIdx := argIdx
	for i := range ids {
		if i > 0 {
			query += ", "
		}
		query += fmt.Sprintf("$%d", idStartIdx+i)
	}
	query += ")"
	args = append(args, ids...)
	_, err := r.db.Exec(query, args...)
	return err
}

func (r *ProductsPostgres) UpdateProductSizes(productId int, removedSizes []int, addedSizes []model.Size) error {
	tx, err := r.db.Begin()
	if err != nil {
		return err
	}
	if len(addedSizes) > 0 {
		if err := r.AddProductSizes(tx, productId, addedSizes); err != nil {
			tx.Rollback()
			return err
		}
	}
	if len(removedSizes) == 0 {
		if err := tx.Commit(); err != nil {
			return err
		}
		return nil
	}
	query := fmt.Sprintf("DELETE FROM %s WHERE", sizesTable)
	queryToGetSizes := fmt.Sprintf("SELECT * FROM %s WHERE", sizesTable)
	argsCounter := 0
	args := make([]interface{}, 0)
	for _, sizeId := range removedSizes {
		query += fmt.Sprintf("id = $%d OR ", argsCounter+1)
		queryToGetSizes += fmt.Sprintf("id = $%d OR ", argsCounter)
		args = append(args, sizeId)
		argsCounter += 1
	}
	query = query[:len(query)-4]
	queryToGetSizes = queryToGetSizes[:len(queryToGetSizes)-4]
	var deletedSizes []model.Size
	if err := r.db.Select(&deletedSizes, queryToGetSizes, args...); err != nil {
		tx.Rollback()
		return err
	}
	_, err = tx.Exec(query, args...)
	if err != nil {
		tx.Rollback()
		return err
	}
	sizesIds := make([]string, 0)
	for _, size := range deletedSizes {
		sizesIds = append(sizesIds, fmt.Sprintf("size = '%s'", size.Name))
	}
	query = fmt.Sprintf("UPDATE %s SET existence = FALSE WHERE %s AND product_id = $1", productsInCartTable, strings.Join(sizesIds, " OR "))
	_, err = tx.Exec(query, productId)
	if err != nil {
		tx.Rollback()
		return err
	}
	if err := tx.Commit(); err != nil {
		return err
	}
	return nil
}

func (r *ProductsPostgres) GetProductSizes(productId int) ([]model.Size, error) {
	query := fmt.Sprintf("SELECT * FROM %s WHERE product_id = $1", sizesTable)
	var sizes []model.Size
	if err := r.db.Select(&sizes, query, productId); err != nil {
		return nil, err
	}
	return sizes, nil
}

func (r *ProductsPostgres) CreateCategory(category model.Category) (int, error) {
	query := fmt.Sprintf("INSERT INTO %s (name) VALUES ($1) RETURNING id", categoriesTable)
	row := r.db.QueryRow(query, category.Name)
	var categoryId int
	if err := row.Scan(&categoryId); err != nil {
		return 0, err
	}
	return categoryId, nil
}

func (r *ProductsPostgres) GetCategories() ([]model.Category, error) {
	query := fmt.Sprintf("SELECT * FROM %s", categoriesTable)
	var categories []model.Category
	if err := r.db.Select(&categories, query); err != nil {
		return nil, err
	}
	return categories, nil
}

func (r *ProductsPostgres) SearchProducts(userId int, userQuery string) ([]model.Product, error) {

	query := fmt.Sprintf(`
	SELECT DISTINCT 
		p.*,
		c.name as category_name,
		EXISTS(SELECT 1 FROM %s lp WHERE lp.product_id = p.id AND lp.user_id = $1) as is_liked 
	FROM %s p
	LEFT JOIN %s c ON p.category_id = c.id WHERE p.name ILIKE $2`,
		likedProductsTable, productsTable, categoriesTable)
	var products []model.Product
	if err := r.db.Select(&products, query, userId, "%"+userQuery+"%"); err != nil {
		return nil, err
	}
	return products, nil
}
