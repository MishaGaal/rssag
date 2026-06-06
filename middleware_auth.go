package main

import (
	"fmt"
	"net/http"
	"rssag/internal/auth"
	"rssag/internal/database"
)

type authHandler func(w http.ResponseWriter, r *http.Request, user database.User)

func (cfg *apiConfig) middlewareAuth(handle authHandler) http.HandlerFunc {
	return func(response http.ResponseWriter, request *http.Request) {
		apiKey, err := auth.GetApiKey(request.Header)
		if err != nil {
			respondWithError(response, 401, fmt.Sprintf("Unauthorized: %s", err))
			return
		}
		usr, err := cfg.DB.GetUsrByAPIKey(request.Context(), apiKey)
		if err != nil {
			respondWithError(response, 400, fmt.Sprintf("Couldn't get user: %s", err))
		}
		handle(response, request, usr)
	}
}
