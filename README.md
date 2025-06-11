# Product API - Go lang

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


---

### <div id="setting-up-the-environment">Setting Up the Environment ⚙️</div>

#### <div>Repository</div>

#### <div>Docker</div>

#### <div>Database</div>

---

### <div id="endpoints">Endpoints 📌</div>

#### <div>Products</div>

#### <div>Users</div>

---

### <div id="scripts">Scripts ⌨️</div>

#### <div>To Run</div>

#### <div>To Build</div>

---

### <div id="architecture">System Architecture 🏛️</div>

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