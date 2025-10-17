package main

import(
	"encoding/json"
	"fmt"
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
		respondWithError(w, http.StatusInternalServerError, fmt.Sprintf("Error decoding parameters: %s", err))
		return
	}

	if len(params.Body) > 140{
		respondWithError(w, http.StatusBadRequest, "Chirp is too long")
		return
	}
	params.Body = profanityFilter(params.Body)
	
	c, err := a.Db.CreateChirp(r.Context(), database.CreateChirpParams{
		Body: params.Body,
		UserID: params.UserID,
	})
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, fmt.Sprintf("Error creating db entry: %s", err))
		return
	}
	chirp := Chirp{
		ID: c.ID,
		CreatedAt: c.CreatedAt,
		UpdatedAt: c.UpdatedAt,
		Body: c.Body,
		UserID: c.UserID,
	}
	respondWithJSON(w, http.StatusCreated, chirp)
}

func (a *apiConfig) getChirpsHandler(w http.ResponseWriter, r *http.Request){
	chirps, err := a.Db.GetChirps(r.Context())
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, fmt.Sprintf("Error getting chirps: %s", err))
		return
	}
	var res []Chirp
	for _, c := range chirps {
		res = append(res, Chirp{
			ID: c.ID,
			CreatedAt: c.CreatedAt,
			UpdatedAt: c.UpdatedAt,
			Body: c.Body,
			UserID: c.UserID,
		})
	}
	respondWithJSON(w, http.StatusOK, res)
}
	
func (a *apiConfig) getChirpByIdHandler(w http.ResponseWriter, r *http.Request){
	id, err := uuid.Parse(r.PathValue("chirp_id"))
	if err != nil {
		respondWithError(w, http.StatusBadRequest, fmt.Sprintf("Error parsing uuid: %s", err))
		return
	}
	c, err := a.Db.GetChirp(r.Context(), id)
	if err != nil {
		respondWithError(w, http.StatusNotFound, fmt.Sprintf("Error getting chirp: %s", err))
		return
	}
	chirp := Chirp{
		ID: c.ID,
		CreatedAt: c.CreatedAt,
		UpdatedAt: c.UpdatedAt,
		Body: c.Body,
		UserID: c.UserID,
	}
	respondWithJSON(w, http.StatusOK, chirp)
}
