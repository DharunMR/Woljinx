package helper

import (
	"backend/internal/models"
	"backend/internal/repository/database"
	"context"
	"fmt"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

func GetAllMovies(ctx context.Context, movies []models.Movie, client *mongo.Client) ([]models.Movie, error) {

	var movieCollection *mongo.Collection = database.OpenCollection("movies", client)

	cursor, err := movieCollection.Find(ctx, bson.D{})
	if err != nil {
		fmt.Println(err)
		return nil, fmt.Errorf("%v error retrieving data", err)
	}
	defer cursor.Close(ctx)

	if err = cursor.All(ctx, &movies); err != nil {
		return nil, fmt.Errorf("%v Failed to decode movies", err)
	}

	return movies, nil

}

func GetOneMovie(ctx context.Context, movieid string, client *mongo.Client) (models.Movie, error) {

	var movieCollection *mongo.Collection = database.OpenCollection("movies", client)

	var movie models.Movie

	err := movieCollection.FindOne(ctx, bson.D{{Key: "imdb_id", Value: movieid}}).Decode(&movie)
	if err != nil {
		return models.Movie{}, fmt.Errorf("%v error retrieving data", err)
	}

	return movie, nil
}

func AddOneMovie(ctx context.Context, movie models.Movie, client *mongo.Client) (*mongo.InsertOneResult, error) {

	var movieCollection *mongo.Collection = database.OpenCollection("movies", client)
	result, err := movieCollection.InsertOne(ctx, movie)
	if err != nil {
		return nil, fmt.Errorf("error retrieving data: %w", err)
	}

	return result, nil
}
