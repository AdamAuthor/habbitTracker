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
	CreateGoogle(ctx context.Context, user *models.User) error
	Create(ctx context.Context, user *models.User) error
	ConfirmRegistration(ctx context.Context, token string) error
	FindByEmail(ctx context.Context, email string) (*models.User, error)
	All(ctx context.Context) ([]*models.User, error)
	ByID(ctx context.Context, id int) (*models.User, error)
	Delete(ctx context.Context, id int) error
	Update(ctx context.Context, user *models.User) error
	FindByConfirmToken(ctx context.Context, token string) (*models.User, error)
}
