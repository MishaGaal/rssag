package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"rssag/internal/database"
	"time"

	"github.com/google/uuid"
)

func (apiCfg *apiConfig) handleCreateUser(response http.ResponseWriter, request *http.Request) {

	type parameters struct {
		Name string `json:"name"`
	}

	decoder := json.NewDecoder(request.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(response, 400, "JSON decoding error")
		return
	}
	usr, err := apiCfg.DB.CreateUser(request.Context(), database.CreateUserParams{
		ID:        uuid.New(),
		CreatedAt: time.Now().UTC(),
		Name:      params.Name,
	})
	if err != nil {
		respondWithError(response, 400, fmt.Sprintf("Could not create user %s", err))
		return
	}
	respondWithJSON(response, 201, dbUserToUser(usr))
}

func (apiCfg *apiConfig) handleGetUserByApiKey(response http.ResponseWriter, request *http.Request, user database.User) {
	respondWithJSON(response, 200, dbUserToUser(user))
}
