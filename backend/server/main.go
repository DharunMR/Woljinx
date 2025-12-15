package main

import (
	"context"
	"crypto/tls"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	mw "backend/internal/api/middlewares"
	"backend/internal/api/router"
	"backend/internal/repository/database"
	"backend/pkg/utils"

	"github.com/joho/godotenv"
)

func main() {

	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	port := os.Getenv("API_PORT")

	// cert := os.Getenv("CERT_FILE")
	// key := os.Getenv("KEY_FILE")

	tlsConfig := &tls.Config{
		MinVersion: tls.VersionTLS12,
	}

	hppOptions := mw.HPPOptions{
		CheckQuery:                  true,
		CheckBody:                   true,
		CheckBodyOnlyForContentType: "application/x-www-form-urlencoded",
		Whitelist:                   []string{"sortBy", "sortOrder", "name", "age", "class"},
	}

	Ratelimiter := mw.NewRateLimiter(10, 30*time.Second)

	client := database.Connect()
	defer func() {
		err := client.Disconnect(context.Background())
		if err != nil {
			log.Fatalf("Failed to disconnect from MongoDB: %v", err)
		}
	}()

	router := router.MainRouter(client)
	jwtMiddleware := mw.MiddlewaresExcludePaths(mw.JWTMiddleware, "/movies", "/register", "/login", "/logout", "/genres", "/refresh")
	secureMux := utils.ApplyMiddlewares(router, jwtMiddleware, Ratelimiter.Middleware, mw.SecurityHeaders, mw.Compression, mw.Hpp(hppOptions), mw.XSSMiddleware, mw.ResponseTimeMiddleware, mw.Cors)

	server := &http.Server{
		Addr:      port,
		Handler:   secureMux,
		TLSConfig: tlsConfig,
	}

	// http2.ConfigureServer(server, &http2.Server{})

	fmt.Println("Server is running on port: " + port)

	err = server.ListenAndServe()
	if err != nil {
		log.Fatalln("Error starting the server", err)
	}

}
