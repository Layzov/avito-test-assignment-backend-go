package postgres

import (
	"avito-test-assignment-backend/internal/models"
	"avito-test-assignment-backend/pkg/response"
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


// #### team/add ####

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
		placeholders = append(
			placeholders, 
			fmt.Sprintf("($%d,$%d,$%d,$%d)", 
			idx, idx+1, idx+2, idx+3))
		args = append(
			args, 
			users[i].UserID, 
			users[i].Username, 
			teamName, 
			users[i].IsActive)
		idx += 4
	}

	query := fmt.Sprintf(`
		INSERT INTO users (user_id, username, team_name, is_active)
		VALUES %s
		ON CONFLICT (user_id)
		DO UPDATE
		SET team_name = EXCLUDED.team_name,
			username = EXCLUDED.username,
			is_active = EXCLUDED.is_active;
		`, 
		strings.Join(placeholders, ","),
	)

	_, err := tx.ExecContext(ctx, query, args...)
	if err != nil {
		return fmt.Errorf("%s exec: %w", op, err)
	}


	return nil
}

// #### team/get ####

func (s *Storage) GetTeam(ctx context.Context, teamName string) (*models.Team, error) {
	const op = "storage.postgres.GetTeamByName"

	var team models.Team
	var user models.User

	err := s.db.QueryRowContext(ctx, `SELECT team_name FROM teams WHERE team_name=$1`, teamName).Scan(&team.TeamName)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, response.ErrNotFound
		}

		return nil, fmt.Errorf("%s: %w", op, err)
	}

	rows, err := s.db.QueryContext(ctx, `SELECT user_id, username, is_active FROM users WHERE team_name=$1`, teamName)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	defer rows.Close()

	for rows.Next() {
		err := rows.Scan(&user.UserID, &user.Username, &user.IsActive)
		if err != nil {
			return nil, fmt.Errorf("%s: %w", op, err)
		}

		user.TeamName = teamName

		team.Members = append(team.Members, user)
	}
	
	return &team, nil
}

func (s *Storage) SetIsActive(ctx context.Context, userID string, isActive bool) (*models.User, error) {
	const op = "storage.postgres.SetIsActive"

	var user models.User
	var isActiveDB bool

	err := s.db.QueryRowContext(ctx, `SELECT is_active FROM users WHERE user_id=$1`, userID).Scan(&isActiveDB)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, response.ErrNotFound
		}

		return nil, fmt.Errorf("%s: %w", op, err)
	}

	if isActiveDB != isActive {
		_, err := s.db.ExecContext(ctx, `UPDATE users SET is_active=$1 WHERE user_id=$2`, isActive, userID)
		if err != nil {
			return nil, fmt.Errorf("%s: %w", op, err)
		}
	}

	err = s.db.QueryRowContext(ctx, 
	`SELECT username, team_name, is_active 
	FROM users WHERE user_id=$1`,userID).
	Scan(
		&user.Username, 
		&user.TeamName, 
		&user.IsActive,
	)

	user.UserID = userID

	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return &user, nil

}