# Users API

RESTful API на Go для управления пользователями и отношениями "друзья".

В данном проекте реализованы CRUD-овые операции для пользователей, а также есть возможность создания связи между ними(дружбы). При удалении пользователя, состоящего в дружбе, реализовано каскадное удаление этой связи.

**Технологии:**
- Go (Golang)
- PostgreSQL (основная база)
- SQLite (база сугубо для тестов)
- GORM (ORM)
- Chi Router (маршрутизатор)
- REST API

**Структура проекта:**
- `model/` — схема базы,в ней описаны 2 структуры данных user и friendship c указанием зависимостей для Gorm
- `repository/` — слой работы с БД
- `service/` — бизнес-логика
- `handler/` — HTTP-хендлеры
- `config/` — инициализация подключения к БД
- `migrations` - хранилище sql-миграции схемы БД

## Установка и запуск

git clone https://github.com/UnendingLoop/users-api.git
cd users-api

# установка зависимостей
go mod tidy

# запуск локально
go run cmd/main.go

## Переменные окружения

- `DATABASE_URL` — строка подключения к PostgreSQL, указывается в .env рядом с main.go

## Примеры API-запросов
# Создание пользователя
curl -X POST http://localhost:8080/users \
  -H "Content-Type: application/json" \
  -d '{
    "name": "John",
    "surname": "Smith",
    "email": "newjohn@example.com"
}'

curl -X POST http://localhost:8080/users \
-H "Content-Type: application/json" \
-d '{
    "name":"Alice",
    "surname":"Smith",
    "email":"alice@example.com"
}'

# Обновление информации о существующем пользователе:
curl -X PUT http://localhost:8080/update/1 \
  -H "Content-Type: application/json" \
  -d '{
    "name": "john",
    "surname": "smith",
    "email": "newjohn@example.com"
}'

# Получение информации о существующем пользователе:
curl -X GET http://localhost:8080/users/1

# Получение списка всех существующих пользователей:
curl http://localhost:8080/users

# Удаление пользователя:
curl -X DELETE http://localhost:8080/delete/1

# Создание дружбы между двумя пользователями:
curl -X POST http://localhost:8080/users/1/make_friend/2

# Просмотр списка друзей пользователя:
curl http://localhost:8080/users/1/friends

# Удаление дружбы:
curl -X DELETE http://localhost:8080/users/1/remove_friend/2


## Тестирование

- `test.db` используется как SQLite in-memory база для тестов
- Тесты запускаются через:
go test ./...


## TODO

- [ ] JWT-авторизация
- [ ] Swagger-документация
- [ ] Интеграционные тесты
- [ ] Пагинация и поиск
