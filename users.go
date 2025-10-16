package main

import(
	"net/http"
	"encoding/json"
	"log"
)

func (a *apiConfig) registerUserHandler(w http.ResponseWriter, r *http.Request){	
	type parameters struct{
		Email string `json:"email"`
	}
	defer r.Body.Close()
	
	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	if err := decoder.Decode(&params); err != nil {
		log.Printf("Error decoding parameters: %s", err)
		w.WriteHeader(500)
		return
	}

	usr, err := a.Db.CreateUser(r.Context(), params.Email)
	if err != nil {
		log.Printf("Error creating new user: %s", err)
		w.WriteHeader(500)
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
