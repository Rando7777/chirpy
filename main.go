package main

import(
	"fmt"
	"net/http"
	"log"
	"sync/atomic"
	"strconv"
)

type apiConfig struct{
	fileserverHits atomic.Int32
}

func (a *apiConfig) middlewareMetricsInc(next http.Handler) http.Handler{
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request){
		a.fileserverHits.Add(1)	
		next.ServeHTTP(w, r)
	})
}


func main(){
	const root = "."
	const port = "8080"
	fileHandler := http.StripPrefix("/app", http.FileServer(http.Dir(root)))
	
	apiCfg := &apiConfig{}

	mux := http.NewServeMux()
	mux.Handle("/app/", apiCfg.middlewareMetricsInc(fileHandler))
	mux.HandleFunc("/healthz", healthzHandler)
	mux.HandleFunc("/metrics", apiCfg.metricsHandler)
	mux.HandleFunc("/reset", apiCfg.metricsResetHandler)

	srv := &http.Server{
		Handler: mux, 
		Addr: ":" + port,
	}
	
	fmt.Println("Starting server at 8080...")
	if err := srv.ListenAndServe(); err != nil{
		log.Fatalf("server stopped: %s\n", err)
	}
}


func healthzHandler(w http.ResponseWriter, r *http.Request){
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(http.StatusText(http.StatusOK)))
}

func (a *apiConfig) metricsHandler(w http.ResponseWriter, r *http.Request){
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	
	countStr := strconv.Itoa(int(a.fileserverHits.Load()))
	w.Write([]byte("Hits: " + countStr))
}

func (a *apiConfig) metricsResetHandler(w http.ResponseWriter, r *http.Request){
	a.fileserverHits.Store(0)

	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(http.StatusText(http.StatusOK)))
}







