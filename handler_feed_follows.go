package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"rssag/internal/database"
	"time"

	"github.com/go-chi/chi"
	"github.com/google/uuid"
)

func (apiCfg *apiConfig) handleCreateFeedFollow(response http.ResponseWriter, request *http.Request, user database.User) {

	type parameters struct {
		FeedId uuid.UUID `json:"feed_id"`
	}

	decoder := json.NewDecoder(request.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(response, 400, "JSON decoding error")
		return
	}
	feedFollow, err := apiCfg.DB.CreateFeedFollows(request.Context(), database.CreateFeedFollowsParams{
		ID:        uuid.New(),
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
		UserID:    user.ID,
		FeedID:    params.FeedId,
	})
	if err != nil {
		respondWithError(response, 400, fmt.Sprintf("Could not create feed follow %s", err))
		return
	}
	respondWithJSON(response, 201, dbFeedFollowsToFeedFollows(feedFollow))
}

func (apiCfg *apiConfig) handleGetUserFollows(response http.ResponseWriter, request *http.Request, user database.User) {
	follows, err := apiCfg.DB.GetFeedFollows(request.Context(), user.ID)
	if err != nil {
		respondWithError(response, 404, "User's feeds not found")
	}
	respondWithJSON(response, 200, dbFollowsToFollows(follows))
}

func (apiCfg *apiConfig) handleDeleteUserFollow(response http.ResponseWriter, request *http.Request, user database.User) {
	feedFollowsIdStr := chi.URLParam(request, "feedFollowID")
	feedFollowsID, err := uuid.Parse(feedFollowsIdStr)
	if err != nil {
		respondWithError(response, 400, "JSON decoding error")
		return
	}
	err = apiCfg.DB.DeleteFeedFollow(request.Context(), database.DeleteFeedFollowParams{
		ID:     feedFollowsID,
		UserID: user.ID,
	})
	if err != nil {
		respondWithError(response, 404, "User's feed follow not found")
		return
	}
	respondWithJSON(response, 204, struct{}{})
}
