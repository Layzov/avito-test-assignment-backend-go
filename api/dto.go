package api

import ()

type TeamMember struct {
	UserID   string `json:"user_id"`
	Username string `json:"username"`
	IsActive bool   `json:"is_active"`
}

type PullRequestShort struct {
	PullRequestID  	string `json:"pull_request_id"`
	PullRequestName string `json:"pull_request_name"`
	AuthorID       	string `json:"author_id"`
	Status 			string `json:"status"`
}

//возможно потом в домен закину
type Team struct {
	TeamName    string       `json:"team_name"`
	Members     []TeamMember `json:"members"`
}