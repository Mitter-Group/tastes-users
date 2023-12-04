# My Fiber App

This is a simple microservice built with Golang and Fiber, following the principles of clean architecture.

## Project Structure

The project is structured as follows:

- `cmd/app/main.go`: This is the entry point of the application. It sets up the router and starts the server.
- `internal/app/handler/login.go`: This file exports a function `Login` which handles the login endpoint. It uses the `LoginService` to perform the login operation.
- `internal/app/router/router.go`: This file exports a function `SetupRoutes` which sets up the routes for the application. It uses the `Login` handler for the login endpoint.
- `internal/app/middleware/middleware.go`: This file exports middleware functions that can be used in the application.
- `internal/domain/login/service.go`: This file exports a `LoginService` interface and a `LoginServiceImpl` struct that implements the interface. The `LoginService` interface defines a `Login` method.
- `internal/domain/login/repository.go`: This file exports a `LoginRepository` interface and a `LoginRepositoryImpl` struct that implements the interface. The `LoginRepository` interface defines a `Login` method.
- `internal/infrastructure/http/client.go`: This file exports a `Client` struct and a `NewClient` function. The `Client` struct has a `Get` method which makes a GET request to a given URL.
- `pkg/config/config.go`: This file exports a `Config` struct and a `NewConfig` function. The `Config` struct holds the configuration of the application.
- `go.mod` and `go.sum`: These files are used to manage the dependencies of the application.

## Running the Application

To run the application, navigate to the `cmd/app` directory and run the following command:

```bash
go run main.go
```

This will start the server on the port specified in the configuration.

## Endpoints

The application has the following endpoint:

- `POST /login`: This endpoint accepts a JSON body with a `provider` field. Depending on the value of the `provider` field, it will call another microservice via a GET request to `/login`.

## Configuration

The configuration of the application can be changed by modifying the `Config` struct in `pkg/config/config.go`.

## Dependencies

The application uses the following dependencies:

- [Fiber](https://github.com/gofiber/fiber): A web framework for Go.
- [Viper](https://github.com/spf13/viper): A library for managing configuration in Go.

## Contributing

Contributions are welcome. Please submit a pull request or create an issue to discuss the changes you want to make.