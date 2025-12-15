package router

import (
	"backend/internal/api/handlers"
	"net/http"

	"go.mongodb.org/mongo-driver/v2/mongo"
)

func MovieRoute(client *mongo.Client) *http.ServeMux {
	mux := http.NewServeMux()

	mux.Handle("GET /movies", handlers.GetMoviesHandler(client))
	mux.Handle("POST /addmovie", handlers.AddMovieHandler(client))
	mux.Handle("GET /movie/{imdb_id}", handlers.GetOneMovieHandler(client))
	return mux
}
