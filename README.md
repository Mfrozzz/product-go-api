# Product API - Go lang

RESTful API for product and user management, developed in Go with Gin, PostgreSQL, and clean architecture. Supports JWT authentication, role-based permissions (admin/user), rate limiting, and is ready for use with Docker.

## Table of Contents 📋
* [Requirements](#requirements)
* [Setting Up the Environment](#setting-up-the-environment)
* [Middlewares](#middlewares)
* [Endpoints](#endpoints)
* [Scripts](#scripts)
* [System Architecture](#architecture)
* [Folder Structure](#folder-structure)
* [PT-BR version](README-ptbr.md)

---

## <div id="requirements">Requirements 📄</div>

- Go 1.24.3 or 1.20+
- Docker
- PostgreSQL
- Insomnia or Postman

---

## <div id="setting-up-the-environment">Setting Up the Environment ⚙️</div>

### <div>Repository</div>

1. **Clone the repository:**

   ```sh
   git clone https://github.com/Mfrozzz/product-go-api.git
   cd product-go-api
   ```

2. **Configure the `.env` file:**

    * Copy the `.env.example` file to `.env` and fill in your settings:

    ```sh
    cp .env.example .env
    ```

    * File example:
    
    ```.env
    JWT_SECRET="YOUR-SECRET-KEY"
    PORT=":8000"

    DB_HOST="go_db"
    DB_PORT=5432
    DB_USER="YOUR-DATABASE-USER"
    DB_PASSWORD="YOUR-DATABASE-PASSWORD"
    DB_NAME="YOUR-DATABASE-NAME"
    ```

    * **Important:** The project depends on the `.env` variables to connect to the database and generate JWT tokens.

3. **Install Go dependencies:**
  ```sh
  go mod tidy
  ```

---

### <div>Docker</div>

* To start the containers, run:

```sh
docker compose up -d
```

* To build the containers, run:

```sh
docker build -t product-go-api .
```

---

### <div>Database</div>

The choice of tool to manage the database is up to the developer, such as DBeaver or any other software of your preference. However, in this case, access to the database will be demonstrated via command line, accessing directly the Docker container where the database is running.

* Access the database:

```sh
docker exec -it go_db bash
psql -d postgres -U postgres
```

* Create tables in the database:

```sh
CREATE TABLE product (
  id SERIAL PRIMARY KEY,
  product_name VARCHAR(50) NOT NULL,
  price NUMERIC(10, 2) NOT NULL
);
CREATE TABLE users (
  id SERIAL PRIMARY KEY,
  email VARCHAR(255) UNIQUE NOT NULL,
  username VARCHAR(255) NOT NULL,
  password VARCHAR(255) NOT NULL,
  role VARCHAR(20) NOT NULL DEFAULT 'user'
);
```

---

## <div id="middlewares">Middlewares ↔️</div>

The project uses three main middlewares to ensure security, access control, and request limiting:

### <div>1. **Auth Middleware**</div>

Responsible for validating the JWT token sent in the `Authorization` header.
- Only allows the user to access routes if authenticated (logged in).
- If the token is valid, extracts the `role` field from the claims and stores it in the request context (`ctx.Set("role", role)`), allowing other middlewares and handlers to know the authenticated user's role.
- If the token is missing or invalid, returns a 401 (Unauthorized) error.

### <div>2. **Rate Limiter Middleware**</div>

Limits the number of requests per IP to prevent abuse (rate limiting).
- Each IP can make up to 3 requests per second, with a maximum burst of 5.
- If the limit is exceeded, returns a 429 (Too Many Requests) error.

### <div>3. **Require Admin Middleware**</div>

Ensures that only users with the admin role can access certain routes.

- Checks the `role` field in the request context.
- If the user is not an admin, returns a 401 (Unauthorized) error and blocks access to the route.

---

## <div id="endpoints">Endpoints 📌</div>

### <div>Products</div>

---

### <div>Users</div>

#### POST `/register`

Registers a new user.

- Request Body:
  ```json
  {
    "username":"Test Example",
    "email": "user@example.com",
    "password": "password123"
  }
  ```

- Response:
  ```json
  {
    "Message": "Register successful",
    "User": {
      "email": "user@example.com",
      "role": "user",
      "user_id": 1,
      "username":"Test Example"
    }
  }
  ```

#### POST `/login`

Authenticates a user and returns a JWT token.

- Request Body:
  ```json
  {
    "email": "user@example.com",
    "password": "password123"
  }
  ```

- Response:
  ```json
  {
    "Message": "Login successful",
    "token": "your_jwt_token"
  }
  ```

#### GET `/api/users/:id_user`

Retrieves information about a specific user.

- Path Params:
  - `id_user`: The user ID.

- Headers:
  - `Authorization`: Bearer `jwt_token`

- Applied Middlewares::
  - [Auth Middleware](./middleware/authMiddleware.go)

- Response:
  ```json
  {
    "id_user": 1,
    "username":"Test Example",
    "email": "user@example.com",
    "password": "your_encrypted_password",
    "role": "user"
  }
  ```

#### PUT `/api/users/:id_user`

Updates information for a specific user.

- Any authenticated user can update their own data (username, email, password).
- Only administrators can change the `role` field of any user.

- Path Params:
  - `id_user`: The user ID.

- Headers:
  - `Authorization`: Bearer `jwt_token`

- Applied Middlewares:
  - [Auth Middleware](./middleware/authMiddleware.go)
  - [Require Admin](./middleware/requireAdmin.go)

- Request Body:
  ```json
  {
    "username": "New Name",
    "email": "newemail@example.com",
    "password": "newPassword123",
    "role": "admin" // Will only be updated if the authenticated user is admin
  }
  ```

- Response:
  ```json
  {
    "id_user": 1,
    "username":"New Name",
    "email": "newemail@example.com",
    "password": "your_new_encrypted_password",
    "role": "admin"
  }
  ```

- Notes:
  - If a regular user tries to change the role field, a 403 (Forbidden) error will be returned.
  - Fields not sent in the JSON remain unchanged.
  - The password field is always saved in encrypted form.

#### DELETE `/api/admin/users/:id_user`

Only administrators can access this endpoint and delete a user from the database.

- Path Params:
  - `id_user`: The user ID.

- Headers:
  - `Authorization`: Bearer `jwt_token`

- Applied Middlewares:
  - [Auth Middleware](./middleware/authMiddleware.go)
  - [Require Admin](./middleware/requireAdmin.go)

- Response:
  ```json
  {
    "Message": "User deleted successfully"
  }
  ```

---

## <div id="scripts">Scripts ⌨️</div>

### <div>To Run</div>

```sh
docker compose up -d
```

or

```sh
cd product-go-api/cmd
go run main.go
```

---

### <div>To Build</div>

```sh
docker build -t product-go-api .
```

or

```sh
cd product-go-api/cmd
go build -o main cmd/main.go
```

---

## <div id="architecture">System Architecture 🏛️</div>

The system architecture follows the Clean Architecture pattern, clearly separating responsibilities into layers. This makes the project easier to maintain, test, and evolve.

```
               ┌───────────────────────────────┐
               │ Controllers (Gin HTTP handler)│
               └───────────────┬───────────────┘
                               │
               ┌───────────────▼───────────────┐
               │   Usecase (Business Rules)    │
               └───────────────┬───────────────┘
                               │
               ┌───────────────▼───────────────┐
               │ Repository (Data Access)      │
               └───────────────┬───────────────┘
                               │
               ┌───────────────▼───────────────┐
               │      Model (Entities)         │
               └───────────────┬───────────────┘
                               │
               ┌───────────────▼───────────────┐
               │   Database (PostgreSQL)       │
               └───────────────────────────────┘
```

---

## <div id="folder-structure">Folder Structure 📁</div>

```
product-go-api/
├── cmd/
|   └── main.go
├── controller/
|   ├── product_controller.go
|   └── user_controller.go
├── db/
|   └── connection.go
├── middleware
|   ├── authMiddleware.go
|   ├── rateLimiter.go
|   └── requireAdmin.go
├── model/
|   ├── product.go
|   ├── response.go
|   └── user.go
├── repository/
|   ├── product_repository.go
|   └── user_repository.go
├── usecase/
|   ├── product_usecase.go
|   └── user_usecase.go
├── .env
├── .env.example
├── .gitignore
├── docker-compose.yml
├── Dockerfile
├── go.mod
├── go.sum
├── README-ptbr.md
└── README.md
```

---