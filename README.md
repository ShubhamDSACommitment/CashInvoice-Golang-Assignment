# CashInvoice-Golang-Assignment
Golang Task Management Service (CashInvoice Assignment)

A scalable RESTful Task Management Service built in Go using clean architecture, JWT-based authentication, MySQL persistence, and concurrent background workers powered by goroutines and channels.

This project is fully containerized using Docker Compose, so you can run everything with a single command.


A **scalable RESTful Task Management Service** built in Go using clean architecture, JWT-based authentication, MySQL persistence, and concurrent background workers powered by goroutines and channels.

This project is fully containerized using **Docker Compose**, so you can run everything with a single command.

---

## 🚀 Features

- RESTful CRUD APIs for tasks
- JWT Authentication & Role-Based Authorization
    - **Users** → can access only their own tasks
    - **Admins** → can access all tasks
- MySQL persistence layer
- Background worker for **auto-completing tasks after X minutes**
- Clean layered architecture:
    - Handlers → Services → Repositories → Database
- Dockerized setup (API + DB)

---

## 🏗 Architecture Overview
```mermaid
flowchart TD
    C[Client<br/>(Postman / Curl)]
    R[Gin Router]
    M[JWT Middleware]
    H[Handlers<br/>(HTTP Layer)]
    S[Services<br/>(Business Logic)]
    REPO[Repositories<br/>(MySQL)]
    DB[(Database)]
    W[Background Worker<br/>(Goroutines + Channel Queue)]

    C --> R
    R --> M
    M --> H
    H --> S
    S --> REPO
    REPO --> DB

    %% Background worker flow
    S -->|Push Task ID| W
    W -->|Conditional Update| DB
```
---

## 🐳 Docker Compose Setup

Docker Compose runs:
- **MySQL 8** container
- **Golang API** container
- Internal network
- Health checks
- Persistent DB volume

---

## ▶️ Getting Started

### Prerequisites

Make sure you have:
- Docker
- Docker Compose
- Git

---

## 🔽 Clone the Repository

```bash
git clone https://github.com/ShubhamDSACommitment/CashInvoice-Golang-Assignment
cd CashInvoice-Golang-Assignment
```


## ▶️ Run the System
```
docker-compose up --build
```

## ✅ Verify Services
```
http://localhost:8080
```

## 🔐 Authentication Flow
### Register User
```POST http://localhost:8080/auth/register```

Request Body : 
```
{
  "email": "user@test.com",
  "password": "password123"
}
```

### Login
``` POST http://localhost:8080/auth/login ```

Request Body :
```
{
  "email": "user@test.com",
  "password": "password123"
}
```
Response :
```
{
  "token": "JWT_TOKEN",
  "role": "user"
}
```

## 📝 Task APIs

All task endpoints require this header:
``` Authorization: Bearer <JWT_TOKEN> ```

### ➕ Create Task
``` POST http://localhost:8080/tasks ```

**Request Body**
```json
{
  "title": "Finish Golang Assignment",
  "description": "Implement background worker"
}
```
### 📋 Get All Tasks
```
GET http://localhost:8080/tasks
```

### 🔍 Get Task by ID
```
GET http://localhost:8080/tasks/{id}
```

### ❌ Delete Task
```
DELETE http://localhost:8080/tasks/{id}
```

## 👑 Admin Features

### Admins can:

- View all tasks

- Delete any task

- Create other admins

### 🔐 Create First Admin (Bootstrap)
### Run this in any terminal: 
You will get logged in into the mysql terminal
```
docker exec -it task_mysql mysql -u taskuser -ptaskpass taskdb
```
Then:
```
INSERT INTO users (id, email, password, role)
VALUES (
  UUID(),
  'admin@system.com',
  '$2a$10$N9qo8uLOickgx2ZMRZo5i.ej2U9U5Q9F8vKcX4F1eO0xS3yYz5b6C',
  'admin'
);
```
This will insert admin user directly into the database 

## ⚡ Background Worker (Concurrency)

### How It Works

- When a task is created, its ID is pushed into a **buffered channel**
- Worker goroutines consume tasks from the queue
- Each worker waits **X minutes** (configurable via `AUTO_COMPLETE_MINUTES`)
- The worker performs an **atomic DB update**:
  - If task is still `pending` or `in_progress` → mark as `completed`
  - If task was deleted or manually completed → skip

---

### 🔐 Thread Safety

- Channels are thread-safe
- Workers run in isolated goroutines

## 🛠 Tech Stack

- Golang

- Gin

- MySQL 8

- Docker & Docker Compose

- JWT 

- bcrypt
