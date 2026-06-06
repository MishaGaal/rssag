package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"rssag/internal/database"
	"time"

	"github.com/google/uuid"
)

func (apiCfg *apiConfig) handleCreateFeed(response http.ResponseWriter, request *http.Request, user database.User) {

	type parameters struct {
		Name string `json:"name"`
		Url  string `json:"url"`
	}

	decoder := json.NewDecoder(request.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(response, 400, "JSON decoding error")
		return
	}
	feed, err := apiCfg.DB.CreateFeed(request.Context(), database.CreateFeedParams{
		ID:        uuid.New(),
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
		Name:      params.Name,
		Url:       params.Url,
		UserID:    user.ID,
	})
	if err != nil {
		respondWithError(response, 400, fmt.Sprintf("Could not create user %s", err))
		return
	}
	respondWithJSON(response, 201, dbFeedToFeed(feed))
}

func (apiCfg *apiConfig) handleGetFeedByApiKeyAndFeedId(
	response http.ResponseWriter,
	request *http.Request,
	user database.User) {
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

	feed, err := apiCfg.DB.GetFeedById(request.Context(), params.FeedId)
	if err != nil {
		respondWithError(response, 404, "Feed not found")
	}
	respondWithJSON(response, 200, dbFeedToFeed(feed))
}

func (apiCfg *apiConfig) handleGetAllFeeds(
	response http.ResponseWriter,
	request *http.Request,
	user database.User) {

	feeds, err := apiCfg.DB.GetAllFeeds(request.Context())
	if err != nil {
		respondWithError(response, 404, "All feeds not found")
	}
	respondWithJSON(response, 200, dbFeedsToFeeds(feeds))
}

func (apiCfg *apiConfig) handleGetUserFeeds(
	response http.ResponseWriter,
	request *http.Request,
	user database.User) {

	feeds, err := apiCfg.DB.GetUserFeeds(request.Context(), user.ID)
	if err != nil {
		respondWithError(response, 404, "User's feeds not found")
	}
	respondWithJSON(response, 200, dbFeedsToFeeds(feeds))
}
