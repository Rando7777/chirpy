package main

import(
	"net/http"
)


func (a *apiConfig) metricsResetHandler(w http.ResponseWriter, r *http.Request){
	a.fileserverHits.Store(0)

	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(http.StatusText(http.StatusOK)))
}
