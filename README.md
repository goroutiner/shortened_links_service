[![codecov](https://codecov.io/gh/goroutiner/shortened_links_service/graph/badge.svg)](https://codecov.io/gh/goroutiner/shortened_links_service)

## 📖 Translations
- [Read in Russian](/README_RU.md)

---

<h3 align="center">
  <div align="center">
    <h1>Shortened Links Service</h1>
  </div>
</h3>

## 📋 Project Description

**Shortened Links Service** is a project that provides a service for shortening URLs. The service supports data storage in PostgreSQL and in-memory modes.

---

## 🚀 Running the Project

### 1️⃣ Installing Dependencies

---

*❗Before running the service, make sure you have **Docker** and **Docker Compose** installed.*

---

### 2️⃣ Environment Configuration

The **environment** variables are set by default, but you can change them in the `compose.yaml` file:

- For the `golang` service:
```yaml
...
environment:
    PORT: ":8080"   
    MODE: "postgres"
    DATABASE_URL: "postgres://root:password@postgres:5432/mydb?sslmode=disable"
...
```
If you need **in-memory** mode, specify `MODE: "in-memory"`.

- For the `postgres` service:
```yaml
...
environment:
    POSTGRES_USER: "root"
    POSTGRES_PASSWORD: "password"
    POSTGRES_DB: "mydb"
...
```

### 3️⃣ Running the Project

The project is started using `docker compose`:

```sh
make run
```

### 4️⃣ Stopping the Service

To stop the containers, run:

```sh
make stop
```

---

## 🔥 API Endpoints

### 1️⃣ Example: Creating a Short Link

**POST** `/api/v1/shorten`

#### **Request Body:**

```json
{
  "original_link": "https://finance.ozon.ru"
}
```

#### **Response Body:**

```json
{
  "short_link": "abc123XYZ_"
}
```

### 2️⃣ Example: Retrieving the Original Link

**GET** `/api/v1/{short_link}`

#### **Response Body:**

```json
{
  "original_link": "https://finance.ozon.ru"
}
```

---

## 🧪 Running Tests

### 1️⃣ Command to Run Unit Tests:

```sh
make unit-tests  
```

### 2️⃣ Command to Run Integration Tests:

---

*❗Before running these tests, make sure you have **Docker** installed and running.* 

--- 

```sh
make integration-tests
```

### 3️⃣ After running all tests, stop the PostgreSQL container and clear the cache:

```sh
make clean
```

---

## 🛠️ Technical Resources

- **Libraries for Database Interaction:** [jmoiron/sqlx](https://github.com/jmoiron/sqlx) and [jackc/pgx](https://github.com/jackc/pgx)

- **Library for Writing Tests:** [stretchr/testify](https://github.com/stretchr/testify)

- **Library for Limiting User RPS:** [golang.org/x/time/rate](https://pkg.go.dev/golang.org/x/time@v0.10.0/rate#pkg-overview)

---
