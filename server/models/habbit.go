package models

import "time"

type Habit struct {
	ID          int       `json:"id" db:"id"`
	Name        string    `json:"name" db:"name"`
	Description string    `json:"description" db:"description"`
	Target      int       `json:"target" db:"target"`
	Measurement string    `json:"measurement" db:"measurement"`
	StartDate   time.Time `json:"startDate" db:"startdate"`
	EndDate     time.Time `json:"endDate" db:"enddate"`
	TotalValue  int       `json:"totalValue" db:"totalvalue"`
	TotalDays   int       `json:"totalDays" db:"totaldays"`
	LastCheckIn time.Time `json:"lastCheckIn" db:"lastcheckin"`
	CheckIns    []CheckIn `json:"checkIns"`
}
