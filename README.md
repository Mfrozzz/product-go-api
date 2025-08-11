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
* [Useful Links](#useful-links)
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

### <div>Tip: How to create an admin user ✉️</div>

To turn a user into an admin directly in the database, run:

```sh
UPDATE users SET role = 'admin' WHERE id = 1;
```

To turn a user into an super_admin directly in the database, run:

```sh
UPDATE users SET role = 'super_admin' WHERE id = 1;
```

---

## <div id="middlewares">Middlewares ↔️</div>

The project uses three main middlewares to ensure security, access control, and request limiting:

### <div id="auth-middleware">1. **Auth Middleware**</div>

Responsible for validating the JWT token sent in the `Authorization` header.
- Only allows the user to access routes if authenticated (logged in).
- If the token is valid, extracts the `role` field from the claims and stores it in the request context (`ctx.Set("role", role)`), allowing other middlewares and handlers to know the authenticated user's role.
- If the token is missing or invalid, returns a 401 (Unauthorized) error.

### <div id="rate-limiter">2. **Rate Limiter Middleware**</div>

Limits the number of requests per IP to prevent abuse (rate limiting).
- Each IP can make up to 3 requests per second, with a maximum burst of 5.
- If the limit is exceeded, returns a 429 (Too Many Requests) error.

### <div id="require-admin">3. **Require Admin Middleware**</div>

Ensures that only users with the admin role can access certain routes.

- Checks the `role` field in the request context.
- If the user is not an admin, returns a 401 (Unauthorized) error and blocks access to the route.

---

## <div id="endpoints">Endpoints 📌</div>

### <div>Products</div>

#### POST `/api/products`

Registers a new product in the system.

- Headers:
  - `Authorization`: Bearer `jwt_token`

- Applied Middlewares:
  - [Auth Middleware](#auth-middleware)

- Request Body:
  ```json
  {
    "name": "Potato",
    "price": 4.45
  }
  ```

- Response:
  ```json
  {
    "id_product": 1,
    "name": "Potato",
    "price": 4.45
  }
  ```

#### GET `/api/products`

Lists all products from the database. You can use parameters, filters, and pagination in the results.

- **Query Parameters**:
  - `page` (optional): Page number, default = 1
  - `limit` (optional): Number of items per page, default = 10
  - `name` (optional): Filter by product name

- **Examples**:
  - List all products (default):
  ```
  GET /api/products
  ```
  - List products with pagination, limit, and name filter:
  ```
  GET /api/products?page=1&limit=5&name=Potato
  ```

- Headers:
  - `Authorization`: Bearer `jwt_token`

- Applied Middlewares:
  - [Auth Middleware](#auth-middleware)

- Response:
  ```json
  [
    {
      "id_product": 1,
      "name": "Potato",
      "price": 4.45
    },
    {
      "id_product": 9,
      "name": "Potato Chips",
      "price": 9
    }
    ...
  ]
  ```

#### GET `/api/products/:id_product`

Retrieves information about a specific product.

- Path Params:
  - `id_product`: The product ID.

- Headers:
  - `Authorization`: Bearer `jwt_token`

- Applied Middlewares:
  - [Auth Middleware](#auth-middleware)

- Response:
  ```json
  {
    "id_product": 1,
    "name": "Potato",
    "price": 4.45
  }
  ```

#### PUT `/api/products/:id_product`

Updates information for a specific product.

- Path Params:
  - `id_product`: The product ID.

- Headers:
  - `Authorization`: Bearer `jwt_token`

- Applied Middlewares:
  - [Auth Middleware](#auth-middleware)

- Object before update:
  ```json
  {
    "id_product": 14,
    "name": "Pasta",
    "price": 10.2
  }
  ```

- Request Body:
  ```json
  {
    "name": "Spaghetti Pasta",
    "price": 13.2
  }
  ```

- Response:
  ```json
  {
    "id_product": 14,
    "name": "Spaghetti Pasta",
    "price": 13.2
  }
  ```

- Notes:
  - Fields not sent in the JSON remain unchanged.


#### DELETE `/api/admin/products/:id_product`

Only administrators can access this endpoint and delete a product from the database.

- Path Params:
  - `id_product`: The product ID

- Headers:
  - `Authorization`: Bearer `jwt_token`

- Applied Middlewares:
  - [Auth Middleware](#auth-middleware)
  - [Require Admin](#require-admin)

- Response:
  ```json
  {
    "Message": "Product deleted successfully"
  }
  ```

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
  - [Auth Middleware](#auth-middleware)

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

#### GET `/api/user/info`

Returns the information of the currently authenticated user (based on the JWT token).

- Headers:
  - `Authorization`: Bearer `jwt_token`

- Applied Middlewares:
  - [Auth Middleware](#auth-middleware)

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

#### GET `/api/admin/users`

Lists all users from the database. You can use parameters, filters, and pagination in the results.

- **Query Parameters**:
  - `page` (optional): Page number, default = 1
  - `limit` (optional): Number of items per page, default = 10
  - `name` (optional): Filter by username

- **Examples**:
  - List all users (default):
  ```
  GET /api/admin/users
  ```
  - List users with pagination, limit, and name filter:
  ```
  GET /api/admin/users?page=1&limit=5&name=user
  ```

- Headers:
  - `Authorization`: Bearer `jwt_token`

- Applied Middlewares:
  - [Auth Middleware](#auth-middleware)
  - [Require Admin](#require-admin)

- Response:
  ```json
  [
    {
      "id_user": 1,
      "username":"Test Example",
      "email": "user@example.com",
      "password": "your_encrypted_password",
      "role": "user"
    },
    {
      "id_user": 2,
      "username":"Test Example 2",
      "email": "user2@example.com",
      "password": "your_encrypted_password",
      "role": "user2"
    }
    ...
  ]
  ```

#### PUT `/api/users/:id_user`

Updates information for a specific user.

- Any authenticated user can update their own data (username, email, password).
- Only **super_admins** are allowed to:
  - Promote a user to `super_admin`
  - Modify the role of another `admin` or `super_admin`

- Path Params:
  - `id_user`: The user ID.

- Headers:
  - `Authorization`: Bearer `jwt_token`

- Applied Middlewares:
  - [Auth Middleware](#auth-middleware)
  - [Require Admin](#require-admin)

- Object before update:
```json
  {
    "id_user": 1,
    "username": "Name",
    "email": "email@example.com",
    "password": "Password123",
    "role": "user"
  }
```

- Request Body:
  ```json
  {
    "username": "New Name",
    "email": "newemail@example.com",
    "password": "newPassword123",
    "role": "admin" // Will only be updated if the authenticated user is admin or super_admin
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
  - If a user without the required permissions tries to change the role field, a 403 (Forbidden) error will be returned.
  - Fields not sent in the JSON remain unchanged.
  - The password field is always saved in encrypted form.
  - Only `super_admins` can:
    - Promote users to `super_admin`
    - Change the role of users with role `admin` or `super_admin`

#### DELETE `/api/admin/users/:id_user`

Only administrators can access this endpoint and delete a user from the database.

- Path Params:
  - `id_user`: The user ID.

- Headers:
  - `Authorization`: Bearer `jwt_token`

- Applied Middlewares:
  - [Auth Middleware](#auth-middleware)
  - [Require Admin](#require-admin)

- Response:
  ```json
  {
    "Message": "User deleted successfully"
  }
  ```

- Notes:
  - Only **super_admins** are allowed to delete users with the role `admin` or `super_admin`.
  - If an `admin` attempts to delete another `admin` or a `super_admin`, a `403 Forbidden` error will be returned.

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

## <div id="useful-links">Useful Links 🔗</div>

- [Go (Golang) Official Documentation](https://golang.org/doc/)
- [Go Modules Reference](https://blog.golang.org/using-go-modules)
- [Gin Web Framework Documentation](https://gin-gonic.com/en/docs/)
- [Docker Documentation](https://docs.docker.com/)
- [PostgreSQL Official Documentation](https://www.postgresql.org/docs/)
- [JWT Introduction](https://jwt.io/introduction)
- [Insomnia](https://insomnia.rest/) / [Postman](https://www.postman.com/) — API testing tools

---

Developed by [Marcos Vinicius Boava](https://github.com/Mfrozzz), using as a base [Go Lab Tutoriais](https://www.youtube.com/watch?v=3p4mpId_ZU8&t=3317s).