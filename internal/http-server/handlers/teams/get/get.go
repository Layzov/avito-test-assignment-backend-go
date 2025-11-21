package get

import (
	"avito-test-assignment-backend/api"
	"avito-test-assignment-backend/pkg/response"
	"avito-test-assignment-backend/pkg/sl"
	"context"
	"errors"
	"log/slog"
	"net/http"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/render"
)

type TeamGetter interface {
	GetTeamService(ctx context.Context, teamName string) (*api.Team, error)
}


type Response struct {
	response.Response
	api.Team `json:"team,omitempty"`
}

func New(log *slog.Logger,teamGetter TeamGetter) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.teams.get.get.New"
		log = log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		teamName := chi.URLParam(r, "team_name")

		team, err :=  teamGetter.GetTeamService(r.Context(), teamName)

		if errors.Is(err, response.ErrNotFound) {
			log.Error("team_name not found")
			w.WriteHeader(http.StatusNotFound)
			render.JSON(w, r, response.Error(string(response.NOT_FOUND),"resource not found"))

			return
		}

		if err != nil {
			log.Error("Failed to get team", sl.Err(err))
			w.WriteHeader(http.StatusInternalServerError)
			render.JSON(w, r, response.Error(string(response.FAILED_REQUEST),"failed to get team"))

			return
		}

		log.Info("Team was received ", slog.Any("team", team))

		responseOK(w, r, team)
	}
}

func responseOK(w http.ResponseWriter, r *http.Request, team *api.Team) {
	render.JSON(w, r, Response{
		Team: *team,
	})
}