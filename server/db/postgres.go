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
	CheckIn() CheckInRepository
}

type HabitRepository interface {
	Create(ctx context.Context, habit *models.Habit) error
	FindAll(ctx context.Context) ([]*models.Habit, error)
	FindByID(ctx context.Context, id int) (*models.Habit, error)
	Update(ctx context.Context, habit *models.Habit) error
	Delete(ctx context.Context, id int) error
}

type CheckInRepository interface {
	AddCheckIn(ctx context.Context, checkIn *models.CheckIn) error
	GetCheckIns(ctx context.Context, checkInID int, userID int) ([]*models.CheckIn, error)
}

type UserRepository interface {
	SetPasswordResetToken(ctx context.Context, email, token string) error
	CreateGoogle(ctx context.Context, user *models.User) error
	Create(ctx context.Context, user *models.User) error
	ConfirmRegistration(ctx context.Context, token string) error
	FindByEmail(ctx context.Context, email string) (*models.User, error)
	Delete(ctx context.Context, id int) error
	UpdateConfirmed(ctx context.Context, user *models.User) error
	FindByConfirmToken(ctx context.Context, token string) (*models.User, error)
	FindByPasswordResetToken(ctx context.Context, token string) (*models.User, error)
	UpdatePassword(ctx context.Context, email string, password string) error
	DeletePasswordResetToken(ctx context.Context, email string) error
}
