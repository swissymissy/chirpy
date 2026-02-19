# Chirpy: A Social Media Backend in Go
Chirpy is a RESTful API and backend server built with Go. It allows users to create accounts, authenticate, and share short text posts ("chirps"). This project focuses on building robust HTTP routing, database integration, and secure user management.

## Features
- User authentication
- Chirp creation and retrieval
- Profanity filtering

## Prerequisites
- Go 1.22+
- A running PostgreSQL instance

## ⚙️ Installation
1. **Clone the repository**
`git clone https://github.com/swissymissy/chirpy.git`
2. **Configure environment variables**
Create a new `.env` file and add it to `.gitignore`. The `.env` file has these  important variables:
```
PORT=8080
DB_URL=postgres://user:password@localhost:5432/chirpy?sslmode=disable
JWT_SECRET=your_super_secret_key_here
PLATFORM="dev"
```
- PORT: Default port is set to 8080, but you can change to the port you want
- DB_URL: The connection string to the database. In this project we use Postgres
- JWT_SECRET: The secret key for your server to generate JWT token for user
- PLATFORM: Set to "dev" to enable development-only endpoints. If omitted, the server defaults to a secure state where administrative resets are disabled
## Usage
**Compile and run**
```
go build -o chirpy
./chirpy
```
- The server will be host on port on 8080 if you keep the default port: `http://localhost:8080`
- Homepage: `http://localhost:8080/app/`
## API Endpoints
Some API endpoints that the server can handle:
1. **GET**
- `GET /api/healthz` : Check server's health
- `GET /admin/metrics` : Show how many times Chirpy has been visited
- `GET /api/chirps` : Get list of chirps in ascending order. It also accepts query parameters such as `author_id` to get chirps of a specific user, or `sort` to reorder the list.
- GET /api/chirps/{chirpID}` : Get information of a chirp by its ID
2.**POST**
- `POST /api/users` : Create new user
- `POST /api/chirps` : Create new chirp
- `POST /api/login` : Authenticate user logging in
- `POST /api/refresh` : Check user's refresh token and access token
- `POST /api/revoke` : Update "revoked" status of user's refresh token
- `POST /api/polka/webhooks` : handle communication with third-party server
3.**PUT**
- `PUT /api/users` : Let authorized user change their password and email
4.**DELETE**
- `DELETE /api/chirps/{chirpID}` : Let authorized user delete their posted chirp
### Safety Feature
This API endpoint is built with a gatekeeper PLATFORM="dev" in `.env` because it is a reset database endpoint. Treat carefully. Returns `403 Forbidden` if PLATFORM is not set to `"dev"`.
- `POST /admin/reset` : Reset users table