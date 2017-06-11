package alarmy

import (
	"net/http"

	"github.com/pressly/chi"
)

func setupRouter() (http.Handler, error) {
	r := chi.NewRouter()
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("welcome"))
	})
	return r, nil
}
