package helper

import (
	"backend/generate"
	"backend/internal/models"
	"backend/internal/repository/database"
	"context"
	"errors"
	"log"
	"net/http"
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

func RegisterUserToDB(ctx context.Context, user *models.User, hashedPassword string, client *mongo.Client) *mongo.InsertOneResult {

	var userCollection *mongo.Collection = database.OpenCollection("users", client)

	count, err := userCollection.CountDocuments(ctx, bson.D{{Key: "email", Value: user.Email}})
	if err != nil {
		log.Printf("%v error retrieving data", err)
		return nil
	}

	if count > 0 {
		log.Printf("%v User already exists", err)
		return nil
	}

	user.UserID = bson.NewObjectID().Hex()
	user.CreatedAt = time.Now()
	user.UpdatedAt = time.Now()
	user.Password = hashedPassword

	result, err := userCollection.InsertOne(ctx, user)
	if err != nil {
		log.Printf("%v Failed to create user", err)
		return nil
	}

	return result

}

func LoginUserToDB(ctx context.Context, client *mongo.Client, userEmail string) models.User {

	var foundUser models.User
	var userCollection *mongo.Collection = database.OpenCollection("users", client)

	err := userCollection.FindOne(ctx, bson.D{{Key: "email", Value: userEmail}}).Decode(&foundUser)
	if err != nil {
		log.Printf("%v Failed to create user", err)
		return models.User{}
	}

	return foundUser

}

func RefreshTokenDB(ctx context.Context, client *mongo.Client, claims generate.SignedDetails) (string, string) {

	var userCollection *mongo.Collection = database.OpenCollection("users", client)

	var user models.User
	err := userCollection.FindOne(ctx, bson.D{{Key: "user_id", Value: claims.UserId}}).Decode(&user)
	if err != nil {
		log.Printf("%v User not found", err)
		return "", ""
	}

	newToken, newRefreshToken, _ := generate.GenerateAllToken(user.Email, user.FirstName, user.LastName, user.Role, user.UserID)
	err = UpdateAllToken(ctx, user.UserID, newRefreshToken, client)
	if err != nil {
		log.Printf("%v Error updating token", err)
		return "", ""
	}

	return newToken, newRefreshToken

}

func UpdateAllToken(ctx context.Context, userId, refreshToken string, client *mongo.Client) error {

	updateAt, _ := time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))

	updateData := bson.M{
		"$set": bson.M{
			"refresh_token": refreshToken,
			"update_at":     updateAt,
		},
	}

	var movieCollection *mongo.Collection = database.OpenCollection("users", client)

	_, err := movieCollection.UpdateOne(ctx, bson.M{"user_id": userId}, updateData)
	return err
}

func GetUserIdFromContext(r *http.Request) (string, error) {
	userId := r.Context().Value(generate.UserIdKey).(string)

	if userId == "" {
		return "", errors.New("userId does not exists in this context")
	}

	return userId, nil

}

func GetRoleFromContext(r *http.Request) (string, error) {
	role := r.Context().Value(generate.RoleKey)

	if role == "" {
		return "", errors.New("role does not exists in this context")
	}

	memberRole, ok := role.(string)

	if !ok {
		return "", errors.New("unable to retrieve userId")
	}

	return memberRole, nil

}
