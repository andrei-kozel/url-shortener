package save

import (
	"errors"
	"log/slog"
	"net/http"

	"github.com/andrei-kozel/url-shortener/internal/lib/api/response"
	"github.com/andrei-kozel/url-shortener/internal/lib/random"
	"github.com/andrei-kozel/url-shortener/internal/storage"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"github.com/go-playground/validator/v10"
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

const aliasLength = 8

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

		log.Info("request body decoded", "req", req)

		if err := validator.New().Struct(req); err != nil {
			validateErr := err.(validator.ValidationErrors)
			log.Error("validation failed", "error", validateErr)

			render.JSON(w, r, response.Error("validation failed"))
			render.JSON(w, r, response.ValidationError(validateErr))

			return
		}

		log.Info("saving url...", "url", req.URL, "alias", req.Alias)

		alias := req.Alias
		if alias == "" {
			alias = random.NewRandomString(aliasLength)
		}

		id, err := urlSaver.SaveURL(req.URL, alias)
		if errors.Is(err, storage.ErrURLExists) {
			log.Error("url already exist", err)

			render.JSON(w, r, response.Error("failed to save url"))

			return
		}
		if err != nil {
			log.Error("failed to save url", err)

			render.JSON(w, r, response.Error("failed to save url"))

			return
		}

		log.Info("url saved", "id", id)

		render.JSON(w, r, Response{
			Response: response.OK(),
			Alias:    alias,
		})
	}
}
