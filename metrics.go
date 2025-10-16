package main

import(
	"net/http"
	"fmt"
)

func (a *apiConfig) metricsHandler(w http.ResponseWriter, r *http.Request){
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	content := fmt.Sprintf(`
<html>
<body>
	<h1>Welcome, Chirpy Admin</h1>
	<p>Chirpy has been visited %d times!</p>
</body>
</html>
	`, a.fileserverHits.Load())

	w.Write([]byte(content))
}

func (a *apiConfig) middlewareMetricsInc(next http.Handler) http.Handler{
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request){
		a.fileserverHits.Add(1)	
		next.ServeHTTP(w, r)
	})
}

