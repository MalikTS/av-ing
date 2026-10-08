# av-ing

REST API сервиса объявлений на Go, Chi и PostgreSQL.

## Возможности

- регистрация и JWT-авторизация;
- создание, просмотр, изменение и архивирование объявлений;
- поиск, фильтры, сортировка и пагинация;
- загрузка изображений;
- комментарии, избранное и отзывы.

Новые объявления сразу получают статус `active`. При внедрении модерации статус
нужно синхронно изменить на `pending` в обработчике создания и в схеме БД.

## Локальный запуск

1. Создайте базу PostgreSQL.
2. Примените `migrations/001_init.up.sql`.
3. Скопируйте `.env.example` в `.env` и заполните значения.
4. Запустите API:

```powershell
go run ./cmd/server
```

По умолчанию сервер доступен на `http://localhost:8080`.

## Основные маршруты

Публичные:

- `POST /register`, `POST /login`;
- `GET /ads`, `GET /ads/{id}`;
- `GET /ads/{id}/comments`, `GET /ads/{id}/reviews`;
- `GET /uploads/*`.

С заголовком `Authorization: Bearer <token>`:

- `GET /me`, `GET /me/favorites`;
- `POST /ads`, `PUT /ads/{id}`, `DELETE /ads/{id}`;
- `POST /ads/{id}/images` с multipart-полем `images`;
- создание и удаление комментариев;
- добавление и удаление избранного;
- создание отзывов.

## Проверка

```powershell
gofmt -w cmd internal
go test ./...
go vet ./...
```
