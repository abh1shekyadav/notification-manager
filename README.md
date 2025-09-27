# Notification Manager

A scalable notification system that allows users to register, authenticate, and send notifications (SMS or Email) via Kafka-backed processing. User data is stored in PostgreSQL, and JWT authentication secures the API.

---

## Features

- **User Registration**: Register users with email and password.
- **JWT Authentication**: Secure endpoints with JWT-based auth.
- **Send Notifications**: Authenticated users can send SMS or email notifications.
- **Scalable Processing**: Notifications are sent to Kafka for scalable, asynchronous processing.
- **PostgreSQL Storage**: User data is persisted in a PostgreSQL database.
- **Modular Architecture**: Clean separation of concerns for user, notification, and auth logic.

---

## Setup

### 1. Prerequisites

- Go 1.24+
- PostgreSQL
- Kafka
- [lib/pq](https://github.com/lib/pq) (Postgres driver)

### 2. Environment Variables

Set the following environment variables:

- `DB_CONN` — PostgreSQL connection string (e.g., `postgres://user:password@localhost:5432/dbname?sslmode=disable`)
- `JWT_SECRET` — Secret key for JWT signing

Example:
```sh
export DB_CONN="postgres://postgres:password@localhost:5432/notificationdb?sslmode=disable"
export JWT_SECRET="your-very-secret-key"
```

### 3. Database Setup

Create the `users` table:

```sql
CREATE TABLE users (
    id UUID PRIMARY KEY,
    email TEXT UNIQUE NOT NULL,
    password TEXT NOT NULL,
    created_at TIMESTAMP NOT NULL
);
```

Create the `notifications` table:

```sql
CREATE TABLE notifications (
    id UUID PRIMARY KEY,
    user_id UUID NOT NULL,
    type VARCHAR(50) NOT NULL,
    payload TEXT NOT NULL,
    status VARCHAR(50) NOT NULL,
    created_at TIMESTAMP DEFAULT now()
);
```

### 4. Kafka Setup

- Ensure Kafka is running and accessible.
- Configure your Kafka topic(s) as needed (e.g., `notifications`).

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
- **Body**:
  ```json
  {
    "type": "email", // or "sms"
    "payload": {
      "to": "recipient@example.com",
      "subject": "Hello",
      "body": "This is a test email."
    }
  }
  ```
- **Response**: Notification status

---

## Security

- All endpoints except `/users/register` and `/auth/login` require a valid JWT in the `Authorization` header.
- Passwords are hashed before storage.

---

## Scaling

- Notification requests are published to Kafka for asynchronous, scalable processing.
- Consumers can be scaled independently to process notifications.

---

## Contributing

Feel free to open issues or submit pull requests!

---

## License

MIT