package service

import (
	"avito-test-assignment-backend/internal/models"
	"avito-test-assignment-backend/pkg/response"
	"context"
	"database/sql"
	"fmt"
)

type Store interface {
	BeginTx(ctx context.Context) (*sql.Tx, error)
	InsertTeamTx(ctx context.Context, tx *sql.Tx, teamName string) (int64, error)
	UpsertUsersTx(ctx context.Context, tx *sql.Tx, teamName string, users []models.User) error
}

type TeamService struct {
	store Store
}

func NewTeamService(store Store) *TeamService {
    return &TeamService{store: store}
}


func (s *TeamService) AddTeamService(ctx context.Context, t models.Team) error {
	const op = "service.AddTeam"
	
	tx, err := s.store.BeginTx(ctx)
	if err != nil {
        return fmt.Errorf("%s: begin tx: %w", op, err)
    }

	defer func() {
		if p := recover(); p != nil {
            _ = tx.Rollback()
            panic(p)
		}
	}()

	rows, err := s.store.InsertTeamTx(ctx, tx, t.TeamName)
    if err != nil {
        _ = tx.Rollback()
        return fmt.Errorf("%s: insert team: %w", op, err)
    }
    if rows == 0 {
        _ = tx.Rollback()
        return fmt.Errorf("%s: %w", op, response.ErrTeamExists) 
    }
	
	if err := s.store.UpsertUsersTx(ctx, tx, t.TeamName, t.Members); err != nil {
        _ = tx.Rollback()
        return fmt.Errorf("%s: upsert users: %w", op, err)
    }

    if err := tx.Commit(); err != nil {
        return fmt.Errorf("%s: commit: %w", op, err)
    }

    return nil
}