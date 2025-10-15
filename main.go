package main

import(
	"fmt"
	"net/http"
	"log"
	"sync/atomic"
	_ "github.com/lib/pq"
	"github.com/joho/godotenv"
	"os"
	"database/sql"
	"github.com/Rando7777/chirpy/internal/database"
)

type apiConfig struct{
	fileserverHits atomic.Int32
	Db *database.Queries
}


func main(){
	godotenv.Load()
	const root = "."
	const port = "8080"
	dbURL := os.Getenv("DB_URL")
	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatalf("Error opening db connection: %s", err)
	}

	apiCfg := &apiConfig{
		Db: database.New(db),
	}

	mux := http.NewServeMux()
	
	fileHandler := apiCfg.middlewareMetricsInc(http.StripPrefix("/app", http.FileServer(http.Dir(root))))
	mux.Handle("/app/", fileHandler)
	
	mux.HandleFunc("GET /api/healthz", healthzHandler)
	mux.HandleFunc("POST /api/validate_chirp", validateChirpHandler)

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



