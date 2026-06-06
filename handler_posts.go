package main

import (
	"net/http"
	"rssag/internal/database"
	"strconv"
)

func (apiCfg *apiConfig) handleGetUserPosts(response http.ResponseWriter, request *http.Request, user database.User) {
	limitStr := request.URL.Query().Get("limit")
	limit := 10
	if specifiedLimit, err := strconv.Atoi(limitStr); err == nil {
		limit = specifiedLimit
	}

	posts, err := apiCfg.DB.GetPostsForUser(request.Context(), database.GetPostsForUserParams{
		UserID: user.ID,
		Limit:  int32(limit),
	})
	if err != nil {
		respondWithError(response, http.StatusInternalServerError, "Couldn't get posts for user")
		return
	}

	respondWithJSON(response, http.StatusOK, databasePostsToPosts(posts))
}
