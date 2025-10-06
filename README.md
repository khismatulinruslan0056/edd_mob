# Effective Mobile API

API для работы с данными о людях. Сервис предоставляет стандартные CRUD-операции, а также обогащает данные о пользователе (возраст, пол, национальность) с помощью внешних публичных API.

##  Ключевые особенности

-   **CRUD API**
-   **Обогащение данных**
-   **Фильтрация и пагинация**
-   **Конфигурация через `.env`**
-   **Поддержка Docker**
-   **Миграции базы данных**
-   **Документация Swagger**
-   **Структурированное логирование**

## 🏛️ Архитектура и структура проекта

Проект следует принципам чистой архитектуры, разделяя логику на слои для лучшей поддерживаемости и тестируемости.

```
.
├── cmd/server/main.go    # Точка входа в приложение
├── internal
│   ├── config            # Конфигурация приложения
│   ├── httpserver        # HTTP-сервер, роутинг, хендлеры, middleware
│   │   ├── handlers      # Обработчики HTTP-запросов (CRUD)
│   │   └── middleware    # Промежуточное ПО (например, логгер)
│   ├── logger            # Пакет для логирования
│   ├── model             # Модели данных (сущности)
│   ├── service           # Бизнес-логика (например, обогащение данных)
│   │   └── enrichment
│   └── storage           # Слой для взаимодействия с базой данных
│       └── pg            # Реализация хранилища для PostgreSQL
├── migrations            # SQL-миграции для базы данных
├── docs                  # Сгенерированная документация Swagger
├── .env                  # Файл с переменными окружения (пример ниже)
├── Dockerfile            # Инструкции для сборки Docker-образа
├── docker-compose.yml    # Файл для оркестрации контейнеров (приложение + БД)
└── Makefile              # Утилиты для сборки, запуска и миграций
```

## Технологический стек

-   **Язык**: Go
-   **База данных**: PostgreSQL
-   **Веб-фреймворк**: Стандартная библиотека `net/http`
-   **Роутер**: Самописный роутер на базе `net/http`
-   **Конфигурация**: `godotenv`, `env`
-   **Миграции**: `golang-migrate`
-   **Документация API**: `swaggo/swag`
-   **Контейнеризация**: Docker, Docker Compose

## Начало работы

### Предварительные требования

-   [Go](https://golang.org/doc/install) (версия 1.21+)
-   [Docker](https://www.docker.com/get-started) и [Docker Compose](https://docs.docker.com/compose/install/)
-   [golang-migrate](https://github.com/golang-migrate/migrate/tree/master/cmd/migrate)

```bash
# Установка golang-migrate (для macOS/Linux)
curl -L https://github.com/golang-migrate/migrate/releases/download/v4.17.0/migrate.linux-amd64.tar.gz | tar xvz
sudo mv migrate /usr/local/bin/
```

### Установка и запуск

1.  **Клонируйте репозиторий:**
    ```bash
    git clone <your-repository-url>
    cd <repository-name>
    ```

2.  **Создайте файл конфигурации:**
    Создайте файл `.env` и заполните его своими данными.

    Пример файла `.env`:
    ```env
    # Отладка
    DEBUG=true

    # Конфигурация PostgreSQL
    DSN_PORT=5432
    DSN_USER=admin
    DSN_PASSWORD=123
    DSN_NAME=EffectiveMobileAPI
    DSN_HOST=localhost # Для локального запуска. При запуске через Docker Compose используйте имя сервиса, например 'db'

    # URL для миграций
    DATABASE_URL=postgres://admin:123@localhost:5432/EffectiveMobileAPI?sslmode=disable

    # Конфигурация HTTP-сервера
    HTTP_ADDR=localhost:7007
    HTTP_TIMEOUT=4s
    HTTP_IDLE_TIMEOUT=30s
    HTTP_USER=admin # Эти данные не используются в текущей реализации, но могут быть добавлены для Basic Auth
    HTTP_SERVER_PASSWORD=secret
    ```

### Способы запуска

#### 1. С помощью Docker Compose (рекомендуемый способ)

```bash
# Убедитесь, что в .env DSN_HOST=db
docker-compose up --build
```

Приложение будет доступно по адресу `http://localhost:7007`.

#### 2. Локальный запуск с помощью Makefile

1.  **Запустите PostgreSQL:** Убедитесь, что у вас есть работающая база данных, и ее параметры соответствуют указанным в `.env`.

2.  **Примените миграции:**
    ```bash
    make migrate-up
    ```

3.  **Запустите приложение:**
    ```bash
    make run
    ```

##  Makefile команды

-   `make build`: Собрать бинарный файл приложения.
-   `make run`: Запустить приложение (после сборки).
-   `make migrate-up`: Применить все доступные миграции.
-   `make migrate-down`: Откатить последнюю примененную миграцию.
-   `make clean`: Удалить собранный бинарник.
-   `make fmt`: Отформатировать код проекта.

##  API Документация

Сервис использует Swagger для документирования API. Интерактивная документация доступна после запуска приложения по адресу:

**[http://localhost:7007/swagger/index.html](http://localhost:7007/swagger/index.html)**

### Основные эндпоинты

| Метод  | Путь             | Описание                                                                  |
| :----- | :--------------- | :------------------------------------------------------------------------ |
| `GET`  | `/people`        | Получить список людей с возможностью фильтрации и пагинации.              |
| `POST` | `/people`        | Добавить нового человека. Данные обогащаются (возраст, пол, национальность). |
| `PUT`  | `/people/{id}`   | Обновить данные человека по его ID.                                       |
| `DELETE`| `/people/{id}`  | Удалить человека по его ID.                                               |
| `GET`  | `/health`        | Проверка работоспособности сервиса.                                       |

#### Пример запроса на создание пользователя (`POST /people`)

```json
{
  "name": "Dmitriy",
  "surname": "Ivanov",
  "patronymic": "Sergeevich"
}
```

#### Пример ответа

```json
{
    "id": 1,
    "message": "user added"
}
```

#### Пример запроса на получение списка пользователей (`GET /people`)

Вы можете использовать query-параметры для фильтрации:

`GET /people?age=30&gender=male&limit=10`