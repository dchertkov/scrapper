package api

import (
	"net/http"

	"github.com/dchertkov/scrapper/pkg/api/logger"
	"github.com/dchertkov/scrapper/pkg/api/service"
	"github.com/dchertkov/scrapper/pkg/api/stat"
	"github.com/dchertkov/scrapper/pkg/types"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
)

type API struct {
	service types.ServiceStore
	stat    types.StatStore
}

func NewHandler(
	service types.ServiceStore,
	stat types.StatStore,
) *API {
	return &API{
		service: service,
		stat:    stat,
	}
}

func (api *API) Handler() http.Handler {
	r := chi.NewRouter()
	r.Use(middleware.Recoverer)
	r.Use(logger.Middleware)

	r.Route("/service", func(r chi.Router) {
		r.Route("/{host}", func(r chi.Router) {
			r.Use(service.ServiceCtx(api.service))

			r.With(stat.Collect(api.stat)).Get("/", service.HandleService())
			r.Get("/stat", stat.HandleStat(api.stat))
		})

		r.Get("/time/min", service.HandleMinTime(api.service))
		r.Get("/time/max", service.HandleMaxTime(api.service))
	})

	return r
}
