# Go Auth

A simple authentication API built with Go and Gin.  
It features user registration, login, JWT-based authentication, and protected routes.

## Features

- User Registration
- User Login
- Protected profile and logout routes
- Password hashing with bcrypt
- JWT Authentication
- Unit tests and repository mocking for controller logic

## Getting Started

### Prerequisites

- Go 1.18+ installed (https://golang.org/dl/)
- (Optional) Docker for containerization

### Installation

Clone the repository:

```sh
git clone https://github.com/your-username/go-auth-app.git
cd go-auth-app
```

Download dependencies:

```sh
go mod tidy
```

### Running the App

```sh
go run main.go
```

By default, the server runs on `localhost:8080`.

### API Endpoints

- `POST /register` — Register a new user  
  **Body:**  
  ```json
  {
    "name": "Your Name",
    "email": "email@example.com",
    "password": "yourpassword"
  }
  ```

- `POST /login` — Login a user  
  **Body:**  
  ```json
  {
    "email": "email@example.com",
    "password": "yourpassword"
  }
  ```

- `GET /profile` — Get the current user profile (JWT required in headers)

- `POST /logout` — Logout (implement according to your JWT invalidation strategy)

### Running Tests

Run all tests with:

```sh
go test ./tests/...
```

## Project Structure

```
.
├── controllers/        # HTTP handler logic
├── models/             # Data models
├── repository/         # Data access layer
├── tests/              # Unit tests and mocks
└── main.go             # App entry point
```

## Contributing

1. Fork the repo
2. Create your feature branch (`git checkout -b feature/your-feature`)
3. Commit your changes (`git commit -am 'feat: add new feature'`)
4. Push to the branch (`git push origin feature/your-feature`)
5. Open a pull request

## License

MIT



