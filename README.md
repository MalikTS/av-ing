# av-ing

Доска объявлений для Республики Ингушетия. Проект состоит из REST API на Go
и адаптивного frontend-прототипа на Nuxt, Vue и TypeScript.

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

## Frontend-прототип

Требуется Node.js 24.11 или новее.

```powershell
cd frontend
npm install
npm run dev
```

Интерфейс откроется на `http://localhost:3000`. Сейчас главная страница работает
на демонстрационных данных: поиск, фильтры, сортировка и избранное интерактивны,
но запросы к REST API будут подключены на следующем этапе.

Проверка frontend:

```powershell
cd frontend
npm run typecheck
npm run build
```

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
