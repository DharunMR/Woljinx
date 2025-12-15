package router

import (
	"backend/internal/api/handlers"
	"net/http"

	"go.mongodb.org/mongo-driver/v2/mongo"
)

func UserRoute(client *mongo.Client) *http.ServeMux {
	mux := http.NewServeMux()

	mux.Handle("POST /register", handlers.RegisterUser(client))
	mux.Handle("POST /login", handlers.LoginUser(client))
	mux.Handle("POST /logout", handlers.LogoutUser(client))
	mux.Handle("POST /refresh", handlers.RefreshToken(client))
	return mux
}
