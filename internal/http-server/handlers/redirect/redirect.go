package redirect

import (
	"errors"
	"log/slog"
	"net/http"

	"github.com/andrei-kozel/url-shortener/internal/storage"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
)

type URLGetter interface {
	GetURL(alias string) (string, error)
}

func New(log *slog.Logger, urlGetter URLGetter) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handler.url.redirect.New"

		log = log.With("op", op)

		alias := chi.URLParam(r, "alias")
		if alias == "" {
			log.Info("alias is empty")
			render.JSON(w, r, "not found")
			return
		}

		resURL, err := urlGetter.GetURL(alias)
		if errors.Is(err, storage.ErrURLNotFound) {
			log.Info("url not found", "alias", alias)
			render.JSON(w, r, "not found")
			return
		}
		if err != nil {
			log.Info("failed to get url", err)
			render.JSON(w, r, "internsl error")
			return
		}

		log.Info("url found", "url", resURL)

		http.Redirect(w, r, resURL, http.StatusFound)
	}
}
