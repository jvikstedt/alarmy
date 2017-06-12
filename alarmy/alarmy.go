package alarmy

import (
	"net/http"
)

type Router interface {
	Setup() (http.Handler, error)
}

func StartServer(addr string, router Router) error {
	handler, err := router.Setup()
	if err != nil {
		return err
	}

	return http.ListenAndServe(addr, handler)
}
