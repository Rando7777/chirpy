package main

import(
	"encoding/json"
	"log"
	"net/http"
	"strings"
	"slices"
)


func profanityFilter(msg string) string{
	profaneWords := []string{"kerfuffle", "sharbert", "fornax"}
	replacement := "****"
	words := strings.Split(msg, " ")
	for i, w := range words {
		if slices.Contains(profaneWords, strings.ToLower(w)){
			words[i] = replacement
		}
	}
	return strings.Join(words, " ")
}

func validateChirpHandler(w http.ResponseWriter, r *http.Request){
	type parameters struct{
		Chirp string `json:"body"`
	}
	defer r.Body.Close()
	
	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	if err := decoder.Decode(&params); err != nil {
		log.Printf("Error decoding parameters: %s", err)
		w.WriteHeader(500)
		return
	}

	if len(params.Chirp) > 140{
		respondWithError(w, 400, "Chirp is too long")
		return
	}

	params.Chirp = profanityFilter(params.Chirp)
	type successResponse struct{
		Cleaned string `json:"cleaned_body"`
	}
	respondWithJSON(w, 200, successResponse{Cleaned: params.Chirp})
}
