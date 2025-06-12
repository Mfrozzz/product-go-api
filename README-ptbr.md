# Product API - Go lang

API RESTful para gerenciamento de produtos e usuários, desenvolvida em Go com Gin, PostgreSQL e arquitetura limpa. Suporta autenticação JWT, controle de permissões por role (admin/user), rate limiting, e está pronta para uso com Docker.

## Sumário 📋
* [Requisitos](#requirements)
* [Configurando Ambiente](#setting-up-the-environment)
* [Middlewares](#middlewares)
* [Endpoints](#endpoints)
* [Scripts](#scripts)
* [Arquitetura do Sistema](#architecture)
* [Estrutura de Pastas](#folder-structure)
* [Versão EN-US](README.md)

---

## <div id="requirements">Requisitos 📄</div>

- Go 1.24.3 ou 1.20+
- Docker
- PostgreSQL
- Insomnia ou Postman

---

## <div id="setting-up-the-environment">Configurando Ambiente ⚙️</div>

### <div>Repositório</div>

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

3. **Instale as dependências Go:**
  ```sh
  go mod tidy
  ```

---

### <div>Docker</div>

* Para subir os containers execute:

```sh
docker compose up -d
```

* Para realizar a build dos containers execute:

```
docker build -t product-go-api .
```

---

### <div>Banco de Dados</div>

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

## <div id="middlewares">Middlewares ↔️</div>

O projeto utiliza três middlewares principais para garantir segurança, controle de acesso e limitação de requisições:

### <div>1. **Auth Middleware**</div>

Responsável por validar o token JWT enviado no header `Authorization`.  
- Só vai permitir que o usuário tenha acesso às rotas caso esteja autenticado (realizado o login).
- Se o token for válido, extrai o campo `role` das claims e armazena no contexto da requisição (`ctx.Set("role", role)`), permitindo que outros middlewares e handlers saibam o papel do usuário autenticado.
- Se o token estiver ausente ou inválido, retorna erro 401 (Unauthorized).

### <div>2. **Rate Limiter Middleware**</div>

Limita o número de requisições por IP para evitar abusos (rate limiting).
- Cada IP pode fazer até 3 requisições por segundo, com um burst máximo de 5.
- Se o limite for excedido, retorna erro 429 (Too Many Requests).

### <div>3. **Require Admin Middleware**</div>

Garante que apenas usuários com papel de admin possam acessar determinadas rotas.
- Verifica o campo `role` no contexto da requisição.
- Se o usuário não for admin, retorna erro 401 (Unauthorized) e bloqueia o acesso à rota.

---

## <div id="endpoints">Endpoints 📌</div>

### <div>Produtos</div>

---

### <div>Usuários</div>

#### POST `/register`

Registra um novo usuário.

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

Autentica um usuário e retorna um token JWT.

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

Obtém as informações de um usuário específico.

- Path Params:
  - `id_user`: O ID do usuário.

- Headers:
  - `Authorization`: Bearer `jwt_token`

- Middlewares Aplicados:
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

Atualiza as informações de um usuário específico.

- Qualquer usuário autenticado pode atualizar seus próprios dados (username, email, password).
- Apenas administradores podem alterar o campo `role` de qualquer usuário.

- Path Params:
  - `id_user`: O ID do usuário.

- Headers:
  - `Authorization`: Bearer `jwt_token`

- Middlewares Aplicados:
  - [Auth Middleware](./middleware/authMiddleware.go)
  - [Require Admin](./middleware/requireAdmin.go)

- Request Body:
  ```json
  {
    "username": "New Name",
    "email": "newemail@example.com",
    "password": "newPassword123",
    "role": "admin" // Só será atualizado se o usuário autenticado for admin
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

- Observações:
  - Se um usuário comum tentar alterar o campo role, receberá erro 403 (Forbidden).
  - Campos não enviados no JSON permanecem inalterados.
  - O campo password sempre será salvo de forma criptografada.

#### DELETE `/api/admin/users/:id_user`

Apenas administradores podem acessar esse endpoint e excluir um usuário do banco de dados.

- Path Params:
  - `id_user`: O ID do usuário.

- Headers:
  - `Authorization`: Bearer `jwt_token`

- Middlewares Aplicados:
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

### <div>Para iniciar</div>

```sh
docker compose up -d
```

ou

```sh
cd product-go-api/cmd
go run main.go
```

---

### <div>Para fazer Build</div>

```sh
docker build -t product-go-api .
```

ou

```sh
cd product-go-api/cmd
go build -o main cmd/main.go
```

---

## <div id="architecture">Arquitetura do sistema 🏛️</div>

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

## <div id="folder-structure">Estrutura de Pastas 📁</div>

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