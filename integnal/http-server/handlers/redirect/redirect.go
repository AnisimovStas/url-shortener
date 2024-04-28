package redirect

import (
	"errors"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"log/slog"
	"net/http"
	resp "url-shortner/integnal/lib/api/response"
	"url-shortner/integnal/lib/logger/sl"
)

//go:generate mockery --name=URLGetter
type URLGetter interface {
	GetURL(alias string) (string, error)
}

func New(log *slog.Logger, URLGetter URLGetter) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const fn = "handlers.redirect.New"

		log = log.With(slog.String("fn", fn),
			slog.String("request_id", middleware.GetReqID(r.Context())))

		alias := chi.URLParam(r, "alias")
		if alias == "" {
			log.Error("alias is empty", sl.Err(errors.New("alias is empty")))

			render.JSON(w, r, resp.Error("alias is empty"))

			return
		}

		url, err := URLGetter.GetURL(alias)
		if err != nil {
			log.Error("failed to get url", sl.Err(err))

			render.JSON(w, r, resp.Error("failed to get url"))

			return
		}

		http.Redirect(w, r, url, http.StatusFound)

	}

}
