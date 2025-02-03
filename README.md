# 🚀 User Management API

## 📌 Overview
This project is a **Golang-based User Management API** using **Echo**, **PostgreSQL**, **GORM**, **RabbitMQ**, and **Swagger**. It supports user creation, retrieval, and message queuing via RabbitMQ.

This project follows the **Domain-Driven Design (DDD) pattern**, making it easy to maintain, test, and extend.

Additionally, **parallelism and graceful shutdown** have been implemented to enhance scalability and reliability in production environments.

---

## ⚡️ Quick Start

### **1️⃣ Clone the Repository**
```sh
git clone https://github.com/your-username/user-api.git
cd user-api
```

### **2️⃣ Create an `.env` File**
Ensure you have a **`.env`** file in the root directory.

```ini
APP_ENV=dev
POSTGRES_USER=test
POSTGRES_PASSWORD=test
POSTGRES_DB=user_db
POSTGRES_HOST=postgres
POSTGRES_PORT=5432
RABBITMQ_URL=amqp://guest:guest@localhost:5672/
```

### **3️⃣ Start the Services**
```sh
make start
```
✅ **This will:**
- Start **PostgreSQL**
- Start **RabbitMQ**
- Run **database migrations**
- Start the **User API server on `http://localhost:8080`**

---

## 📝 API Endpoints

### **1️⃣ User Endpoints**

| Method | Endpoint | Description |
|--------|---------|-------------|
| `POST` | `/users` | Create a new user |
| `GET`  | `/users` | Get all users with pagination & filters |
| `PUT`  | `/users/{id}` | Update user by ID |
| `DELETE` | `/users/{id}` | Delete user by ID |

### **2️⃣ Health Check Endpoint**
| Method | Endpoint | Description |
|--------|---------|-------------|
| `GET`  | `/health` | Returns API and database status |

Example Response:
```json
{
  "status": "healthy"
}
```

### **3️⃣ Swagger API Documentation**
Swagger documentation is available at:
👉 [**http://localhost:8080/swagger/index.html**](http://localhost:8080/swagger/index.html)

---

## 🛠 RabbitMQ Setup

### **1️⃣ Access RabbitMQ Admin Panel**
RabbitMQ's management interface is available at:
👉 [**http://localhost:15672/**](http://localhost:15672/)

🔑 **Login Credentials:**
- **Username:** `guest`
- **Password:** `guest`

### **2️⃣ RabbitMQ Queues Used**
| Queue Name     | Purpose |
|----------------|---------|
| `user.created` | Handles user creation events |
| `user.updated` | Handles user update events |
| `user.deleted` | Handles user deletion events |

---

## 🧪 Running Tests

### **1️⃣ Run Unit Tests**
```sh
make test-unit
```

### **2️⃣ Run Integration Tests**
```sh
make test-integration
```
✅ **This will:**
- Start **PostgreSQL (Test DB)**
- Run **GORM migrations**
- Execute **integration tests**

### **3️⃣ Run API Integration Tests**
```sh
make test-api
```
✅ **This will:**
- Start **PostgreSQL (Test DB)**
- Run **GORM migrations**
- Execute **API integration tests**
- Test `getUsers` endpoint to ensure pagination and filtering work correctly

---

## 🔄 Stopping the Services
```sh
make stop
```
✅ **This will:**
- Stop all running containers
- Preserve database data

---

## 💡 Additional Notes
- The project follows **RESTful API standards**.
- Uses **JWT authentication (if implemented)**.
- Logs are managed using **Zerolog**.
- All database interactions use **GORM ORM**.
- Implements **Domain-Driven Design (DDD) pattern** for better maintainability and testability.

---

## 💡 Things to add
- Store logs in Kibana for centralized log management.
- Improve log structure for better readability and debugging.
- Increase Test Coverage.
- User [cursor pagination](https://planetscale.com/blog/mysql-pagination) instead of OFFSET-based pagination.