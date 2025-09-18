# Project Conventions

## Architecture Patterns

### Layered Architecture
```
cmd/           → Application entry points
internal/      → Private application code
  ├── api/     → HTTP/gRPC handlers
  ├── service/ → Business logic
  ├── repo/    → Data access layer
  └── model/   → Domain models
pkg/           → Public libraries
```

### Dependency Injection
- Use interfaces for dependencies
- Inject dependencies through constructors
- Avoid global state
```go
type Service struct {
    repo Repository
    cache Cache
    logger Logger
}

func NewService(repo Repository, cache Cache, logger Logger) *Service {
    return &Service{
        repo: repo,
        cache: cache,
        logger: logger,
    }
}
```

## API Design

### RESTful Endpoints
- Use standard HTTP methods
- Resource-based URLs
- Consistent response format
```
GET    /api/v1/users          → List users
GET    /api/v1/users/{id}     → Get user
POST   /api/v1/users          → Create user
PUT    /api/v1/users/{id}     → Update user
DELETE /api/v1/users/{id}     → Delete user
```

### Request/Response Models
- Use separate DTOs for API
- Validate input thoroughly
- Return consistent error format
```go
type CreateUserRequest struct {
    Username string `json:"username" validate:"required,min=3,max=50"`
    Email    string `json:"email" validate:"required,email"`
    Password string `json:"password" validate:"required,min=8"`
}

type UserResponse struct {
    ID        string    `json:"id"`
    Username  string    `json:"username"`
    Email     string    `json:"email"`
    CreatedAt time.Time `json:"created_at"`
}

type ErrorResponse struct {
    Error   string            `json:"error"`
    Message string            `json:"message"`
    Details map[string]string `json:"details,omitempty"`
}
```

## Database Conventions

### Migrations
- Use sequential migration files
- Include up and down migrations
- Version control all migrations
```
migrations/
  ├── 001_create_users_table.up.sql
  ├── 001_create_users_table.down.sql
  ├── 002_add_email_index.up.sql
  └── 002_add_email_index.down.sql
```

### Query Patterns
- Use prepared statements
- Handle NULL values properly
- Use transactions for consistency
```go
func (r *UserRepo) Create(ctx context.Context, user *User) error {
    tx, err := r.db.BeginTx(ctx, nil)
    if err != nil {
        return err
    }
    defer tx.Rollback()
    
    query := `
        INSERT INTO users (id, username, email, created_at)
        VALUES ($1, $2, $3, $4)
    `
    
    _, err = tx.ExecContext(ctx, query, 
        user.ID, 
        user.Username, 
        user.Email, 
        user.CreatedAt,
    )
    if err != nil {
        return err
    }
    
    return tx.Commit()
}
```

## Configuration Management

### Environment Variables
- Use environment variables for config
- Provide sensible defaults
- Validate configuration on startup
```go
type Config struct {
    Port        string `env:"PORT" envDefault:"8080"`
    DatabaseURL string `env:"DATABASE_URL" envDefault:"postgres://localhost/app"`
    LogLevel    string `env:"LOG_LEVEL" envDefault:"info"`
    
    Redis struct {
        Host     string `env:"REDIS_HOST" envDefault:"localhost"`
        Port     string `env:"REDIS_PORT" envDefault:"6379"`
        Password string `env:"REDIS_PASSWORD"`
    }
}

func LoadConfig() (*Config, error) {
    cfg := &Config{}
    if err := env.Parse(cfg); err != nil {
        return nil, fmt.Errorf("failed to parse config: %w", err)
    }
    
    if err := cfg.Validate(); err != nil {
        return nil, fmt.Errorf("invalid config: %w", err)
    }
    
    return cfg, nil
}
```

## Logging

### Structured Logging
- Use structured logging (e.g., zap, logrus)
- Include request ID for tracing
- Log at appropriate levels
```go
logger.Info("user created",
    zap.String("user_id", user.ID),
    zap.String("username", user.Username),
    zap.String("request_id", requestID),
)

logger.Error("failed to save user",
    zap.Error(err),
    zap.String("user_id", user.ID),
    zap.String("request_id", requestID),
)
```

### Log Levels
- DEBUG: Detailed debugging information
- INFO: General informational messages
- WARN: Warning messages
- ERROR: Error messages
- FATAL: Fatal errors that cause shutdown

## Security

### Authentication
- Use secure token generation
- Implement proper session management
- Hash passwords with bcrypt
```go
func HashPassword(password string) (string, error) {
    bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
    return string(bytes), err
}

func CheckPassword(password, hash string) bool {
    err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
    return err == nil
}
```

### Input Validation
- Validate all user input
- Sanitize data before storage
- Use parameterized queries
```go
func ValidateEmail(email string) error {
    if len(email) > 255 {
        return errors.New("email too long")
    }
    
    emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
    if !emailRegex.MatchString(email) {
        return errors.New("invalid email format")
    }
    
    return nil
}
```

## Git Workflow

### Branch Naming
```
feature/user-authentication
bugfix/memory-leak
hotfix/security-patch
chore/update-dependencies
```

### Commit Messages
- Use conventional commits format
- Keep subject line under 50 characters
- Include body for complex changes
```
feat: add user authentication endpoint

- Implement JWT token generation
- Add password hashing with bcrypt
- Include rate limiting for login attempts
```

### Pull Request Process
1. Create feature branch from main
2. Write tests for new functionality
3. Ensure all tests pass
4. Update documentation
5. Request code review
6. Squash commits before merge

## Development Workflow

### Local Development
```bash
# Start development environment
make dev

# Run tests
make test

# Run linter
make lint

# Build application
make build

# Clean build artifacts
make clean
```

### Code Review Checklist
- [ ] Code follows style guide
- [ ] Tests cover new functionality
- [ ] Documentation is updated
- [ ] No sensitive data in code
- [ ] Error handling is appropriate
- [ ] Performance impact considered
- [ ] Security implications reviewed

## Monitoring and Observability

### Metrics
- Track key business metrics
- Monitor system performance
- Set up alerting thresholds
```go
var (
    requestDuration = prometheus.NewHistogramVec(
        prometheus.HistogramOpts{
            Name: "http_request_duration_seconds",
            Help: "Duration of HTTP requests in seconds",
        },
        []string{"method", "endpoint", "status"},
    )
)
```

### Health Checks
- Implement health endpoints
- Check critical dependencies
- Return appropriate status codes
```go
func HealthCheck(db *sql.DB, cache Cache) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        if err := db.PingContext(r.Context()); err != nil {
            w.WriteHeader(http.StatusServiceUnavailable)
            json.NewEncoder(w).Encode(map[string]string{
                "status": "unhealthy",
                "error":  "database unavailable",
            })
            return
        }
        
        w.WriteHeader(http.StatusOK)
        json.NewEncoder(w).Encode(map[string]string{
            "status": "healthy",
        })
    }
}
```