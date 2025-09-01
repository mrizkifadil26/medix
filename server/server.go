package server

import (
	"net/http"

	"github.com/mrizkifadil26/medix/utils/logger"
)

func logRequest(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		logger.Info("[" + r.Method + "] " + r.URL.Path)
		handler.ServeHTTP(w, r)
	})
}

func Serve(distPath string, port string) error {
	logger.Info("Serving " + distPath + " at http://localhost:" + port)

	fs := http.FileServer(http.Dir(distPath))
	return http.ListenAndServe(":"+port, logRequest(fs))
}
