package main

import (
	"crypto/tls"
	"embed"
	"fmt"
	"log"
	"net/http"
	"os"

	mw "backend/internal/api/middlewares"
	"backend/internal/api/router"

	"github.com/joho/godotenv"
)

var envFile embed.FS

func loadEnvFromEmbeddedFile() {
	content, err := envFile.ReadFile(".env")
	if err != nil {
		log.Fatalf("Error reading .env file: %v", err)
	}

	tempfile, err := os.CreateTemp("", ".env")
	if err != nil {
		log.Fatalf("Error creating temp .env file: %v", err)
	}
	defer os.Remove(tempfile.Name())

	_, err = tempfile.Write(content)
	if err != nil {
		log.Fatalf("Error writing to temp .env file: %v", err)
	}

	err = tempfile.Close()
	if err != nil {
		log.Fatalf("Error closing temp file: %v", err)
	}

	err = godotenv.Load(tempfile.Name())
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}
}

func main() {

	port := os.Getenv("API_PORT")

	cert := os.Getenv("CERT_FILE")
	key := os.Getenv("KEY_FILE")

	tlsConfig := &tls.Config{
		MinVersion: tls.VersionTLS12,
	}

	hppOptions := mw.HPPOptions{
		CheckQuery:                  true,
		CheckBody:                   true,
		CheckBodyOnlyForContentType: "application/x-www-form-urlencoded",
		Whitelist:                   []string{"sortBy", "sortOrder", "name", "age", "class"},
	}

	router := router.MainRouter()
	jwtMiddleware := mw.MiddlewaresExcludePaths(mw.JWTMiddleware, "/execs/login", "/execs/forgotpassword", "/execs/resetpassword/reset")
	secureMux := utils.ApplyMiddlewares(router, mw.SecurityHeaders, mw.Compression, mw.Hpp(hppOptions), mw.XSSMiddleware, jwtMiddleware, mw.ResponseTimeMiddleware, mw.Cors)

	server := &http.Server{
		Addr:      port,
		Handler:   secureMux,
		TLSConfig: tlsConfig,
	}

	// http2.ConfigureServer(server, &http2.Server{})

	fmt.Println("Server is running on port: " + port)

	err := server.ListenAndServeTLS(cert, key)
	if err != nil {
		log.Fatalln("Error starting the server", err)
	}

}
