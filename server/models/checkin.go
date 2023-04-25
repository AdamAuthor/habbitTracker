package models

import "time"

type CheckIn struct {
	ID      int       `json:"id" db:"id"`
	HabitID int       `json:"habitId" db:"habit_id"`
	Date    time.Time `json:"date" db:"date"`
	Value   int       `json:"value" db:"value"`
	UserID  int       `json:"userId" db:"userid"`
}
