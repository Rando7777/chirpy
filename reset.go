package main

import(
	"net/http"
	"log"
)


func (a *apiConfig) metricsResetHandler(w http.ResponseWriter, r *http.Request){
	if a.Platform != "dev"{
		w.WriteHeader(http.StatusForbidden)
	    w.Write([]byte(http.StatusText(http.StatusForbidden)))
		return
	}

	a.fileserverHits.Store(0)
	if err := a.Db.DeleteUsers(r.Context()); err != nil {
		log.Printf("Error decoding parameters: %s", err)
		w.WriteHeader(500)
		return
	}

	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(http.StatusText(http.StatusOK)))
}
