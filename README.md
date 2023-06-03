# Database Server

This is a web application written in Go, named `database-server`. It provides a simple API for creating and listing emails. This application uses a MySQL database for storage and employs `go-chi/chi` for HTTP routing and middleware.

## Features

- REST API endpoints to create and list emails.
- Utilizes the Go `embed` package to include static web content for serving.
- Supports both development and production environments.
- MySQL database integration for data persistence.
- Middleware for common tasks like logging, recovery, and CORS.

## Running the Server

You can start the server by executing the `main.go` in the `cmd` directory:

```
bash
cd cmd/database-server
go build
./database-server
```

By default, the server will start on port 8080 and in production. You can specify a different port using the -door flag:

```
./database-server -door=9090
```

You can also run the server in development mode using the -dev flag:

```
./database-server -dev
```

In development mode, the server expects to find the static web content in the pkg/webserver/www directory and server will load environment variables from the .env file.

## API Endpoints

The server provides the following endpoints:

    POST /create-email
    GET /list-emails

Each of these endpoints are defined and handled in the handlers package. The middlewares package contains necessary middleware for the application such as DBConnected which ensures the database connection is alive before proceeding with the request.

The server also serves static web content at the root path (/), which includes a form for creating new email entries.

## Static Website

The static website portion of this project is powered by [GrapesJS](https://grapesjs.com/demo.html). GrapesJS is an open-source, multi-purpose, Web Builder Framework which combines different tools and features with the goal to help you (or users of your application) to build HTML templates without any knowledge of coding.



## Built With

- [Go](https://golang.org/) - The Go programming language
- [MySQL](https://www.mysql.com/) - MySQL database
- [Chi](https://github.com/go-chi/chi) - Lightweight, idiomatic and composable router for building Go HTTP services
- [CORS](https://github.com/rs/cors) - Go net/http configurable handler to handle CORS requests
- [Godotenv](https://github.com/joho/godotenv) - Go port of Ruby's dotenv library (Loads environment variables from .env)
- [GrapesJS](https://grapesjs.com/) - Open-source, multi-purpose, Web Builder Framework used for building the HTML templates of the static website.


Enjoy using the database-server!


License

This project is licensed under the MIT License.
