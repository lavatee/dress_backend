package repository

import (
	"fmt"

	"github.com/jmoiron/sqlx"
	"github.com/lavatee/dresscode_backend/internal/model"
)

type AuthPostgres struct {
	db *sqlx.DB
}

func NewAuthPostgres(db *sqlx.DB) *AuthPostgres {
	return &AuthPostgres{
		db: db,
	}
}

func (r *AuthPostgres) CreateAdmin(name string, email string, password string) error {
	query := fmt.Sprintf("INSERT INTO %s (name, email, password_hash, role) VALUES ($1, $2, $3, $4)", usersTable)
	_, err := r.db.Exec(query, name, email, password, adminRole)
	return err
}

func (r *AuthPostgres) CreateUser(user model.User) (int, error) {
	var userId int
	query := fmt.Sprintf("INSERT INTO %s (name, email, password_hash, role) VALUES ($1, $2, $3, $4) RETURNING id", usersTable)
	row := r.db.QueryRow(query, user.Name, user.Email, user.Password, customerRole)

	if err := row.Scan(&userId); err != nil {
		return 0, err
	}

	return userId, nil
}

func (r *AuthPostgres) SignIn(email, password_hash string) (int, error) {
	var userId int
	query := fmt.Sprintf("SELECT id FROM %s WHERE email = $1 AND password_hash = $2", usersTable)
	row := r.db.QueryRow(query, email, password_hash)

	if err := row.Scan(&userId); err != nil {
		return 0, err
	}

	return userId, nil
}

func (r *AuthPostgres) NewAdmin(thisAdminId int, newAdminId int) error {
	query := fmt.Sprintf("UPDATE %s SET role = $1 WHERE id = $2", usersTable)
	_, err := r.db.Exec(query, adminRole, newAdminId)
	return err
}

func (r *AuthPostgres) NewBuyer(thisAdminId int, newBuyerId int) error {
	query := fmt.Sprintf("UPDATE %s SET role = $1 WHERE id = $2", usersTable)
	_, err := r.db.Exec(query, buyerRole, newBuyerId)
	return err
}

func (r *AuthPostgres) IsAdmin(userId int) bool {
	var isAdmin bool
	query := fmt.Sprintf("SELECT 1 FROM %s WHERE id = $1 AND role = $2", usersTable)
	row := r.db.QueryRow(query, userId, adminRole)

	if err := row.Scan(&isAdmin); err != nil {
		return false
	}

	return isAdmin
}

func (r *AuthPostgres) IsBuyer(userId int) bool {
	var isBuyer bool
	query := fmt.Sprintf("SELECT 1 FROM %s WHERE id = $1 AND role = $2", usersTable)
	row := r.db.QueryRow(query, userId, buyerRole)

	if err := row.Scan(&isBuyer); err != nil {
		return false
	}

	return isBuyer
}

func (r *AuthPostgres) GetUserRole(userId int) (string, error) {
	var role string
	query := fmt.Sprintf("SELECT role FROM %s WHERE id = $1", usersTable)
	row := r.db.QueryRow(query, userId)
	if err := row.Scan(&role); err != nil {
		return "", err
	}
	return role, nil
}

func (r *AuthPostgres) GetUser(userId int) (model.User, error) {
	var user model.User
	query := fmt.Sprintf("SELECT * FROM %s WHERE id = $1", usersTable)
	row := r.db.QueryRow(query, userId)
	if err := row.Scan(&user.ID, &user.Name, &user.Email, &user.Password, &user.Role); err != nil {
		return model.User{}, err
	}
	return user, nil
}

func (r *AuthPostgres) RemoveBuyer(thisAdminId int, buyerId int) error {
	query := fmt.Sprintf("UPDATE %s SET role = $1 WHERE id = $2", usersTable)
	_, err := r.db.Exec(query, customerRole, buyerId)
	return err
}
