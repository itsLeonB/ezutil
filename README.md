# EZUtil

[![CI](https://github.com/itsLeonB/ezutil/workflows/CI/badge.svg)](https://github.com/itsLeonB/ezutil/actions)
[![Tests](https://github.com/itsLeonB/ezutil/workflows/Tests/badge.svg)](https://github.com/itsLeonB/ezutil/actions)
[![Go Report Card](https://goreportcard.com/badge/github.com/itsLeonB/ezutil)](https://goreportcard.com/report/github.com/itsLeonB/ezutil)
[![codecov](https://codecov.io/gh/itsLeonB/ezutil/branch/main/graph/badge.svg)](https://codecov.io/gh/itsLeonB/ezutil)
[![Go Reference](https://pkg.go.dev/badge/github.com/itsLeonB/ezutil.svg)](https://pkg.go.dev/github.com/itsLeonB/ezutil)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

A comprehensive, production-ready Go utility library that provides common functionality for web applications built with Gin, GORM, and other popular Go frameworks. EZUtil is designed to accelerate development by providing well-tested, reusable components for modern Go web applications.

## 🚀 Features

### 🌐 HTTP & Web Utilities
- **Gin Parameter Extraction**: Type-safe parameter parsing with `GetPathParam[T]`, `GetQueryParam[T]`
- **Request Binding**: Simplified JSON/form data binding with validation
- **Response Helpers**: Standardized JSON response utilities
- **Middleware Support**: Common middleware implementations for web applications
- **Routing Utilities**: Simplified routing helpers and patterns

### 🗄️ Database & ORM
- **GORM Integration**: Seamless database connection management and configuration
- **Query Scopes**: Reusable, composable query scopes for common database operations
- **Transaction Management**: Robust transaction utilities with nested transaction support
- **Multi-Database Support**: MySQL and PostgreSQL drivers with automatic configuration
- **Connection Pooling**: Optimized database connection management

### 🔧 Configuration Management
- **Environment-Based Config**: Automatic configuration loading from environment variables
- **Type-Safe Parsing**: Built-in validation and type conversion for configuration values
- **Database Auto-Configuration**: Automatic database connection setup with fallback defaults
- **Application Settings**: Centralized management of app-level configuration
- **Flexible Loading**: Support for loading configuration with or without database dependency

### 🔐 Authentication & Security
- **JWT Service**: Complete JWT token creation, verification, and management
- **Secure Token Handling**: Built-in token expiration and refresh capabilities
- **Authentication Middleware**: Ready-to-use authentication middleware for Gin
- **Security Best Practices**: Implements industry-standard security patterns

### 🛠️ General Utilities
- **Type-Safe Parsing**: Generic `Parse[T]` function for string-to-type conversion
- **String Utilities**: Random string generation, manipulation, and validation
- **Time Management**: Date/time formatting, manipulation, and timezone handling
- **Slice Operations**: Functional programming utilities (`MapSlice`, `MapSliceWithError`)
- **Error Handling**: Structured error types with HTTP context and stack traces
- **Template Integration**: Seamless integration with Templ template engine
- **UUID Support**: UUID generation, parsing, and validation utilities

## 📦 Installation

```bash
go get github.com/itsLeonB/ezutil
```

**Requirements:**
- Go 1.23 or higher
- Compatible with Go 1.23, and 1.24

## 🏃 Quick Start

### Basic Web Application

```go
package main

import (
    "log"
    "time"
    
    "github.com/gin-gonic/gin"
    "github.com/itsLeonB/ezutil"
)

func main() {
    // Load configuration from environment
    defaults := ezutil.Config{
        App: &ezutil.App{
            Env:        "development",
            Port:       "8080",
            Timeout:    30 * time.Second,
            ClientUrls: []string{"http://localhost:8080"},
            Timezone:   "UTC",
        },
        Auth: &ezutil.Auth{
            SecretKey:      "your-secret-key",
            TokenDuration:  24 * time.Hour,
            CookieDuration: 7 * 24 * time.Hour,
            Issuer:         "your-app",
            URL:            "http://localhost:8080",
        },
    }
    
    config := ezutil.LoadConfig(defaults)
    
    // Create Gin router
    r := gin.Default()
    
    // Use EZUtil helpers for type-safe parameter extraction
    r.GET("/user/:id", func(c *gin.Context) {
        userID, exists, err := ezutil.GetPathParam[int](c, "id")
        if err != nil {
            c.JSON(400, gin.H{"error": "Invalid user ID format"})
            return
        }
        if !exists {
            c.JSON(400, gin.H{"error": "User ID is required"})
            return
        }
        
        c.JSON(200, gin.H{"user_id": userID, "message": "User found"})
    })
    
    // Start server
    log.Printf("Server starting on port %s", config.App.Port)
    r.Run(":" + config.App.Port)
}
```

### Configuration Setup

Create a `.env` file or set environment variables:

```env
# Application Configuration
APP_ENV=production
APP_PORT=8080
APP_TIMEOUT=30s
APP_TIMEZONE=UTC
APP_CLIENTURLS=http://localhost:8080,https://yourdomain.com

# Database Configuration
SQLDB_HOST=localhost
SQLDB_PORT=5432
SQLDB_NAME=myapp
SQLDB_USER=username
SQLDB_PASSWORD=password
SQLDB_DRIVER=postgres

# Authentication Configuration
AUTH_SECRETKEY=your-super-secret-jwt-key
AUTH_TOKENDURATION=24h
AUTH_COOKIEDURATION=168h
AUTH_ISSUER=myapp
AUTH_URL=https://yourdomain.com
```

### Advanced Usage Examples

#### Database Operations with Transactions

```go
// Initialize transactor
transactor := ezutil.NewTransactor(db)

// Use transactions with automatic rollback on error
err := transactor.WithinTransaction(ctx, func(ctx context.Context) error {
    tx, err := ezutil.GetTxFromContext(ctx)
    if err != nil {
        return err
    }
    
    // Perform database operations within transaction
    user := User{Name: "John Doe", Email: "john@example.com"}
    if err := tx.Create(&user).Error; err != nil {
        return err // Transaction will be rolled back automatically
    }
    
    // Nested transactions are supported
    return transactor.WithinTransaction(ctx, func(ctx context.Context) error {
        innerTx, err := ezutil.GetTxFromContext(ctx)
        if err != nil {
            return err
        }
        
        profile := Profile{UserID: user.ID, Bio: "Software Developer"}
        return innerTx.Create(&profile).Error
    })
})

if err != nil {
    log.Printf("Transaction failed: %v", err)
}
```

#### Type-Safe String Parsing

```go
// Parse various types from strings
userID, err := ezutil.Parse[int]("123")
if err != nil {
    log.Printf("Invalid user ID: %v", err)
}

isActive, err := ezutil.Parse[bool]("true")
price, err := ezutil.Parse[float64]("29.99")
uuid, err := ezutil.Parse[uuid.UUID]("550e8400-e29b-41d4-a716-446655440000")

// Generate secure random strings
apiKey, err := ezutil.GenerateRandomString(32)
if err != nil {
    log.Printf("Failed to generate API key: %v", err)
}
```

#### Functional Slice Operations

```go
// Transform slices with type safety
numbers := []int{1, 2, 3, 4, 5}
doubled := ezutil.MapSlice(numbers, func(n int) int {
    return n * 2
})
// Result: [2, 4, 6, 8, 10]

// Handle errors during transformation
strings := []string{"1", "2", "invalid", "4"}
numbers, err := ezutil.MapSliceWithError(strings, func(s string) (int, error) {
    return ezutil.Parse[int](s)
})
if err != nil {
    log.Printf("Conversion failed: %v", err)
}
```

#### Advanced Error Handling

```go
// Create structured application errors
appErr := ezutil.NewAppError(
    "VALIDATION_ERROR",
    "Invalid input provided",
    http.StatusBadRequest,
    map[string]string{
        "field": "email",
        "issue": "invalid format",
    },
)

// Use in Gin handlers with automatic error response
r.POST("/users", func(c *gin.Context) {
    var user User
    if err := c.ShouldBindJSON(&user); err != nil {
        ezutil.HandleError(c, ezutil.NewAppError(
            "BIND_ERROR",
            "Invalid JSON payload",
            http.StatusBadRequest,
            nil,
        ))
        return
    }
    
    // Process user...
})
```

#### JWT Authentication

```go
// Create JWT service
jwtService := ezutil.NewJWTService(config.Auth.SecretKey, config.Auth.Issuer)

// Generate tokens
claims := map[string]interface{}{
    "user_id": 123,
    "role":    "admin",
}

token, err := jwtService.GenerateToken(claims, config.Auth.TokenDuration)
if err != nil {
    log.Printf("Token generation failed: %v", err)
}

// Verify tokens
parsedClaims, err := jwtService.VerifyToken(token)
if err != nil {
    log.Printf("Token verification failed: %v", err)
}
```

## 🏗️ Project Structure

```
ezutil/
├── .github/
│   └── workflows/           # CI/CD workflows
│       ├── ci.yml          # Main CI pipeline
│       ├── test.yml        # Extended testing with security scans
│       └── lint.yml        # Code linting
├── test/                   # Comprehensive test suite
│   ├── *_test.go          # Test files for each module
│   ├── go.mod             # Test module dependencies
│   └── go.sum             # Test dependency checksums
├── config/                 # Configuration structures
├── internal/              # Internal utilities
├── config_loader.go       # Environment configuration loading
├── errors.go              # Error handling utilities
├── gin_*.go              # Gin framework utilities
├── gorm_*.go             # GORM database utilities
├── http_utils.go         # HTTP utilities
├── services.go           # Service layer utilities (JWT, etc.)
├── slice_utils.go        # Slice manipulation utilities
├── sql_utils.go          # SQL utilities
├── string_utils.go       # String manipulation utilities
├── templ_utils.go        # Template utilities
├── time_utils.go         # Time/date utilities
├── uuid_utils.go         # UUID utilities
├── Makefile              # Build and test automation
├── go.mod                # Go module definition
└── README.md             # This file
```

## 🧪 Testing

EZUtil includes a comprehensive test suite with over 200 individual test cases covering all exported functions and methods. The tests are organized in a separate `test/` directory using the `ezutil_test` package for proper isolation.

### Test Coverage

The test suite provides comprehensive coverage including:
- ✅ **All exported functions and methods** - 100% coverage of public API
- ✅ **Happy path scenarios** - Normal operation testing
- ✅ **Error conditions** - Comprehensive error handling validation
- ✅ **Edge cases** - Boundary condition testing
- ✅ **Database operations** - Using in-memory SQLite for isolation
- ✅ **HTTP request/response handling** - Complete web layer testing
- ✅ **JWT token operations** - Authentication flow testing
- ✅ **Configuration loading** - Environment variable processing
- ✅ **Transaction management** - Database transaction testing
- ✅ **Type safety** - Generic function validation

### Running Tests

Use the provided Makefile commands for various testing scenarios:

```bash
# Show all available commands
make help

# Run all tests (quick)
make test

# Run tests with verbose output
make test-verbose

# Run tests with coverage report
make test-coverage

# Generate HTML coverage report
make test-coverage-html

# Clean test cache and run fresh tests
make test-clean
```

### Continuous Integration

The project uses GitHub Actions for comprehensive CI/CD:

#### **Main CI Pipeline** (`ci.yml`)
- **Multi-version testing**: Go 1.23, 1.24
- **Automated testing**: Full test suite execution
- **Coverage reporting**: Automatic upload to Codecov
- **Build verification**: Cross-version compatibility

#### **Extended Testing** (`test.yml`)
- **Comprehensive testing**: All test scenarios
- **Security scanning**: Gosec static analysis
- **Dependency verification**: Module integrity checks
- **SARIF reporting**: Security findings integration

#### **Code Quality** (`lint.yml`)
- **Static analysis**: golangci-lint integration
- **Code formatting**: Automated style checking
- **Best practices**: Go idiom enforcement

### Test Organization

```
test/
├── config_loader_test.go    # Configuration loading tests
├── errors_test.go           # Error handling tests
├── gin_utils_test.go        # Gin utilities tests
├── gorm_scopes_test.go      # Database scope tests
├── gorm_transactor_test.go  # Transaction management tests
├── http_utils_test.go       # HTTP utility tests
├── services_test.go         # Service layer tests (JWT, etc.)
├── slice_utils_test.go      # Slice operation tests
├── sql_utils_test.go        # SQL utility tests
├── string_utils_test.go     # String manipulation tests
├── templ_utils_test.go      # Template utility tests
├── time_utils_test.go       # Time/date utility tests
└── uuid_utils_test.go       # UUID utility tests
```

## 📚 Dependencies

EZUtil builds upon several excellent Go packages:

### Core Dependencies
- **[Gin](https://github.com/gin-gonic/gin)** `v1.9.1` - HTTP web framework
- **[GORM](https://gorm.io/)** `v1.25.5` - ORM library for database operations
- **[Eris](https://github.com/rotisserie/eris)** `v0.8.1` - Error handling and stack traces
- **[JWT](https://github.com/golang-jwt/jwt)** `v5.2.0` - JSON Web Token implementation

### Database Drivers
- **[MySQL Driver](https://github.com/go-sql-driver/mysql)** - MySQL database support
- **[PostgreSQL Driver](https://github.com/lib/pq)** - PostgreSQL database support

### Utility Libraries
- **[Templ](https://github.com/a-h/templ)** `v0.2.543` - Template engine integration
- **[Envconfig](https://github.com/kelseyhightower/envconfig)** `v1.4.0` - Environment variable configuration
- **[UUID](https://github.com/google/uuid)** `v1.4.0` - UUID generation and parsing

### Development Dependencies
- **[Testify](https://github.com/stretchr/testify)** `v1.8.4` - Testing assertions and mocks
- **[SQLite Driver](https://github.com/mattn/go-sqlite3)** - In-memory testing database

## 🔧 Configuration Reference

### Application Configuration (`APP_*`)

| Variable | Type | Default | Description |
|----------|------|---------|-------------|
| `APP_ENV` | string | `development` | Application environment |
| `APP_PORT` | string | `3000` | Server port number |
| `APP_TIMEOUT` | duration | `10s` | Request timeout |
| `APP_CLIENTURLS` | []string | `["http://localhost:3000"]` | Allowed client URLs |
| `APP_TIMEZONE` | string | `America/New_York` | Application timezone |

### Database Configuration (`SQLDB_*`)

| Variable | Type | Required | Description |
|----------|------|----------|-------------|
| `SQLDB_HOST` | string | ✅ | Database host |
| `SQLDB_PORT` | string | ✅ | Database port |
| `SQLDB_NAME` | string | ✅ | Database name |
| `SQLDB_USER` | string | ✅ | Database username |
| `SQLDB_PASSWORD` | string | ✅ | Database password |
| `SQLDB_DRIVER` | string | ✅ | Database driver (`mysql` or `postgres`) |

### Authentication Configuration (`AUTH_*`)

| Variable | Type | Default | Description |
|----------|------|---------|-------------|
| `AUTH_SECRETKEY` | string | `default-secret` | JWT signing key |
| `AUTH_TOKENDURATION` | duration | `30m` | JWT token lifetime |
| `AUTH_COOKIEDURATION` | duration | `12h` | Cookie lifetime |
| `AUTH_ISSUER` | string | `default-issuer` | JWT issuer |
| `AUTH_URL` | string | `http://localhost:3000` | Authentication service URL |

## 🚀 Performance & Best Practices

### Database Optimization
- **Connection Pooling**: Automatic connection pool management
- **Transaction Efficiency**: Nested transaction support with proper rollback
- **Query Optimization**: Reusable scopes for common query patterns
- **Type Safety**: Compile-time type checking for database operations

### Security Features
- **JWT Security**: Secure token generation with configurable expiration
- **Input Validation**: Built-in parameter validation and sanitization
- **Error Handling**: Structured errors without sensitive information leakage
- **HTTPS Support**: Ready for production HTTPS deployment

### Development Experience
- **Type Safety**: Generic functions for compile-time type checking
- **Error Context**: Rich error information with stack traces
- **Testing Support**: Comprehensive test utilities and mocks
- **Documentation**: Extensive inline documentation and examples

## 🤝 Contributing

Contributions are welcome! Please feel free to submit a Pull Request. For major changes, please open an issue first to discuss what you would like to change.

### Development Setup

1. **Fork and Clone**
   ```bash
   git clone https://github.com/yourusername/ezutil.git
   cd ezutil
   ```

2. **Install Dependencies**
   ```bash
   go mod download
   cd test && go mod download
   ```

3. **Run Tests**
   ```bash
   make test-verbose
   ```

4. **Run Linting**
   ```bash
   make lint
   ```

### Contribution Guidelines

1. **Code Quality**: Ensure all tests pass and maintain test coverage
2. **Documentation**: Update documentation for new features
3. **Commit Messages**: Use clear, descriptive commit messages
4. **Pull Requests**: Include description of changes and test results

### Development Workflow

1. Create your feature branch (`git checkout -b feature/amazing-feature`)
2. Make your changes and add tests
3. Ensure all tests pass (`make test`)
4. Run linting (`make lint`)
5. Commit your changes (`git commit -m 'Add some amazing feature'`)
6. Push to the branch (`git push origin feature/amazing-feature`)
7. Open a Pull Request

## 📄 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## 👨‍💻 Author

**Ellion Blessan** - [itsLeonB](https://github.com/itsLeonB)

## 🙏 Acknowledgments

- The Go community for excellent libraries and tools
- Contributors who help improve this project
- Users who provide feedback and bug reports

---

**EZUtil** - Making Go web development easier, one utility at a time. 🚀

*Built with ❤️ for the Go community*
