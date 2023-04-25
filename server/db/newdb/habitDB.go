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

func (r *HabitRepository) Create(ctx context.Context, habit *models.Habit) error {
	query := `INSERT INTO habits(name, description, target, measurement, startdate, enddate) VALUES ($1, $2, $3, $4, $5, $6) RETURNING id`
	err := r.conn.GetContext(ctx, &habit.ID, query, habit.Name, habit.Description, habit.Target, habit.Measurement, habit.StartDate, habit.EndDate)
	if err != nil {
		return err
	}
	return nil
}

func (r *HabitRepository) FindByID(ctx context.Context, id int) (*models.Habit, error) {
	habit := new(models.Habit)
	query := `SELECT * FROM habits WHERE id = $1`
	err := r.conn.GetContext(ctx, habit, query, id)
	if err != nil {
		return nil, err
	}

	return habit, nil
}

func (r *HabitRepository) FindAll(ctx context.Context) ([]*models.Habit, error) {
	var habits []*models.Habit
	err := r.conn.SelectContext(ctx, &habits, "SELECT * FROM habits")
	return habits, err
}

func (r *HabitRepository) Update(ctx context.Context, habit *models.Habit) error {
	query := `UPDATE habits SET name=:name, description=:description, target=:target, measurement=:measurement, 
		startdate=:startdate, enddate=:enddate, totalvalue=:totalvalue, totaldays=:totaldays, lastcheckin=:lastcheckin WHERE id=:id`
	_, err := r.conn.NamedExec(query, habit)
	return err
}

func (r *HabitRepository) Delete(ctx context.Context, id int) error {
	query := `DELETE FROM habits WHERE id = $1`
	_, err := r.conn.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}

	return nil
}
