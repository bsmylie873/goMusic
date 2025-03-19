## GoMusic API
A RESTful API for managing music information including albums, artists, and bands with JWT-based authentication.

### Features
JWT Authentication: Secure endpoints with JSON Web Tokens
User Management: Registration, login, and profile management
Music Catalog: CRUD operations for albums, artists, and bands
SQLite Database: Simple data persistence
Test Coverage: Comprehensive test suites for all services

### Setup Instructions
#### 1. Clone the repository
`git clone https://github.com/yourusername/goMusic.git`
`cd goMusic`

#### 2. Set environment variables
`export JWT_SECRET_KEY=your-secret-key-here`

#### 3.Run the application
`go run main.go`

The server will start on localhost:8082

### API Endpoints
#### Authentication
* POST /register - Register a new user
```
{
"username": "user",
"password": "password123",
"email": "user@example.com"
}
```

* POST /login - Login and get JWT token
```
{
"username": "user",
"password": "password123"
}
```

* GET /profile - Get authenticated user profile (protected)

#### Albums
* GET /albums - Get all albums
* GET /albums/{id} - Get album by ID
* POST /albums - Create new album (protected)
* PUT /albums/{id} - Update album (protected)
* DELETE /albums/{id} - Delete album (protected)

#### Artists
* GET /artists - Get all artists
* GET /artists/{id} - Get artist by ID
* POST /artists - Create new artist (protected)
* PUT /artists/{id} - Update artist (protected)
* DELETE /artists/{id} - Delete artist (protected)

#### Bands
* GET /bands - Get all bands
* GET /bands/{id} - Get band by ID
* POST /bands - Create new band (protected)
* PUT /bands/{id} - Update band (protected)
* DELETE /bands/{id} - Delete band (protected)

### Authentication
The API uses JWT (JSON Web Tokens) for authentication. To access protected endpoints:
1. Register or login to get a token
2. Add the token to request headers:`Authorization: Bearer your-jwt-token`

### Testing
Run the test suite with 

`go test ./...`

The project includes comprehensive tests for all services, including authentication, user management, and CRUD operations.

### Dependencies
* Go 1.21+
* github.com/golang-jwt/jwt/v4 - JWT implementation
* github.com/DATA-DOG/go-sqlmock - Database mocking for tests
* github.com/stretchr/testify - Test assertions

### License
MIT