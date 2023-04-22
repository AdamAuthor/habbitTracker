package newdb

import (
	"awesomeProject/server/db"
	"awesomeProject/server/models"
	"context"
	"database/sql"
	"fmt"
	"github.com/jmoiron/sqlx"
)

func (D *DB) User() db.UserRepository {
	if D.user == nil {
		D.user = NewUserRepository(D.conn)
	}

	return D.user
}

type UserRepository struct {
	conn *sqlx.DB
}

func (u *UserRepository) DeletePasswordResetToken(ctx context.Context, email string) error {
	query := "UPDATE users SET password_reset_token = NULL, password_reset_token_created_at = NULL WHERE email = $1"
	_, err := u.conn.ExecContext(ctx, query, email)
	if err != nil {
		return err
	}
	return nil
}

func (u *UserRepository) UpdatePassword(ctx context.Context, email string, password string) error {
	query := `UPDATE users SET password = $1 WHERE email = $2`
	_, err := u.conn.ExecContext(ctx, query, password, email)
	if err != nil {
		return err
	}

	return nil
}

func (u *UserRepository) FindByPasswordResetToken(ctx context.Context, token string) (*models.User, error) {
	user := new(models.User)
	err := u.conn.GetContext(ctx, user, "SELECT * FROM users WHERE password_reset_token = $1", token)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	return user, nil
}

func (u *UserRepository) SetPasswordResetToken(ctx context.Context, email, token string) error {
	_, err := u.conn.ExecContext(ctx, "UPDATE users SET password_reset_token = $1, password_reset_token_created_at = NOW() WHERE email = $2", token, email)
	return err
}

func (u UserRepository) FindByEmail(ctx context.Context, email string) (*models.User, error) {
	con := new(models.User)
	if err := u.conn.Get(con, "SELECT * FROM users WHERE email = $1", email); err != nil {
		return nil, err
	}
	return con, nil
}

func (u *UserRepository) FindByConfirmToken(ctx context.Context, token string) (*models.User, error) {
	user := &models.User{}
	if err := u.conn.GetContext(ctx, user, "SELECT * FROM users WHERE confirm_token = $1 LIMIT 1", token); err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("User not found with the provided confirmation token")
		}
		return nil, err
	}
	return user, nil
}

func (u UserRepository) ConfirmRegistration(ctx context.Context, token string) error {
	var userID int
	err := u.conn.QueryRow("SELECT id FROM users WHERE confirm_token = $1", token).Scan(&userID)

	if err != nil {
		if err == sql.ErrNoRows {
			return fmt.Errorf("user not found")
		}
		return err
	}

	_, err = u.conn.ExecContext(ctx, "UPDATE users SET confirmed = true, updated_at = NOW() WHERE id = $1", userID)
	if err != nil {
		return err
	}

	return nil
}

func (u UserRepository) CreateGoogle(ctx context.Context, user *models.User) error {
	_, err := u.conn.Exec("INSERT INTO users(email, confirmed, nickname) VALUES ($1, $2, $3)", user.Email, user.Confirmed, user.Nickname)

	if err != nil {
		return err
	}
	return nil
}

func (u UserRepository) Create(ctx context.Context, user *models.User) error {
	_, err := u.conn.Exec("INSERT INTO users(email, password, confirm_token, nickname) VALUES ($1, $2, $3, $4)", user.Email, user.Password, user.ConfirmToken, user.Nickname)

	if err != nil {
		return err
	}
	return nil
}

func (u UserRepository) All(ctx context.Context) ([]*models.User, error) {
	//TODO implement me
	panic("implement me")
}

func (u UserRepository) ByID(ctx context.Context, id int) (*models.User, error) {
	//TODO implement me
	panic("implement me")
}

func (u UserRepository) Update(ctx context.Context, anime *models.User) error {
	//TODO implement me
	panic("implement me")
}

func (u UserRepository) Delete(ctx context.Context, id int) error {
	//TODO implement me
	panic("implement me")
}

func NewUserRepository(conn *sqlx.DB) *UserRepository {
	return &UserRepository{conn: conn}
}
