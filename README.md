## Ticketon Users API -  A Microservice for User Management (Go & Docker)

**Introduction**

This document outlines the Ticketon Users API, a microservice for user management, JWT authentication, and secure routes. 

**Features**

* User Registration
* User Login with JWT Token Generation
* JWT Token Validation for Secure Routes
* Built with Go and MySQL
* Dockerized for Easy Deployment

**Prerequisites**

Before running the project, ensure you have the following installed:

* Docker
* Docker Compose
* Go (for local development)

**Getting Started**

1. **Clone the Repository**

```bash
git clone [https://github.com/yourusername/ticketon-users-api.git](https://github.com/yourusername/ticketon-users-api.git)
cd ticketon-users-api
```
2. **Environment Setup**
Create a .env file with the following environment variables:

```bash
Copy code
DB_HOST=ticketon-users-db
DB_DRIVER=mysql
DB_USER=dbuser
DB_PASSWORD=dbpasswd
DB_NAME=ticketon_users_db
DB_PORT=3306
JWT_SECRET=your-secret-key
```
Make sure to replace your-secret-key with a strong secret key for JWT token signing.

3. **Docker Compose**
To build and start the services, run the following command:

```bash
Copy code
docker-compose up --build
```
This will spin up both the users-api and MySQL services. The API will be available at http://localhost:8080.

4. **Database Migration**
The MySQL container will automatically initialize the database on startup. You can customize the database schema in your Go application if necessary.

**API Endpoints**
1. Register User
Endpoint: ```POST /api/users```

Description: Registers a new user by saving their details to the database.

Request Body:
```json
{
  "firstname": "Joey",
  "lastname": "Ramone",
  "dni": 12345678,
  "email": "joeyramone@gmail.com",
  "password": "pass",
  "phone": "54388134920"
}
```
Response:
```json
{
  "account_id": 2,
  "email": "joeyramone@gmail.com",
  "user_id": 0
}
```
**2. Login User**
Endpoint: ```POST /api/v1/login```

Description: Authenticates the user and returns a JWT token.

Request Body:
```json
{
  "username": "johndoe",
  "password": "password123"
}
```
Response:
```json
{
  "token": "your-jwt-token-here"
}
```
**3. Protected Route Example**
Endpoint: ```GET /api/v1/profile```

Description: Retrieves the user's profile. This route is protected and requires a valid JWT token.

*Headers*:
```http
Authorization: Bearer your-jwt-token-here
```
*Response*:
```json
{
  "user_id": 1,
  "username": "johndoe"
}
```
**JWT Token**
The API uses JWT tokens for user authentication. After a successful login, the API returns a token, which must be sent with every request to protected routes via the Authorization header:

```http
Authorization: Bearer your-jwt-token-here
```
*Token Validation*
To validate a JWT token, the API decodes it using the JWT_SECRET and verifies its integrity. If the token is valid, access to the protected route is granted.

**Running Tests**
If you've included unit tests or integration tests, you can run them using:

```bash
go test ./...
```

**Project Structure**
```bash
.
├── Dockerfile
├── docker-compose.yml
├── .env
├── src
│   ├── main.go            # Entry point for the API
│   ├── handlers           # Request handlers for routes
│   ├── models             # Database models
│   ├── routes             # API routes
│   ├── middlewares        # JWT and other middleware
│   └── config             # Database and app configuration
└── README.md
```

**License**
This project is licensed under the MIT License - see the LICENSE file for details.
