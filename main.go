package main

import(
	"fmt"
	"net/http"
	"log"
)

func healthzHandler(w http.ResponseWriter, r *http.Request){
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(http.StatusText(http.StatusOK)))
}

func main(){
	const root = "."
	const port = "8080"

	mux := http.NewServeMux()
	mux.Handle("/app/", http.StripPrefix("/app", http.FileServer(http.Dir(root))))
	mux.HandleFunc("/healthz", healthzHandler)

	srv := &http.Server{
		Handler: mux, 
		Addr: ":" + port,
	}
	
	fmt.Println("Starting server at 8080...")
	if err := srv.ListenAndServe(); err != nil{
		log.Fatalf("server stopped: %s\n", err)
	}
}

