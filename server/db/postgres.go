package db

import (
	"awesomeProject/server/models"
	"context"
)

type Database interface {
	Connect(url string) error
	Close() error
	Habit() HabitRepository
	User() UserRepository
}

type HabitRepository interface {
	Create(ctx context.Context, habit *models.Habit) error
	All(ctx context.Context) ([]*models.Habit, error)
	ByID(ctx context.Context, id int) (*models.Habit, error)
	Update(ctx context.Context, habit *models.Habit) error
	Delete(ctx context.Context, id int) error
}

type UserRepository interface {
	Create(ctx context.Context, anime *models.Habit) error
	All(ctx context.Context) ([]*models.Habit, error)
	ByID(ctx context.Context, id int) (*models.Habit, error)
	Update(ctx context.Context, anime *models.Habit) error
	Delete(ctx context.Context, id int) error
}
