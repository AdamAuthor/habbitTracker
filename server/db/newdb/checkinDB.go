package newdb

import (
	"awesomeProject/server/db"
	"awesomeProject/server/models"
	"context"
	"github.com/jmoiron/sqlx"
)

func (D *DB) CheckIn() db.CheckInRepository {
	if D.checkIn == nil {
		D.checkIn = NewChekInRepository(D.conn)
	}

	return D.checkIn
}

type CheckInRepository struct {
	conn *sqlx.DB
}

func NewChekInRepository(conn *sqlx.DB) *CheckInRepository {
	return &CheckInRepository{conn: conn}
}

func (r *CheckInRepository) AddCheckIn(ctx context.Context, checkIn *models.CheckIn) error {
	query := `INSERT INTO checkins(habitid, date, value, userid) VALUES ($1, $2, $3, $4) RETURNING id`
	err := r.conn.GetContext(ctx, &checkIn.ID, query, checkIn.HabitID, checkIn.Date, checkIn.Value, checkIn.UserID)
	if err != nil {
		return err
	}

	return nil
}

func (r *CheckInRepository) GetCheckIns(ctx context.Context, habitID int, userID int) ([]*models.CheckIn, error) {
	query := `SELECT * FROM checkins WHERE habitid = $1 AND userid = $2`
	var checkIns []*models.CheckIn
	err := r.conn.SelectContext(ctx, &checkIns, query, habitID, userID)
	if err != nil {
		return nil, err
	}

	return checkIns, nil
}
