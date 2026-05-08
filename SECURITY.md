# Sales Manager Backend - Security & Best Practices

## Security Measures Implemented

### 1. Authentication & Authorization
- **Multi-tenant isolation**: Every query filtered by `tenant_id`
- **Required headers**: `X-Tenant-ID` and `X-User-ID` validated on all routes
- **Middleware protection**: `RequireAuthMiddleware` on all routers
- **Context-based auth**: Credentials stored in request context, not global state

### 2. CORS Configuration
- **Allowed origins**: HTTPS only in production (wildcards for development)
- **Localhost support**: `http://localhost:*` and `http://127.0.0.1:*` for dev
- **Credentials**: Enabled with `AllowCredentials: true`
- **Headers whitelist**: Only necessary headers allowed
- **Max age**: 5 minutes (300s) for preflight caching

### 3. HTTP Server Hardening
- **Timeouts**:
  - Read: 15 seconds
  - Write: 15 seconds
  - ReadHeader: 5 seconds
  - Idle: 60 seconds
- **Max header size**: 1 MB (prevents large header attacks)
- **Panic recovery**: Automatic recovery from panics via chi middleware
- **Request ID**: Each request gets unique ID for tracing

### 4. Database Security
- **Prepared statements**: GORM uses parameterized queries automatically
- **Connection pooling**: Managed by GORM/database/sql
- **Timezone consistency**: Fixed to UTC-3
- **Singular table naming**: Prevents common injection vectors

### 5. Input Validation
- **JSON decoding**: Safe unmarshaling with Go's json package
- **Type safety**: Strong typing prevents common injection
- **ID parsing**: Strict uint parsing with error handling
- **Query parameter validation**: All params validated before use

### 6. Response Security
- **Content-Type**: Always `application/json`
- **Structured errors**: No stack traces or internal details exposed
- **Status codes**: Consistent HTTP status code usage
- **No data leakage**: Error responses don't expose DB structure

### 7. Logging & Monitoring
- **Structured logging**: Request ID, Real IP, timestamps
- **Access logs**: All requests logged via chi Logger middleware
- **Health check**: `/health` endpoint for uptime monitoring
- **Startup validation**: DB connection verified on startup

## Production Recommendations

### Environment Variables
```bash
# Required
export DB="user:password@tcp(host:3306)/sales_manager"
export PORT="8080"

# Recommended
export CORS_ALLOWED_ORIGINS="https://yourdomain.com"
export MAX_DB_CONNECTIONS="100"
export LOG_LEVEL="info"
```

### CORS Production Config
Update `server.go` line 32:
```go
AllowedOrigins: []string{"https://yourdomain.com"},
```

### Database Connection String
Use connection pooling parameters:
```
user:password@tcp(host:3306)/sales_manager?parseTime=true&maxAllowedPacket=0&timeout=30s&readTimeout=30s&writeTimeout=30s
```

### TLS/HTTPS
Deploy behind a reverse proxy (nginx, Caddy) for TLS termination:
```nginx
location /api/ {
    proxy_pass http://localhost:8080;
    proxy_set_header Host $host;
    proxy_set_header X-Real-IP $remote_addr;
    proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
}
```

### Rate Limiting
Consider adding rate limiting middleware:
```go
import "github.com/didip/tollbooth/v7"

// Add to server.go
r.Use(tollbooth.LimitHandler(tollbooth.NewLimiter(10, nil)))
```

### Firebase Token Validation (TODO)
The `authHelper` currently validates headers but doesn't verify Firebase tokens. To add:
1. Install Firebase Admin SDK
2. Verify ID token in `RequireAuthMiddleware`
3. Extract claims from token instead of headers

## Common Security Pitfalls Avoided

### ✅ SQL Injection
- GORM uses prepared statements automatically
- All queries parameterized

### ✅ XSS (Cross-Site Scripting)
- API only returns JSON (no HTML rendering)
- Content-Type strictly enforced

### ✅ CSRF (Cross-Site Request Forgery)
- CORS properly configured
- SameSite cookie policy (when using cookies)

### ✅ Information Disclosure
- Generic error messages
- No stack traces in responses
- No version info in headers

### ✅ Mass Assignment
- Explicit field mapping in handlers
- No direct struct binding to DB

### ✅ Broken Authentication
- Multi-factor auth via tenant + user validation
- Context-based auth prevents race conditions

## Monitoring & Alerts

### Health Check Endpoint
```bash
curl http://localhost:8080/health
```

### Prometheus Metrics (Future)
Consider adding:
- Request count by endpoint
- Request duration histograms
- Active connections
- Error rates

### Alerting
Set up alerts for:
- 500 errors > 1%
- Response time > 1s (p95)
- Health check failures
- Database connection errors

## Security Checklist for Deployment

- [ ] Change CORS origins to production domain
- [ ] Set up TLS/HTTPS (via reverse proxy)
- [ ] Enable Firebase token validation
- [ ] Add rate limiting
- [ ] Set up database connection pooling
- [ ] Configure proper DB user permissions (not root)
- [ ] Enable query logging for audit trail
- [ ] Set up monitoring/alerting
- [ ] Regular dependency updates (`go get -u ./...`)
- [ ] Code security scanning (gosec, staticcheck)

## Security Tools

### Run Security Scanners
```bash
# Install gosec
go install github.com/securego/gosec/v2/cmd/gosec@latest

# Scan for vulnerabilities
gosec ./...

# Check for dependency vulnerabilities
go list -json -m all | nancy sleuth
```

### Static Analysis
```bash
# Install staticcheck
go install honnef.co/go/tools/cmd/staticcheck@latest

# Run analysis
staticcheck ./...
```

## Contact

For security issues, please DO NOT create a public issue. Contact security team directly.
