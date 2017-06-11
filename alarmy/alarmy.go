package alarmy

import (
	"net/http"
)

func StartServer(addr string) error {
	router, err := setupRouter()
	if err != nil {
		return err
	}

	return http.ListenAndServe(addr, router)
}
