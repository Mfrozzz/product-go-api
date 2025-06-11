# Product API - Go lang

API RESTful para gerenciamento de produtos e usuários, desenvolvida em Go com Gin, PostgreSQL e arquitetura limpa. Suporta autenticação JWT, controle de permissões por role (admin/user), rate limiting, e está pronta para uso com Docker.

## Sumário 📋
* [Requisitos](#requirements)
* [Configurando Ambiente](#setting-up-the-environment)
* [Endpoints](#endpoints)
* [Scripts](#scripts)
* [Arquitetura do Sistema](#architecture)
* [Estrutura de Pastas](#folder-structure)
* [Versão EN-US](README.md)

---

### <div id="requirements">Requisitos 📄</div>

- Go 1.24.3 ou 1.20+
- Docker
- PostgreSQL
- Insomnia ou Postman

---

### <div id="setting-up-the-environment">Configurando Ambiente ⚙️</div>

#### <div>Repositório</div>

1. **Clone o repositório:**
   ```sh
   git clone https://github.com/Mfrozzz/product-go-api.git
   cd product-go-api
   ```

2. **Consigure o arquivo `.env`:**
    * Copie o arquivo `.env.example` para `.env` e preencha com as suas configurações:
    ```sh
    cp .env.example .env
    ```

    * Exemplo do arquivo:
    ```.env
    JWT_SECRET="YOUR-SECRET-KEY"
    PORT=":8000"

    DB_HOST="go_db"
    DB_PORT=5432
    DB_USER="YOUR-DATABASE-USER"
    DB_PASSWORD="YOUR-DATABASE-PASSWORD"
    DB_NAME="YOUR-DATABASE-NAME"
    ```

    * **Importante:** O projeto depende das variáveis do `.env` para conectar ao banco e gerar tokens JWT.

#### <div>Docker</div>
* Para subir os containers execute:

```sh
docker compose up -d
```

* Para realizar a build dos containers execute:

```
docker build -t product-go-api .
```

#### <div>Banco de Dados</div>

A escolha da ferramenta para gerenciar o banco de dados fica a critério do desenvolvedor, podendo ser utilizado, por exemplo, o DBeaver ou outro software de sua preferência. No entanto, neste caso, o acesso ao banco será demonstrado via linha de comando, acessando diretamente o container Docker onde o banco de dados está em execução.

* Acesso ao Banco de dados:

```sh
docker exec -it go_db bash
psql -d postgres -U postgres
```

* Criar Tabelas no Banco de dados:
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

#### <div>Produtos</div>

#### <div>Usuários</div>

---

### <div id="scripts">Scripts ⌨️</div>

#### <div>Para iniciar</div>

```sh
docker compose up -d
```

ou

```sh
cd product-go-api/cmd
go run main.go
```

#### <div>Para fazer Build</div>
```sh
docker build -t product-go-api .
```

ou

```sh
cd product-go-api/cmd
go build -o main cmd/main.go
```

---

### <div id="architecture">Arquitetura do sistema 🏛️</div>

A arquitetura do sistema segue o padrão Clean Architecture, separando claramente as responsabilidades em camadas. Isso facilita a manutenção, testes e evolução do projeto.

```
           ┌───────────────────────────────┐
           │ Controllers (Gin HTTP handler)│
           └───────────────┬───────────────┘
                           │
           ┌───────────────▼───────────────┐
           │   Usecase (Regras de negócio) │
           └───────────────┬───────────────┘
                           │
           ┌───────────────▼───────────────┐
           │ Repository (Acesso a dados)   │
           └───────────────┬───────────────┘
                           │
           ┌───────────────▼───────────────┐
           │      Model (Entidades)        │
           └───────────────┬───────────────┘
                           │
           ┌───────────────▼───────────────┐
           │   Database (PostgreSQL)       │
           └───────────────────────────────┘
```

---

### <div id="folder-structure">Estrutura de Pastas 📁</div>

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