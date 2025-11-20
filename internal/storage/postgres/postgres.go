package postgres

import (
	"avito-test-assignment-backend/internal/models"
	"context"
	"database/sql"
	"fmt"
	"strings"

	_ "github.com/lib/pq"
)

type Storage struct {
	db *sql.DB
}

func New(storagePath string) (*Storage, error) {
	const op = "storage.postgres.New"

    db, err := sql.Open("postgres", storagePath)
    if err != nil {
        return nil, fmt.Errorf("%s: %w", op, err)
    }

    return &Storage{db: db}, nil
}

func (s *Storage) BeginTx(ctx context.Context) (*sql.Tx, error) {
	return s.db.BeginTx(ctx, nil)
}

func (s *Storage) InsertTeamTx(ctx context.Context, tx *sql.Tx, teamName string) (int64, error) {
	const op = "storage.postgres.InsertTeamTx"

	res, err := tx.ExecContext(ctx, `INSERT INTO teams (team_name) VALUES ($1) ON CONFLICT DO NOTHING`, teamName)
	if err != nil {
		return 0, fmt.Errorf("%s: %w", op, err)
	}
	n, err := res.RowsAffected()
	if err != nil {
		return 0, fmt.Errorf("%s rows affected: %w", op, err)
	}
	return n, nil
}

func (s *Storage) UpsertUsersTx(ctx context.Context, tx *sql.Tx, teamName string, users []models.User) error {
	const op = "storage.postgres.UpsertUsersTx"

	if len(users) == 0 {
		return nil
	}

	placeholders := make([]string, 0, len(users))
	args := make([]interface{}, 0, len(users)*4)
	idx := 1
	for i := range users {
		placeholders = append(placeholders, fmt.Sprintf("($%d,$%d,$%d,$%d)", idx, idx+1, idx+2, idx+3))
		args = append(args, users[i].UserID, users[i].Username, teamName, users[i].IsActive)
		idx += 4
	}
	query := fmt.Sprintf(`
		INSERT INTO users (user_id, username, team_name, is_active)
		VALUES %s
		ON CONFLICT (user_id)
		DO NOTHING`, 
		strings.Join(placeholders, ","))

	if _, err := tx.ExecContext(ctx, query, args...); err != nil {
		return fmt.Errorf("%s exec: %w", op, err)
	}
	return nil
}

