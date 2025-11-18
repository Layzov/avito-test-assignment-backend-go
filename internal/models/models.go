package models

import "time"

type PRStatus string

const (
	PR_OPEN PRStatus = "OPEN"
	PR_MERGED PRStatus = "MERGED"
)

type User struct {
	UserID string `db:"user_id"`
	Username string `db:"username"`
	TeamName string `db:"team_name"`
	IsActive bool `db:"is_active"`
}

type PullRequest struct {
	PullRequestID 	string `db:"pull_request_id"`
	PullRequestName	string `db:"pull_request_name"`
	AuthorID 		string `db:"author_id"`
	Status 			PRStatus `db:"status"`
	// ??? Хз пока что как сделать ревьюерс
	CreatedAt 		time.Time `db:"created_at"`
	MergedAt 		*time.Time `db:"merged_at,omitempty"` // TODO: подумать над омитемпти
}

type Team struct {
	TeamName string `db:"team_name"`
	Members []User 
}