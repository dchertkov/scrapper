package stat

import (
	"net/http"

	"github.com/dchertkov/scrapper/pkg/api/context"
	"github.com/dchertkov/scrapper/pkg/api/response"
	"github.com/dchertkov/scrapper/pkg/types"
)

func Collect(stat types.StatStore) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			s, ok := context.FromService(r.Context())
			if !ok {
				response.NotFound(w, nil)
				return
			}

			stat.Add(s.Host)

			next.ServeHTTP(w, r)
		})
	}
}

func HandleStat(stat types.StatStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		s, ok := context.FromService(r.Context())
		if !ok {
			response.NotFound(w, nil)
			return
		}

		v, err := stat.Find(s.Host)
		if err != nil {
			response.NotFound(w, err)
			return
		}

		response.JSON(w, v, http.StatusOK)
	}
}
