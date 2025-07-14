# EZUtil

A comprehensive Go utility library that provides common functionality for web applications built with Gin, GORM, and other popular Go frameworks.

## Features

### 🌐 HTTP & Web Utilities
- **Gin Utilities**: Parameter extraction, request binding, and response helpers
- **Gin Middlewares**: Common middleware implementations for web applications
- **Gin Routing**: Simplified routing utilities and helpers
- **HTTP Utils**: General HTTP-related utility functions

### 🗄️ Database & ORM
- **GORM Utilities**: Database connection management and configuration
- **GORM Scopes**: Reusable query scopes for common database operations
- **GORM Transactor**: Transaction management utilities
- **SQL Utils**: General SQL utility functions

### 🔧 Configuration Management
- **Environment-based Configuration**: Load configuration from environment variables
- **Database Configuration**: Automatic database connection setup (MySQL/PostgreSQL)
- **Application Configuration**: Centralized app settings management

### 🔐 Authentication & Security
- **JWT Service**: Token creation and verification
- **Authentication Utilities**: User authentication helpers

### 🛠️ General Utilities
- **String Utils**: String parsing, random string generation, and manipulation
- **Time Utils**: Date/time formatting and manipulation functions
- **Slice Utils**: Functional programming utilities for slice operations
- **Error Handling**: Structured error types with HTTP context
- **Template Utils**: Template rendering utilities (Templ integration)

## Installation

```bash
go get github.com/itsLeonB/ezutil
```

## Quick Start

### Basic Usage

```go
package main

import (
    "github.com/gin-gonic/gin"
    "github.com/itsLeonB/ezutil"
)

func main() {
    // Load configuration from environment
    config, err := ezutil.LoadConfig()
    if err != nil {
        panic(err)
    }

    // Create Gin router
    r := gin.Default()
    
    // Use ezutil helpers
    r.GET("/user/:id", func(c *gin.Context) {
        userID, exists, err := ezutil.GetPathParam[int](c, "id")
        if err != nil {
            c.JSON(400, gin.H{"error": "Invalid user ID"})
            return
        }
        if !exists {
            c.JSON(400, gin.H{"error": "User ID required"})
            return
        }
        
        c.JSON(200, gin.H{"user_id": userID})
    })
    
    r.Run(":8080")
}
```

### Configuration Setup

Create a `.env` file or set environment variables:

```env
# Application Configuration
APP_NAME=MyApp
APP_ENV=development
APP_PORT=8080

# Database Configuration
DB_DRIVER=postgres
DB_HOST=localhost
DB_PORT=5432
DB_NAME=myapp
DB_USER=username
DB_PASSWORD=password

# JWT Configuration
JWT_SECRET=your-secret-key
JWT_EXPIRES_IN=24h
```

### Using Utilities

#### String Parsing
```go
// Parse string to different types
userID, err := ezutil.Parse[int]("123")
isActive, err := ezutil.Parse[bool]("true")
uuid, err := ezutil.Parse[uuid.UUID]("550e8400-e29b-41d4-a716-446655440000")

// Generate random strings
randomStr, err := ezutil.GenerateRandomString(32)
```

#### Slice Operations
```go
// Transform slices functionally
numbers := []int{1, 2, 3, 4, 5}
doubled := ezutil.MapSlice(numbers, func(n int) int {
    return n * 2
})

// Map with error handling
strings := []string{"1", "2", "invalid", "4"}
numbers, err := ezutil.MapSliceWithError(strings, func(s string) (int, error) {
    return ezutil.Parse[int](s)
})
```

#### Time Utilities
```go
// Get start and end of day
startOfDay, err := ezutil.GetStartOfDay(2024, 1, 15)
endOfDay, err := ezutil.GetEndOfDay(2024, 1, 15)

// Format time with null handling
formatted := ezutil.FormatTimeNullable(time.Now(), "2006-01-02 15:04:05")
```

#### Error Handling
```go
// Create structured application errors
appErr := ezutil.NewAppError(
    "VALIDATION_ERROR",
    "Invalid input provided",
    http.StatusBadRequest,
    map[string]string{"field": "email", "issue": "invalid format"},
)

// Use in Gin handlers
if err != nil {
    ezutil.HandleError(c, appErr)
    return
}
```

## Dependencies

This library builds upon several excellent Go packages:

- **[Gin](https://github.com/gin-gonic/gin)**: HTTP web framework
- **[GORM](https://gorm.io/)**: ORM library for database operations
- **[Eris](https://github.com/rotisserie/eris)**: Error handling and stack traces
- **[JWT](https://github.com/golang-jwt/jwt)**: JSON Web Token implementation
- **[Templ](https://github.com/a-h/templ)**: Template engine integration
- **[Envconfig](https://github.com/kelseyhightower/envconfig)**: Environment variable configuration
- **[UUID](https://github.com/google/uuid)**: UUID generation and parsing

## Project Structure

```
ezutil/
├── config/                 # Configuration structures
├── internal/              # Internal utilities
├── .github/               # GitHub workflows and templates
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
└── go.mod                # Go module definition
```

## Testing

EZUtil includes a comprehensive test suite covering all exported functions and methods. The tests are located in the `test/` directory and use the `ezutil_test` package.

### Running Tests

Use the provided Makefile commands to run tests:

```bash
# Show all available commands
make help

# Run all tests (quick)
make test

# Run tests with verbose output
make test-verbose

# Run tests with coverage report
make test-coverage

# Clean test cache and run tests
make test-clean
```

### Test Coverage

The test suite provides comprehensive coverage including:
- ✅ All exported functions and methods
- ✅ Happy path scenarios and error conditions
- ✅ Edge cases and input validation
- ✅ Database operations (using in-memory SQLite)
- ✅ HTTP request/response handling
- ✅ JWT token operations
- ✅ Configuration loading and validation

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request. For major changes, please open an issue first to discuss what you would like to change.

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add some amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## Author

**Ellion Blessan** - [itsLeonB](https://github.com/itsLeonB)

---

*EZUtil - Making Go web development easier, one utility at a time.*
