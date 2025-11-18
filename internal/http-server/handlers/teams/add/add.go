package add

import (
	"avito-test-assignment-backend/api"
	"avito-test-assignment-backend/pkg/response"
	"avito-test-assignment-backend/pkg/sl"
	"log/slog"
	"net/http"

	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/render"
)

type TeamAdder interface {
	AddTeam(t api.Team) error
}

type Request struct {
	api.Team
}

type Response struct {
	response.Response
	api.Team
}

func New(log *slog.Logger, teamAdder TeamAdder) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.teams.post.add.New"
		log = log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		var req Request

		if err := render.DecodeJSON(r.Body, &req); err != nil {
			log.Error("Failed to decode request body", sl.Err(err))
			w.WriteHeader(http.StatusBadRequest)
			render.JSON(w, r, response.Error("Failed to decode request"))

			return
		}

		log.Info("Request body decoded", slog.Any("request", req))
		
		err := teamAdder.AddTeam(req.Team)
		if err != nil {
			log.Error("Failed to add team", sl.Err(err))
			w.WriteHeader(http.StatusInternalServerError)
			render.JSON(w, r, response.Error("Failed to add team"))

			return
		}

		log.Info("Team added", slog.Any("team", req.Team))

		responseOK(w, r, req.Team)
	}	

}

func responseOK(w http.ResponseWriter, r *http.Request, team api.Team) {
	render.JSON(w, r, Response{
		Response: response.OK(),
		Team: team,
	})
}
