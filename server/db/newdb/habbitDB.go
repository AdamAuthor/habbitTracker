package newdb

import (
	"awesomeProject/server/db"
	"awesomeProject/server/models"
	"context"
	"github.com/jmoiron/sqlx"
)

func (D *DB) Habit() db.HabitRepository {
	if D.habit == nil {
		D.habit = NewHabitRepository(D.conn)
	}

	return D.habit
}

type HabitRepository struct {
	conn *sqlx.DB
}

func NewHabitRepository(conn *sqlx.DB) *HabitRepository {
	return &HabitRepository{conn: conn}
}

func (h HabitRepository) Create(ctx context.Context, anime *models.Habit) error {
	//TODO implement me
	panic("implement me")
}

func (h HabitRepository) All(ctx context.Context) ([]*models.Habit, error) {
	//TODO implement me
	panic("implement me")
}

func (h HabitRepository) ByID(ctx context.Context, id int) (*models.Habit, error) {
	//TODO implement me
	panic("implement me")
}

func (h HabitRepository) Update(ctx context.Context, anime *models.Habit) error {
	//TODO implement me
	panic("implement me")
}

func (h HabitRepository) Delete(ctx context.Context, id int) error {
	//TODO implement me
	panic("implement me")
}
