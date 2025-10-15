package main

import(
	"net/http"
	"log"
	"encoding/json"
)

func respondWithError (w http.ResponseWriter, code int, msg string){
	type errorResponse struct{
		Error string `json:"error"`
	}
	respErr := errorResponse{Error: msg}	
	dat, err := json.Marshal(respErr)
	if err != nil {
		log.Printf("Error marshalling JSON: %s", err)
		w.WriteHeader(500)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(dat)
}

func respondWithJSON(w http.ResponseWriter, code int, payload any){
	dat, err := json.Marshal(payload)
	if err != nil {
		log.Printf("Error marshalling JSON: %s", err)
		w.WriteHeader(500)
		return
	}
	
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(dat)
}
