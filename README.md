# SisterGo — Laravel-inspired Go starter

This project is a Go port of the Laravel PHP app (sisterlaravel12).

## Features

- Gin HTTP router
- GORM (Postgres) for models & migrations
- JWT authentication
- User, Role, Menu, Submenu, Settings, and Student models
- Dashboard with system info and statistics
- User activity logging
- Admin panel functionality

## Quick Start (Local)

1. Copy `.env.example` to `.env` and update values.

2. Using Docker Compose:

```bash
docker-compose up --build
```

3. The server will be available at `http://localhost:8080`.

## API Endpoints

### Public Endpoints

- `POST /api/register` - Register a new user `{name, username, email, password}`
- `POST /api/login` - Login `{identifier, password}` (identifier can be username or email)
- `GET /api/health` - Health check

### Protected Endpoints (Requires Bearer Token)

#### User Profile

- `GET /api/me` - Get current user
- `GET /api/profile` - Get user profile
- `PUT /api/profile` - Update profile `{name, email}`
- `POST /api/profile/change-password` - Change password `{current_password, new_password, confirm_password}`

#### Dashboard

- `GET /api/dashboard` - Get dashboard data with stats and recent activity

#### Students

- `GET /api/students` - List all students
- `POST /api/students` - Create student `{name, email, nis}`
- `GET /api/students/:id` - Get student by ID
- `PUT /api/students/:id` - Update student
- `DELETE /api/students/:id` - Delete student

### Admin Endpoints (Requires role_id == 1)

#### User Management

- `GET /api/admin/users` - List users with pagination
- `POST /api/admin/users` - Create user `{name, username, email, password, role_id}`
- `GET /api/admin/users/:id` - Get user by ID
- `PUT /api/admin/users/:id` - Update user
- `DELETE /api/admin/users/:id` - Delete user (cannot delete admin)

#### Role Management

- `GET /api/admin/roles` - List all roles
- `POST /api/admin/roles` - Create role `{role}`
- `GET /api/admin/roles/:id` - Get role by ID
- `PUT /api/admin/roles/:id` - Update role
- `DELETE /api/admin/roles/:id` - Delete role

#### Menu Management

- `GET /api/admin/menus` - List all menus with submenus
- `POST /api/admin/menus` - Create menu `{menu, icon, ordering}`
- `PUT /api/admin/menus/:id` - Update menu
- `DELETE /api/admin/menus/:id` - Delete menu

#### Submenu Management

- `GET /api/admin/submenus` - List all submenus
- `POST /api/admin/submenus` - Create submenu `{menu_id, title, url, is_active, ordering}`
- `PUT /api/admin/submenus/:id` - Update submenu
- `DELETE /api/admin/submenus/:id` - Delete submenu

#### Settings

- `GET /api/admin/settings` - Get all settings
- `POST /api/admin/settings` - Create/Update setting `{name, value, is_active}`

#### User Activity & Logs

- `GET /api/admin/user-activity` - Get user login activity
- `GET /api/admin/logs` - Get all activity logs with pagination
- `GET /api/admin/logs/user/:user_id` - Get logs for specific user
- `POST /api/admin/logs` - Create log entry `{user_id, activity}`

## Environment Variables

| Variable | Description | Default |
|----------|-------------|---------|
| APP_PORT | Server port | 8080 |
| APP_ENV | Environment (development/production) | production |
| DB_HOST | Database host | localhost |
| DB_PORT | Database port | 5432 |
| DB_USER | Database user | postgres |
| DB_PASSWORD | Database password | password |
| DB_NAME | Database name | sistergo |
| JWT_SECRET | JWT signing secret | secret |

## Usage Examples

### Register a new user

```bash
curl -X POST http://localhost:8080/api/register \
  -H "Content-Type: application/json" \
  -d '{"name":"Admin","username":"admin","email":"admin@example.com","password":"secretpass"}'
```

### Login

```bash
curl -X POST http://localhost:8080/api/login \
  -H "Content-Type: application/json" \
  -d '{"identifier":"admin","password":"secretpass"}'
```

### Use the token for authenticated routes

```bash
curl -H "Authorization: Bearer <TOKEN>" http://localhost:8080/api/profile
```

### Get dashboard data

```bash
curl -H "Authorization: Bearer <TOKEN>" http://localhost:8080/api/dashboard
```

### Admin: List users

```bash
curl -H "Authorization: Bearer <TOKEN>" http://localhost:8080/api/admin/users
```

## Default Admin User

On first startup, the system seeds a default admin user:
- Username: `admin`
- Email: `admin@example.com`
- Password: `password`
- Role: Admin (role_id: 1)
