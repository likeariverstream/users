## Users service

### Start:

Create .env file

.env file example:
```
APP_PORT=4000
POSTGRES_PASSWORD=example
POSTGRES_USER=user
POSTGRES_DB=users
POSTGRES_PORT=5433
ADMINER_PORT=8080
GOOSE_DRIVER=postgres
GOOSE_DBSTRING=postgres://user:example@postgres:5432/users
GOOSE_MIGRATION_DIR=./migrations
DB_DATA_SOURCE_NAME=postgres://user:example@postgres:5432/users?sslmode=disable
```

```shell
 docker-compose up -d
 ```
Swagger:

/swagger/index.html

Adminer server:

`users_postgres`

### Developing:

Create .env file

Install dependencies

PostgreSQL start:
```shell
 docker-compose -f docker-compose.dev.yml up -d
 ```
Adminer server `users_postgres`

Testing:

```go test -v```

Swagger generate:

```swag init -g main.go```

Migrations:
[Goose](https://github.com/pressly/goose)

Driver:

```shell
  GOOSE_DRIVER=postgres
```

Create .sql migration:

```shell
  goose create migration_name sql
```

Migrations up:
```shell
  goose up
```
Migrations down:
```shell
  goose down
```

Current DB version:
```shell
  goose version
```

