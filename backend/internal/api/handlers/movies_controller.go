package handlers

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"backend/internal/api/helper"
	models "backend/internal/models"

	"github.com/go-playground/validator/v10"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

func GetMoviesHandler(client *mongo.Client) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {

		var ctx, cancel = context.WithTimeout(r.Context(), 100*time.Second)
		defer cancel()

		var movies []models.Movie
		movies, err := helper.GetAllMovies(ctx, movies, client)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		response := struct {
			Status string         `json:"status"`
			Count  int            `json:"count"`
			Data   []models.Movie `json:"data"`
		}{
			Status: "success",
			Count:  len(movies),
			Data:   movies,
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	}
}

func GetOneMovieHandler(client *mongo.Client) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {

		var ctx, cancel = context.WithTimeout(r.Context(), 100*time.Second)
		defer cancel()

		id := r.PathValue("imdb_id")

		movie, err := helper.GetOneMovie(ctx, id, client)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(movie)
	}
}

func AddMovieHandler(client *mongo.Client) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {

		var ctx, cancel = context.WithTimeout(r.Context(), 100*time.Second)
		defer cancel()

		var newMovie models.Movie
		err := json.NewDecoder(r.Body).Decode(&newMovie)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		var validate = validator.New()
		err = validate.Struct(newMovie)
		if err != nil {
			http.Error(w, "Validation failed: "+err.Error(), http.StatusBadRequest)
			return
		}

		result, err := helper.AddOneMovie(ctx, newMovie, client)
		if err != nil {
			http.Error(w, "Validation failed: "+err.Error(), http.StatusBadRequest)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(result)
	}
}
