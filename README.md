Разработать REST API на Go для создания, редактирования и удаления пользователей.
Testing:

```go test -v```

Для разработки: 
```shell
 docker-compose -f docker-compose.dev.yml up -d
 ```
Хост для adminer `users_postgres`

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

Framework на ваше усмотрение (можно без него);

Использовать PostgreSQL;

Разместить проект на GitHub.

Архитектура, тесты, контейнеризация.

Методы API:

POST /users – создать пользователя

GET /users/ – получить информацию о пользователе

PUT /users/ – обновить данные пользователя