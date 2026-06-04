package main

import "net/http"

func handleReadiness(response http.ResponseWriter, request *http.Request) {
	respondWithJSON(response, 200, struct{}{})
}
