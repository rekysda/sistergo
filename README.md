# SisterGo — Laravel-inspired Go starter

This project is a starting point to port a Laravel PHP app (sisterlaravel12) into Go.

Features:
- Gin HTTP router
- GORM (Postgres) for models & migrations
- JWT authentication
- Basic User and Student models

Quick start (local):

1. Copy `.env.example` to `.env` and update values.

2. Using Docker Compose:

```powershell
docker-compose up --build
```

3. The server will be available at `http://localhost:8080`.

API endpoints:
- POST /api/register {name, email, password}
- POST /api/login {email, password}
- GET /api/health
- GET /api/me (Bearer token)
- Students (Bearer token):
  - GET /api/students
  - POST /api/students {name, email, nis}
  - GET /api/students/:id
  - PUT /api/students/:id
  - DELETE /api/students/:id
   - Admin endpoints (role_id==1)
     - GET /api/admin/roles
     - POST /api/admin/roles
     - PUT /api/admin/roles/:id
     - DELETE /api/admin/roles/:id
     - GET /api/admin/menus
     - POST /api/admin/menus
     - PUT /api/admin/menus/:id
     - DELETE /api/admin/menus/:id
     - POST /api/admin/submenus
     - PUT /api/admin/submenus/:id
     - DELETE /api/admin/submenus/:id
     - GET /api/admin/settings
     - POST /api/admin/settings

Notes & next steps:
- You can map models and features from the original Laravel app into this skeleton: controllers, services, web templates, policies, etc.
- Implement role-based access, validation, and richer user features as needed.

If you'd like, I can start porting specific features from your Laravel repo (models, controllers, views) into this Go version — tell me which modules to prioritize.

Simple usage examples:

Register a new user (username is required for legacy behavior):

```powershell
curl -X POST http://localhost:8080/api/register -H "Content-Type: application/json" -d '{"name":"Admin","username":"admin","email":"admin@example.com","password":"secretpass"}'
```

Login and extract token:

```powershell
curl -X POST http://localhost:8080/api/login -H "Content-Type: application/json" -d '{"identifier":"admin","password":"secretpass"}'
```

Use the token for authenticated routes:

```powershell
curl -H "Authorization: Bearer <TOKEN>" http://localhost:8080/api/profile
```
