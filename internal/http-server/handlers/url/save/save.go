package save

import (
	"net/http"

	"github.com/andrei-kozel/url-shortener/internal/lib/api/response"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"golang.org/x/exp/slog"
)

type Request struct {
	URL   string `json:"url" validate:"required,url"`
	Alias string `json:"alias,omitempty"`
}

type Response struct {
	response.Response
	Alias string `json:"alias,omitempty"`
}

type URLSaver interface {
	SaveURL(urlToSave string, alias string) (int64, error)
}

func New(log *slog.Logger, urlSaver URLSaver) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.url.save.New"

		log = log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		var req Request
		err := render.DecodeJSON(r.Body, &req)
		if err != nil {
			log.Error("failed to decode request", err)
			render.JSON(w, r, response.Error("failed to decode request"))
			return
		}

		log.Info("saving url...", "url", req.URL, "alias", req.Alias)

		id, err := urlSaver.SaveURL(req.URL, req.Alias)
		if err != nil {
			log.Error("failed to save url", err)
			render.JSON(w, r, response.Error("failed to save url"))
			return
		}

		render.JSON(w, r, Response{
			Response: response.OK(),
			Alias:    req.Alias,
		})
	}
}
