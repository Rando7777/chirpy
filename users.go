package main

import(
	"net/http"
	"encoding/json"
	"fmt"
	"github.com/Rando7777/chirpy/internal/auth"
	"github.com/Rando7777/chirpy/internal/database"
)

func (a *apiConfig) registerUserHandler(w http.ResponseWriter, r *http.Request){	
	defer r.Body.Close()

	type parameters struct{
		Email string `json:"email"`
		Password string `json:"password"`
	}
	
	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	if err := decoder.Decode(&params); err != nil {
		respondWithError(w, http.StatusInternalServerError, fmt.Sprintf("Error decoding parameters: %s", err))
		return
	}
	hash, err := auth.HashPassword(params.Password)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, fmt.Sprintf("Error hasing password: %s", err))
	}
	usrParams := database.CreateUserParams{
		Email: params.Email,
		HashedPassword: hash,
	}
	usr, err := a.Db.CreateUser(r.Context(), usrParams)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, fmt.Sprintf("Error creating new user: %s", err))
		return
	}
	jsonUser := User{
		ID: usr.ID,
		CreatedAt: usr.CreatedAt,
		UpdatedAt: usr.UpdatedAt,
		Email: usr.Email,
	}

	respondWithJSON(w, 201, jsonUser)
}

func (a *apiConfig) loginHandler(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	type parameters struct{
		Email string `json:"email"`
		Password string `json:"password"`
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	if err := decoder.Decode(&params); err != nil {
		respondWithError(w, http.StatusInternalServerError, fmt.Sprintf("Error decoding parameters: %s", err))
		return
	}
	u, err := a.Db.GetUserByEmail(r.Context(), params.Email)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Incorrect email or password")
		return
	}
	match, err := auth.CheckPasswordHash(params.Password, u.HashedPassword)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Incorrect email or password")
		return
	}
	if !match{
		respondWithError(w, http.StatusUnauthorized, "Incorrect email or password")
		return
	}
	
	jsonUser := User{
		ID: u.ID,
		CreatedAt: u.CreatedAt,
		UpdatedAt: u.UpdatedAt,
		Email: u.Email,
	}

	respondWithJSON(w, 200, jsonUser)
}
