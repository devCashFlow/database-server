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
	"github.com/devcashflow/database-server/pkg/database/mysql"
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
	dev       = flag.Bool("dev", false, "development mode")
	door      = flag.Int("door", 8080, "the door the server is going to run. Default 8080")
	tokenAuth *jwtauth.JWTAuth
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
	mysqlUser := os.Getenv("MYSQL_USER")
	mysqlPass := os.Getenv("MYSQL_PASS")
	mysqlHost := os.Getenv("MYSQL_HOST")
	mysqlDoor := os.Getenv("MYSQL_DOOR")
	mysqlDBName := os.Getenv("MYSQL_DATABASE_NAME")

	db, err := mysql.Connect(mysqlUser, mysqlPass, mysqlHost, mysqlDoor, mysqlDBName)
	if err != nil {
		log.Fatal(err)
	}

	// Initialize server with db instance
	server, err := handlers.New(db)
	if err != nil {
		log.Fatalf("Failed to create server: %v", err)
	}

	// Define routes
	router := chi.NewRouter()
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

	router.Handle("/*", fileServer)

	// Start server
	log.Printf("Server listening on port %d", *door)
	log.Fatal(http.ListenAndServe(fmt.Sprintf("localhost:%d", *door), router))
}
