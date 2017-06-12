package alarmy

import "net/http"

func StartServer(addr string, router Router) error {
	handler, err := router.Setup()
	if err != nil {
		return err
	}

	return http.ListenAndServe(addr, handler)
}
