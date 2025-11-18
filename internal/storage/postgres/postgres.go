package postgres

import (
	"avito-test-assignment-backend/api"
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

func (s *Storage) AddTeam(t api.Team) error {
    const op = "storage.postgres.AddTeam"

    tx, err := s.db.Begin()
	if err != nil {
		return fmt.Errorf("%s: begin tx: %w", op, err)
	}

    defer func() {
		if p := recover(); p != nil {
			_ = tx.Rollback()
			panic(p)
		} else if err != nil {
			_ = tx.Rollback()
		} else {
			if cerr := tx.Commit(); cerr != nil {
				err = fmt.Errorf("%s: commit tx: %w", op, cerr)
			}   
		}
	}()

    if _, err = tx.Exec(
		`INSERT INTO teams (team_name) VALUES ($1) ON CONFLICT DO NOTHING`,
		t.TeamName,
	); err != nil {
		return fmt.Errorf("%s: insert team: %w", op, err)
	}

    if len(t.Members) == 0 {
		return nil
	}

    placeholders := make([]string, 0, len(t.Members))
	args := make([]interface{}, 0, len(t.Members)*4)
	idx := 1
	for i := range t.Members {
		// если id пустой — хз пока что делаем
		placeholders = append(placeholders, fmt.Sprintf("($%d,$%d,$%d,$%d)", idx, idx+1, idx+2, idx+3))
		args = append(args,
			t.Members[i].UserID,
			t.Members[i].Username,
			t.TeamName,
			t.Members[i].IsActive,
		)
		idx += 4
	}

	query := fmt.Sprintf(`
		INSERT INTO users (user_id, username, team_name, is_active)
		VALUES %s
		ON CONFLICT (user_id)
		DO NOTHING
	`, strings.Join(placeholders, ","))

	if _, err = tx.Exec(query, args...); err != nil {
		return fmt.Errorf("%s: bulk insert/update users: %w", op, err)
	}

	return nil
}

