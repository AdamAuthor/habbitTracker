package newdb

import (
	"awesomeProject/server/db"
	"awesomeProject/server/models"
	"context"
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

func (u UserRepository) Create(ctx context.Context, anime *models.Habit) error {
	//TODO implement me
	panic("implement me")
}

func (u UserRepository) All(ctx context.Context) ([]*models.Habit, error) {
	//TODO implement me
	panic("implement me")
}

func (u UserRepository) ByID(ctx context.Context, id int) (*models.Habit, error) {
	//TODO implement me
	panic("implement me")
}

func (u UserRepository) Update(ctx context.Context, anime *models.Habit) error {
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
