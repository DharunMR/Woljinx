package handlers

import (
	"backend/generate"
	"backend/internal/api/helper"
	"backend/internal/models"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/go-playground/validator/v10"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

func RegisterUser(client *mongo.Client) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {

		var ctx, cancel = context.WithTimeout(r.Context(), 100*time.Second)
		defer cancel()

		var user models.User

		err := json.NewDecoder(r.Body).Decode(&user)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		validate := validator.New()

		if err := validate.Struct(user); err != nil {
			http.Error(w, "Validation failed: "+err.Error(), http.StatusBadRequest)
			return
		}

		hashedPassword, err := generate.HashPassword(user.Password)
		if err != nil {
			http.Error(w, "Hashing failed: "+err.Error(), http.StatusBadRequest)
			return
		}

		result := helper.RegisterUserToDB(ctx, &user, hashedPassword, client)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(result)

	}
}

func LoginUser(client *mongo.Client) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {

		var ctx, cancel = context.WithTimeout(r.Context(), 100*time.Second)
		defer cancel()

		var userLogin models.User

		err := json.NewDecoder(r.Body).Decode(&userLogin)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		foundUser := helper.LoginUserToDB(ctx, client, userLogin.Email)
		err = generate.VerifyPassword(userLogin.Password, foundUser.Password)
		if err != nil {
			log.Printf("%v Invalid email or Password", err)
			return
		}

		token, refreshToken, err := generate.GenerateAllToken(foundUser.Email, foundUser.FirstName, foundUser.LastName, foundUser.Role, foundUser.UserID)
		if err != nil {
			log.Printf("%v Failed to generate tokens", err)
			return
		}

		err = helper.UpdateAllToken(ctx, foundUser.UserID, refreshToken, client)
		if err != nil {
			log.Printf("%v Failed to update tokens", err)
			return
		}

		http.SetCookie(w, &http.Cookie{
			Name:     "access_token",
			Value:    token,
			Path:     "/",
			MaxAge:   86400,
			Secure:   true,
			HttpOnly: true,
			SameSite: http.SameSiteNoneMode,
		})

		http.SetCookie(w, &http.Cookie{
			Name:     "refresh_token",
			Value:    refreshToken,
			Path:     "/",
			MaxAge:   604800,
			Secure:   true,
			HttpOnly: true,
			SameSite: http.SameSiteNoneMode,
		})

		json.NewEncoder(w).Encode(models.UserResponse{
			UserId:          foundUser.UserID,
			FirstName:       foundUser.FirstName,
			LastName:        foundUser.LastName,
			Email:           foundUser.Email,
			Role:            foundUser.Role,
			FavouriteGenres: foundUser.FavouriteGenres,
		})
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
	}
}

func LogoutUser(client *mongo.Client) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {

		var ctx, cancel = context.WithTimeout(r.Context(), 100*time.Second)
		defer cancel()

		var UserLogout struct {
			UserId string `json:"user_id"`
		}

		err := json.NewDecoder(r.Body).Decode(&UserLogout)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		err = helper.UpdateAllToken(ctx, UserLogout.UserId, "", client)
		if err != nil {
			fmt.Printf("Couldnt clean token %v", err)
			return
		}

		http.SetCookie(w, &http.Cookie{
			Name:     "access_token",
			Value:    "",
			Path:     "/",
			MaxAge:   -1,
			Secure:   true,
			HttpOnly: true,
			SameSite: http.SameSiteNoneMode,
		})

		http.SetCookie(w, &http.Cookie{
			Name:     "refresh_token",
			Value:    "",
			Path:     "/",
			MaxAge:   -1,
			Secure:   true,
			HttpOnly: true,
			SameSite: http.SameSiteNoneMode,
		})

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

	}
}

func RefreshToken(client *mongo.Client) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {

		var ctx, cancel = context.WithTimeout(r.Context(), 100*time.Second)
		defer cancel()

		refreshToken, err := r.Cookie("refresh_token")
		if err != nil {
			fmt.Printf("Couldnt retrieve token from cookie%v", err)
			return
		}

		claims, err := generate.ValidateToken(refreshToken.Value)
		if err != nil {
			fmt.Println("error", err.Error())
			fmt.Printf("Invalid or expired refresh token %v", err)
		}

		newToken, newRefreshToken := helper.RefreshTokenDB(ctx, client, *claims)

		http.SetCookie(w, &http.Cookie{
			Name:     "access_token",
			Value:    newToken,
			Path:     "/",
			Domain:   "localhost",
			MaxAge:   86400, // 1 day
			Secure:   true,
			HttpOnly: true,
		})

		http.SetCookie(w, &http.Cookie{
			Name:     "access_token",
			Value:    newRefreshToken,
			Path:     "/",
			Domain:   "localhost",
			MaxAge:   604800, // 1 day
			Secure:   true,
			HttpOnly: true,
		})

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

	}
}
