package service

import (
	"net/http"

	"github.com/dchertkov/scrapper/pkg/api/context"
	"github.com/dchertkov/scrapper/pkg/api/response"
	"github.com/dchertkov/scrapper/pkg/types"

	"github.com/go-chi/chi"
)

const hostKey = "host"

func ServiceCtx(service types.ServiceStore) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			s, err := service.Find(chi.URLParam(r, hostKey))
			if err != nil {
				response.NotFound(w, err)
				return
			}
			ctx := context.ToService(r.Context(), s)

			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func HandleService() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		s, ok := context.FromService(r.Context())
		if !ok {
			response.NotFound(w, nil)
			return
		}
		response.JSON(w, s, http.StatusOK)
	}
}

func HandleMinTime(service types.ServiceStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		s, err := service.MinTime()
		if err != nil {
			response.NotFound(w, err)
			return
		}
		response.JSON(w, s, http.StatusOK)
	}
}

func HandleMaxTime(service types.ServiceStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		s, err := service.MaxTime()
		if err != nil {
			response.NotFound(w, err)
			return
		}
		response.JSON(w, s, http.StatusOK)
	}
}
