package models

type Habit struct {
	ID          int    `json:"id" db:"id"`
	Name        string `json:"name" db:"name"`
	Description string `json:"description" db:"description"`
	Target      int    `json:"target" db:"target"`
	Frequency   int    `json:"frequency" db:"frequency"`
	StartDate   string `json:"startDate" db:"startDate"`
	EndDate     string `json:"endDate" db:"endDate"`
}
