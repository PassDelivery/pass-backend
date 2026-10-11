# PassDelivery — Backend

REST API сервиса доставки еды PassDelivery. Командный проект в рамках курса ООП (МГТУ «СТАНКИН», ИИТ)

## Стек

- **Go** — язык реализации
- **PostgreSQL** — хранилище данных
- **Docker Compose** — локальный запуск
- **JWT** — аутентификация, **bcrypt** — хранение паролей

## Архитектура

Трёхслойная архитектура: `handler → service → repository`.

```
internal/
  core/              — общая инфраструктура (логгер, middleware, http-сервер, jwt, errors)
  <feature>/          — домен, репозиторий, сервис, http-хендлер на каждую фичу (users, restaurants, ...)
```

## Запуск

Требуется установленный Docker и Docker Compose.

```bash
git clone https://github.com/PassDelivery/pass-backend.git
make env-up
make migrate-up
make forwarder-up
make pass-run
```

Сервер поднимется на `http://localhost:5050`.

### Переменные окружения

Скопировать `.env.example` в `.env` и при необходимости поправить значения (порт, строка подключения к БД, секрет для JWT).

## API

Контракт API описан в OpenAPI-спецификации: [`pass-docs/openapi.yaml`](https://github.com/PassDelivery/pass-docs/blob/main/openapi.yaml).

## Схема данных

ERD базы данных: [`pass-docs/erd.md`](https://github.com/PassDelivery/pass-docs/blob/main/erd.md).

## Git-флоу

Разработка — через feature-ветки (`feature/<название>`) и Pull Request в `main`. Прямые коммиты в `main` запрещены (branch protection).
