# Notification Manager

A scalable notification system that allows users to register, authenticate, and send notifications (SMS or Email) via Kafka-backed processing. User data is stored in PostgreSQL, and JWT authentication secures the API.

---

## Features

- **User Registration**: Register users with email and password.
- **JWT Authentication**: Secure endpoints with JWT-based auth.
- **Send Notifications**: Authenticated users can send SMS or email notifications.
- **Scalable Processing**: Notifications are sent to separate Kafka topics for SMS and Email with multiple partitions and consumers for scalable, asynchronous processing.
- **PostgreSQL Storage**: User data is persisted in a PostgreSQL database with updated_at timestamps.
- **Health Checks & Metrics**: Exposes `/healthz` for health checks and `/metrics` for Prometheus metrics.
- **Modular Architecture**: Clean separation of concerns for user, notification, and auth logic.

---

## Setup

### 1. Prerequisites

- Go 1.24+
- PostgreSQL
- Kafka
- [lib/pq](https://github.com/lib/pq) (Postgres driver)

### 2. Environment Variables

Configure the application exclusively via environment variables:

- `DB_CONN` — PostgreSQL connection string (e.g., `postgres://user:password@localhost:5432/dbname?sslmode=disable`)
- `JWT_SECRET` — Secret key for JWT signing
- `KAFKA_BROKERS` — Comma-separated list of Kafka brokers (e.g., `localhost:9092`)
- `KAFKA_SMS_TOPIC` — Kafka topic for SMS notifications (e.g., `sms_notifications`)
- `KAFKA_EMAIL_TOPIC` — Kafka topic for Email notifications (e.g., `email_notifications`)
- `KAFKA_DLQ_TOPIC` — Kafka Dead Letter Queue topic (e.g., `notifications_dlq`)

Example:
```sh
export DB_CONN="postgres://postgres:password@localhost:5432/notificationdb?sslmode=disable"
export JWT_SECRET="your-very-secret-key"
export KAFKA_BROKERS="localhost:9092"
export KAFKA_SMS_TOPIC="sms_notifications"
export KAFKA_EMAIL_TOPIC="email_notifications"
export KAFKA_DLQ_TOPIC="notifications_dlq"
```

### 3. Database Setup

Create the `users` table:

```sql
CREATE TABLE users (
    id UUID PRIMARY KEY,
    email TEXT UNIQUE NOT NULL,
    password TEXT NOT NULL,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL
);
```

Create the `notifications` table:

```sql
CREATE TABLE notifications (
    id UUID PRIMARY KEY,
    user_id UUID NOT NULL REFERENCES users(id),
    type VARCHAR(50) NOT NULL,
    payload TEXT NOT NULL,
    status VARCHAR(50) NOT NULL,
    created_at TIMESTAMP DEFAULT now(),
    updated_at TIMESTAMP DEFAULT now()
);
```

### 4. Kafka Setup

- Ensure Kafka is running and accessible.
- Configure your Kafka topics with multiple partitions for scalability:
  - SMS notifications topic (e.g., `sms_notifications`)
  - Email notifications topic (e.g., `email_notifications`)
  - Dead Letter Queue topic (e.g., `notifications_dlq`)

---

### 5. Docker Compose

You can use Docker Compose to quickly spin up PostgreSQL and Kafka for local development.

Start services with:

```sh
docker-compose up -d --build
```

---

## Running the Application

```sh
go run main.go
```

The server will start on `:8080`.

---

## API Endpoints

### User Registration

- **POST** `/users/register`
- **Body**:
  ```json
  {
    "email": "user@example.com",
    "password": "yourpassword"
  }
  ```
- **Response**: User info (without password)

### User Login

- **POST** `/auth/login`
- **Body**:
  ```json
  {
    "email": "user@example.com",
    "password": "yourpassword"
  }
  ```
- **Response**: JWT token

### Get User by Email

- **GET** `/users?email=user@example.com`
- **Headers**: `Authorization: Bearer <token>`
- **Response**: User info (without password)

### Send Notification

- **POST** `/notify`
- **Headers**: `Authorization: Bearer <token>`
- **Body (Email Notification)**:
  ```json
  {
    "type": "email",
    "payload": {
      "to": "recipient@example.com",
      "subject": "Hello",
      "body": "This is a test email."
    }
  }
  ```
- **Body (SMS Notification)**:
  ```json
  {
    "type": "sms",
    "payload": {
      "to": "+1234567890",
      "message": "This is a test SMS."
    }
  }
  ```
- **Response**: Notification status

### Health Check

- **GET** `/healthz`
- **Response**: Service health status

### Prometheus Metrics

- **GET** `/metrics`
- **Response**: Prometheus metrics for monitoring

---

## Security

- All endpoints except `/users/register` and `/auth/login` require a valid JWT in the `Authorization` header.
- Passwords are hashed before storage.

---

## Scaling

- Notification requests are published to separate Kafka topics for SMS and Email for asynchronous, scalable processing.
- Multiple partitions and consumers enable horizontal scaling.
- Dead Letter Queue topic handles failed notification messages for later inspection.

---

## Contributing

Feel free to open issues or submit pull requests!

---

## License

MIT