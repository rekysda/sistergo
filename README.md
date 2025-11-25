# SisterGo — Laravel-inspired Go starter

This project is a Go port of the [sisterlaravel12](https://github.com/rekysda/sisterlaravel12) Laravel PHP application.

## Features

- **Gin HTTP router** - Fast and lightweight HTTP framework
- **GORM (Postgres)** - ORM for models & migrations
- **JWT authentication** - Secure token-based authentication
- **Role-based access control** - Admin and user roles
- **CORS support** - Cross-Origin Resource Sharing enabled
- Complete user management system
- Menu and submenu management
- Settings management
- Activity logging

## Quick Start

### Prerequisites

- Go 1.21 or higher
- PostgreSQL database
- Docker (optional)

### Local Development

1. Copy `.env.example` to `.env` and update values:

```bash
cp .env.example .env
```

2. Set up your database credentials in `.env`:

```
APP_PORT=8080
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=password
DB_NAME=sistergo
JWT_SECRET=your-secret-key
```

3. Run the application:

```bash
go run cmd/server/main.go
```

### Using Docker Compose

```bash
docker-compose up --build
```

The server will be available at `http://localhost:8080`.

## API Endpoints

### Public Endpoints

| Method | Endpoint | Description |
|--------|----------|-------------|
| POST | `/api/register` | Register a new user |
| POST | `/api/login` | Login and get JWT token |
| GET | `/api/health` | Health check |

### Protected Endpoints (Bearer Token Required)

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/api/me` | Get current user info |
| GET | `/api/dashboard` | Get dashboard statistics |
| GET | `/api/profile` | Get user profile |
| PUT | `/api/profile` | Update user profile |
| POST | `/api/profile/change-password` | Change password |

### Student Management

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/api/students` | List all students |
| POST | `/api/students` | Create a student |
| GET | `/api/students/:id` | Get a student |
| PUT | `/api/students/:id` | Update a student |
| DELETE | `/api/students/:id` | Delete a student |

### Admin Endpoints (Admin Role Required)

#### User Management

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/api/admin/users` | List all users |
| POST | `/api/admin/users` | Create a user |
| GET | `/api/admin/users/:id` | Get a user |
| PUT | `/api/admin/users/:id` | Update a user |
| DELETE | `/api/admin/users/:id` | Delete a user |

#### Role Management

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/api/admin/roles` | List all roles |
| POST | `/api/admin/roles` | Create a role |
| GET | `/api/admin/roles/:id` | Get a role |
| PUT | `/api/admin/roles/:id` | Update a role |
| DELETE | `/api/admin/roles/:id` | Delete a role |

#### Menu Management

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/api/admin/menus` | List all menus |
| POST | `/api/admin/menus` | Create a menu |
| PUT | `/api/admin/menus/:id` | Update a menu |
| DELETE | `/api/admin/menus/:id` | Delete a menu |

#### Submenu Management

| Method | Endpoint | Description |
|--------|----------|-------------|
| POST | `/api/admin/submenus` | Create a submenu |
| PUT | `/api/admin/submenus/:id` | Update a submenu |
| DELETE | `/api/admin/submenus/:id` | Delete a submenu |

#### Settings Management

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/api/admin/settings` | Get all settings |
| POST | `/api/admin/settings` | Update a setting |

#### Activity & Logs

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/api/admin/user-activity` | Get user activity |
| GET | `/api/admin/logs` | List activity logs |
| POST | `/api/admin/logs` | Create a log entry |

## Usage Examples

### Register a new user

```bash
curl -X POST http://localhost:8080/api/register \
  -H "Content-Type: application/json" \
  -d '{"name":"Admin","username":"admin","email":"admin@example.com","password":"secretpass"}'
```

### Login and get token

```bash
curl -X POST http://localhost:8080/api/login \
  -H "Content-Type: application/json" \
  -d '{"identifier":"admin","password":"secretpass"}'
```

### Use the token for authenticated routes

```bash
curl -H "Authorization: Bearer <TOKEN>" http://localhost:8080/api/profile
```

### Get dashboard statistics

```bash
curl -H "Authorization: Bearer <TOKEN>" http://localhost:8080/api/dashboard
```

### Admin: List all users

```bash
curl -H "Authorization: Bearer <TOKEN>" http://localhost:8080/api/admin/users
```

## Default Admin Account

When the database is first seeded, a default admin account is created:

- **Username:** admin
- **Email:** admin@example.com
- **Password:** password

**Note:** Change this password immediately in production!

## Project Structure

```
sistergo/
├── cmd/
│   └── server/
│       └── main.go          # Application entry point
├── internal/
│   ├── config/
│   │   └── config.go        # Configuration management
│   ├── controllers/
│   │   ├── auth.go          # Authentication handlers
│   │   ├── dashboard.go     # Dashboard handler
│   │   ├── log.go           # Activity log handlers
│   │   ├── menu.go          # Menu handlers
│   │   ├── profile.go       # Profile handlers
│   │   ├── role.go          # Role handlers
│   │   ├── setting.go       # Settings handlers
│   │   ├── student.go       # Student handlers
│   │   └── user.go          # User management handlers
│   ├── database/
│   │   └── db.go            # Database connection & seeding
│   ├── middleware/
│   │   ├── auth_middleware.go   # JWT authentication
│   │   ├── cors.go              # CORS middleware
│   │   └── role_middleware.go   # Role-based access control
│   ├── models/
│   │   ├── student.go
│   │   ├── user.go
│   │   ├── user_access_menu.go
│   │   ├── user_access_submenu.go
│   │   ├── user_log.go
│   │   ├── user_menu.go
│   │   ├── user_role.go
│   │   ├── user_submenu.go
│   │   └── web_setting.go
│   ├── routes/
│   │   └── routes.go        # Route definitions
│   └── utils/
│       ├── hash.go          # Password hashing
│       └── jwt.go           # JWT utilities
├── .env.example
├── docker-compose.yml
├── Dockerfile
├── go.mod
└── README.md
```

## License

This project is open-sourced software.
