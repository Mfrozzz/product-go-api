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
* [Links Úteis](#useful-links)
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

### <div>Dica: Como criar um usuário admin ✉️</div>

Para transformar um usuário em admin diretamente pelo banco de dados, execute:

```sh
UPDATE users SET role = 'admin' WHERE id = 1;
```

---

## <div id="middlewares">Middlewares ↔️</div>

O projeto utiliza três middlewares principais para garantir segurança, controle de acesso e limitação de requisições:

### <div id="auth-middleware">1. **Auth Middleware**</div>

Responsável por validar o token JWT enviado no header `Authorization`.  
- Só vai permitir que o usuário tenha acesso às rotas caso esteja autenticado (realizado o login).
- Se o token for válido, extrai o campo `role` das claims e armazena no contexto da requisição (`ctx.Set("role", role)`), permitindo que outros middlewares e handlers saibam o papel do usuário autenticado.
- Se o token estiver ausente ou inválido, retorna erro 401 (Unauthorized).

### <div id="rate-limiter">2. **Rate Limiter Middleware**</div>

Limita o número de requisições por IP para evitar abusos (rate limiting).
- Cada IP pode fazer até 3 requisições por segundo, com um burst máximo de 5.
- Se o limite for excedido, retorna erro 429 (Too Many Requests).

### <div id="require-admin">3. **Require Admin Middleware**</div>

Garante que apenas usuários com papel de admin possam acessar determinadas rotas.
- Verifica o campo `role` no contexto da requisição.
- Se o usuário não for admin, retorna erro 401 (Unauthorized) e bloqueia o acesso à rota.

---

## <div id="endpoints">Endpoints 📌</div>

### <div>Produtos</div>

#### POST `/api/products`

Registra no sistema um novo produto.

- Headers:
  - `Authorization`: Bearer `jwt_token`

- Middlewares Aplicados:
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

Faz a listagem de todos os produtos do banco de dados.
Pode ser usado parâmetros, filtros e paginação nos resultados

- **Parâmetros de Busca**:
  - `page` (opcional): Número da página, valor padrão = 1
  - `limit` (opcional): Número de itens por página, valor padrão = 10
  - `name` (opcional): Filtro pelo nome do produto

- **Exemplos**:
  - Listar todos os produtos (padrão):
  ```
  GET /api/products
  ```
  - Listar produtos com paginação, limite e filtro por nome:
  ```
  GET /api/products?page=1&limit=5&name=Potato
  ```

- Headers:
  - `Authorization`: Bearer `jwt_token`

- Middlewares Aplicados:
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

Obtém as informações de um produto específico.

- Path Params:
  - `id_product`: O ID do produto.

- Headers:
  - `Authorization`: Bearer `jwt_token`

- Middlewares Aplicados:
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

Atualiza as informações de um produto específico.

- Path Params:
  - `id_product`: O ID do produto.

- Headers:
  - `Authorization`: Bearer `jwt_token`

- Middlewares Aplicados:
  - [Auth Middleware](#auth-middleware)

- Objeto antes da atualização:
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

- Observações:
  - Campos não enviados no JSON permanecem inalterados.


#### DELETE `/api/admin/products/:id_product`

Apenas administradores podem acessar esse endpoint e excluir um produto do banco de dados.

- Path Params:
  - `id_product`: O ID do produto.

- Headers:
  - `Authorization`: Bearer `jwt_token`

- Middlewares Aplicados:
  - [Auth Middleware](#auth-middleware)
  - [Require Admin](#require-admin)

- Response:
  ```json
  {
    "Message": "Product deleted successfully"
  }
  ```

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

Retorna as informações do usuário autenticado (baseado no Token JWT).

- Headers:
  - `Authorization`: Bearer `jwt_token`

- Middlewares Aplicados:
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

Retorna todos os usuários do banco de dados. Pode ser usado parâmetros, filtros e paginação nos resultados

- **Parâmetros de Busca**:
  - `page` (opcional): Número da página, valor padrão = 1
  - `limit` (opcional): Número de itens por página, valor padrão = 10
  - `name` (opcional): Filtro pelo nome do usuário

- **Exemplos**:
  - Listar todos os usuários (padrão):
  ```
  GET /api/admin/users
  ```
  - Listar usuários com paginação, limite e filtro por nome:
  ```
  GET /api/admin/users?page=1&limit=5&name=user
  ```

- Headers:
  - `Authorization`: Bearer `jwt_token`

- Middlewares Aplicados:
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

Atualiza as informações de um usuário específico.

- Qualquer usuário autenticado pode atualizar seus próprios dados (username, email, password).
- Apenas administradores podem alterar o campo `role` de qualquer usuário.

- Path Params:
  - `id_user`: O ID do usuário.

- Headers:
  - `Authorization`: Bearer `jwt_token`

- Middlewares Aplicados:
  - [Auth Middleware](#auth-middleware)
  - [Require Admin](#require-admin)

- Objeto antes da atualização:
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
  - [Auth Middleware](#auth-middleware)
  - [Require Admin](#require-admin)

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

## <div id="useful-links">Links Úteis 🔗</div>

- [Documentação Oficial do Go (Golang)](https://golang.org/doc/)
- [Referência de módulos Go](https://blog.golang.org/using-go-modules)
- [Documentação do Gin Web Framework](https://gin-gonic.com/en/docs/)
- [Documentação do Docker](https://docs.docker.com/)
- [Documentação Oficial do PostgreSQL](https://www.postgresql.org/docs/)
- [Introdução ao JWT](https://jwt.io/introduction)
- [Insomnia](https://insomnia.rest/) / [Postman](https://www.postman.com/) — Ferramentas de testes de API

---

Desenvolvido por [Marcos Vinicius Boava](https://github.com/Mfrozzz), usando como base [Go Lab Tutoriais](https://www.youtube.com/watch?v=3p4mpId_ZU8&t=3317s).