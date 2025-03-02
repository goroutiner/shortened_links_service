<h3 align="center">
  <div align="center">
    <h1>Shortened Links Service </h1>
  </div>
  </a>
</h3>

## 📋 Описание проекта

**Shortened Links Service** - это проект, представляющий собой сервис для сокращения URL-адресов. Сервис поддерживает хранение данных в PostgreSQL и in-memory режимах.

---

## 🚀 Запуск проекта

### 1️⃣ Установка зависимостей

---

*❗Перед запуском сервиса убедитесь, что у вас установлен **Docker** и **Docker Compose**.*

---

### 2️⃣ Конфигурация окружения 

Переменные окружения **environment** установлены по умолчанию, но вы их можете изменить в файле `compose.yaml`:

- Для сервиса `golang`:
```yaml
...
environment:
    PORT: ":8080"   
    MODE: "postgres"
    DATABASE_URL: "postgres://root:password@postgres:5432/mydb?sslmode=disable"
...
```
Если необходим **in-memory** режим, то укажите `MODE: "in-memory"`.

- Для сервиса `postgres`:
```yaml
...
environment:
    POSTGRES_USER: "root"
    POSTGRES_PASSWORD: "password"
    POSTGRES_DB: "mydb"
...
```
### 3️⃣ Запуск проекта

Проект запускается с помощью `docker compose`:

```sh
 make run
```

### 4️⃣ Остановка сервиса

Для остановки работы контейнеров выполните:

```sh
 make stop
```

---

## 🔥 API Эндпоинты

### 1️⃣ Пример создание короткой ссылки

**POST** `/api/v1/shorten`

#### **Тело запроса:**

```json
{
  "original_link": "https://finance.ozon.ru"
}
```

#### **Тело ответа:**

```json
{
  "short_link": "abc123XYZ_"
}
```

### 2️⃣ Пример получение оригинальной ссылки

**GET** `/api/v1/{short_link}`

#### **Тело ответа:**

```json
{
  "original_link": "https://finance.ozon.ru"
}
```

---

## 🧪 Запуск тестов

### 1️⃣ Команда для запуска unit-тестов:

```sh
make unit-tests  
```

### 2️⃣ Команда для запуска integration-тестов:

---

*❗Перед запуском данных тестов вы должны быть уверены, что у вас установлен и запущен **Docker**.* 

--- 

```sh
make integration-tests
```

### 3️⃣ После выполнения всех тестов нужно остановить контейнер с PosgreSQL и почистить кеш: 

```sh
make clean
```

---

## 🛠️ Технические ресурсы

- **Библиотеки для взаимодействия с БД:** [jmoiron/sqlx](https://github.com/jmoiron/sqlx) и [ackc/pgx](https://github.com/jackc/pgx)

- **Библиотека для написания тестов:** [stretchr/testify](https://github.com/stretchr/testify)

- **Библиотека для ограничения RPS пользователей сервиса:** [golang.org/x/time/rate](https://pkg.go.dev/golang.org/x/time@v0.10.0/rate#pkg-overview)
