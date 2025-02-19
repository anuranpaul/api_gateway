# API Gateway Service

A robust API Gateway built with Go, featuring role-based access control, rate limiting, metrics, and user management.

## Features

- 🔐 JWT-based Authentication
- 👥 Role-Based Access Control (RBAC)
- 🚦 Rate Limiting (Token and IP-based)
- 📊 Prometheus Metrics
- 👤 User Management with PostgreSQL
- 🔄 Reverse Proxy for Microservices

## API Endpoints

### User Management
- `POST /api/users` - Create user
- `GET /api/users` - List all users (Admin only)
- `GET /api/users/:id` - Get user details
- `PUT /api/users/:id` - Update user
- `DELETE /api/users/:id` - Delete user

### Metrics
- `GET /metrics` - Prometheus metrics (Admin only)

### Service Routes
- `/users/*` - User service endpoints
- `/admin/*` - Admin service endpoints

## Testing

### Authentication

- User access
```bash
curl -H "Authorization: Bearer <user-token>" http://localhost:8080/users/test
```
- Admin access
```bash
curl -H "Authorization: Bearer <admin-token>" http://localhost:8080/admin/dashboard
```

## Data Flow

1. **Client Request Flow**:
   ```mermaid
   sequenceDiagram
       Client->>API Gateway: HTTP Request
       API Gateway->>Auth Middleware: Verify JWT
       Auth Middleware->>Rate Limiter: Check Limits
       Rate Limiter->>Reverse Proxy: Route Request
       Reverse Proxy->>Microservice: Forward Request
       Microservice->>Client: Response
   ```

2. **User Management Flow**:
   ```mermaid
   sequenceDiagram
       Client->>API Gateway: User CRUD Request
       API Gateway->>Auth Middleware: Verify Permissions
       Auth Middleware->>User Management: Process Request
       User Management->>PostgreSQL: Database Operation
       PostgreSQL->>Client: Response
   ```

## Security Features

1. **Authentication**: JWT-based token authentication
2. **RBAC**: Role-based access control (admin/user roles)
3. **Rate Limiting**: 
   - Token-based: 20 requests per 30 seconds
   - IP-based: 30 requests per 60 seconds
4. **Metrics Protection**: Admin-only access to metrics

## Contributing

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add some amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request




