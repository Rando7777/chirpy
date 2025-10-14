package main

import(
	"fmt"
	"net/http"
	"log"
	"sync/atomic"
)

type apiConfig struct{
	fileserverHits atomic.Int32
}


func main(){
	const root = "."
	const port = "8080"
	
	apiCfg := &apiConfig{}
	mux := http.NewServeMux()
	
	fileHandler := apiCfg.middlewareMetricsInc(http.StripPrefix("/app", http.FileServer(http.Dir(root))))
	mux.Handle("/app/", fileHandler)
	
	mux.HandleFunc("GET /api/healthz", healthzHandler)

	mux.HandleFunc("GET /admin/metrics", apiCfg.metricsHandler)
	mux.HandleFunc("POST /admin/reset", apiCfg.metricsResetHandler)

	srv := &http.Server{
		Handler: mux, 
		Addr: ":" + port,
	}
	
	fmt.Println("Starting server at 8080...")
	if err := srv.ListenAndServe(); err != nil{
		log.Fatalf("server stopped: %s\n", err)
	}
}



