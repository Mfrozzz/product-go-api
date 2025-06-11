# Product API - Go lang

RESTful API for product and user management, developed in Go with Gin, PostgreSQL, and clean architecture. Supports JWT authentication, role-based permissions (admin/user), rate limiting, and is ready for use with Docker.

## Table of Contents 📋
* [Requirements](#requirements)
* [Setting Up the Environment](#setting-up-the-environment)
* [Endpoints](#endpoints)
* [Scripts](#scripts)
* [System Architecture](#architecture)
* [Folder Structure](#folder-structure)
* [PT-BR version](README-ptbr.md)

---

### <div id="requirements">Requirements 📄</div>

- Go 1.24.3 or 1.20+
- Docker
- PostgreSQL
- Insomnia or Postman

---

### <div id="setting-up-the-environment">Setting Up the Environment ⚙️</div>

#### <div>Repository</div>

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

---

#### <div>Docker</div>

* To start the containers, run:

```sh
docker compose up -d
```

* To build the containers, run:

```sh
docker build -t product-go-api .
```

---

#### <div>Database</div>

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

### <div id="endpoints">Endpoints 📌</div>

#### <div>Products</div>

---

#### <div>Users</div>

---

### <div id="scripts">Scripts ⌨️</div>

#### <div>To Run</div>

```sh
docker compose up -d
```

or

```sh
cd product-go-api/cmd
go run main.go
```

---

#### <div>To Build</div>

```sh
docker build -t product-go-api .
```

or

```sh
cd product-go-api/cmd
go build -o main cmd/main.go
```

---

### <div id="architecture">System Architecture 🏛️</div>

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

### <div id="folder-structure">Folder Structure 📁</div>

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