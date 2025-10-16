package main

import(
	"encoding/json"
	"log"
	"net/http"
	"strings"
	"slices"
	"github.com/google/uuid"
	"github.com/Rando7777/chirpy/internal/database"
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

func (a *apiConfig) createChirpHandler(w http.ResponseWriter, r *http.Request){
	type parameters struct{
		Body string `json:"body"`
		UserID uuid.UUID `json:"user_id"`
	}
	defer r.Body.Close()
	
	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	if err := decoder.Decode(&params); err != nil {
		log.Printf("Error decoding parameters: %s", err)
		w.WriteHeader(500)
		return
	}

	if len(params.Body) > 140{
		respondWithError(w, 400, "Chirp is too long")
		return
	}
	params.Body = profanityFilter(params.Body)
	
	c, err := a.Db.CreateChirp(r.Context(), database.CreateChirpParams{
		Body: params.Body,
		UserID: params.UserID,
	})
	if err != nil {
		log.Printf("Error creating db entry: %s", err)
		w.WriteHeader(500)
		return
	}
	chirp := Chirp{
		ID: c.ID,
		CreatedAt: c.CreatedAt,
		UpdatedAt: c.UpdatedAt,
		Body: c.Body,
		UserID: c.UserID,
	}
	respondWithJSON(w, 201, chirp)
}
