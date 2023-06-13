// pkg/webserver/webserver.go

package webserver

import (
	"embed"
	"flag"
	"fmt"
	"io/fs"
	"log"
	"net/http"
	"os"
	"path/filepath"

	handlers "github.com/devcashflow/database-server/pkg/api"
	"github.com/devcashflow/database-server/pkg/database"
	middlewares "github.com/devcashflow/database-server/pkg/middlewares"

	"github.com/go-chi/jwtauth/v5"
	"github.com/joho/godotenv"
	"github.com/rs/cors"

	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/chi/v5"
)

const (
	// token duration for a signed in session in hours.
	tokenDuration = 24
)

var (
	dev          = flag.Bool("dev", false, "development mode")
	door         = flag.Int("door", 3000, "the door the server is going to run. (Default 3000)")
	databaseType = flag.String("databaseType", "postgres", "the door the server is going to run. (Default PostgreSQL)")
	tokenAuth    *jwtauth.JWTAuth
	//go:embed www/*
	www embed.FS
)

func Start() {
	// Initialize database
	flag.Parse()
	var fileServer http.Handler
	if *dev {
		// load env file
		err := godotenv.Load()
		if err != nil {
			log.Fatal("Error loading .env file")
		}

		// Use local directory in development
		exe, err := os.Executable()
		if err != nil {
			log.Fatalf("Failed to determine executable path: %v", err)
		}

		// assumes running from cmd/database-server/main.go
		exeDir := filepath.Dir(exe)
		wwwDir := filepath.Join(exeDir, "..", "..", "pkg", "webserver", "www")
		log.Printf("Server wwwDir on %s", wwwDir)
		fileServer = http.FileServer(http.Dir(wwwDir))
	} else {
		fmt.Printf("[Production]\n")
		// Use embedded files in production
		subFS, err := fs.Sub(www, "www")
		if err != nil {
			log.Fatal(err)
		}
		fileServer = http.FileServer(http.FS(subFS))
	}

	// get env variables
	config := database.Config{
		User:    os.Getenv("MYSQL_USER"),
		Pass:    os.Getenv("MYSQL_PASS"),
		Host:    os.Getenv("MYSQL_HOST"),
		Port:    os.Getenv("MYSQL_DOOR"),
		Name:    os.Getenv("MYSQL_DATABASE_NAME"),
		SSLMODE: os.Getenv("SSLMODE"),
		Type:    *databaseType,
	}
	db, err := database.Connect(config)
	if err != nil {
		log.Fatalf("Failed to Connect to database server: %v", err)
	}

	// Initialize handlers with db instance
	server, err := handlers.New(db)
	if err != nil {
		log.Fatalf("Failed to create handlers: %v", err)
	}

	// Define routes
	router := chi.NewRouter()
	// Add Cors
	// XXX make it less public
	router.Use(cors.New(cors.Options{
		AllowedOrigins: []string{"*"},
		AllowedMethods: []string{
			http.MethodPost,
			http.MethodGet,
		},
		AllowedHeaders:   []string{"*"},
		AllowCredentials: false,
	}).Handler)

	router.Use(middleware.RequestID)
	router.Use(middleware.Logger)
	router.Use(middleware.Recoverer)

	router.Use(middlewares.DBConnected(db))

	router.Post("/create-email", server.HandleCreateEmail)
	router.Get("/list-emails", server.HandleListEmails)
	router.Get("/version", server.HandleVersion)

	router.Handle("/*", fileServer)

	// Start server
	log.Printf("Server listening on port %d", *door)
	log.Fatal(http.ListenAndServe(fmt.Sprintf("localhost:%d", *door), router))
}
