# Recommended Go Packages

This document outlines the key Go packages recommended for implementing the crypto trading bot, along with their purposes and benefits.

## Core Web Framework

### [github.com/gin-gonic/gin](https://github.com/gin-gonic/gin)
Gin is a high-performance HTTP web framework that serves as the foundation for our REST API.

**Benefits:**
- Extremely fast thanks to its use of [httprouter](https://github.com/julienschmidt/httprouter)
- Middleware support for authentication, logging, etc.
- Built-in request binding and validation
- Good error management

**Installation:**
```bash
go get -u github.com/gin-gonic/gin
```

**Example Usage:**
```go
func main() {
    router := gin.Default()
    
    router.GET("/ping", func(c *gin.Context) {
        c.JSON(200, gin.H{
            "message": "pong",
        })
    })
    
    router.Run(":8080")
}
```

## WebSocket Support

### [github.com/gorilla/websocket](https://github.com/gorilla/websocket)
This package provides a complete implementation of the WebSocket protocol, essential for real-time data from MEXC.

**Benefits:**
- Full RFC-6455 compliance
- Support for both client and server implementations
- Efficient message handling

**Installation:**
```bash
go get -u github.com/gorilla/websocket
```

**Example Usage with Gin:**
```go
var upgrader = websocket.Upgrader{
    ReadBufferSize:  1024,
    WriteBufferSize: 1024,
    CheckOrigin: func(r *http.Request) bool {
        return true // Allow all connections in development
    },
}

func handleWebSocket(c *gin.Context) {
    conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
    if err != nil {
        log.Printf("Failed to upgrade: %v", err)
        return
    }
    defer conn.Close()
    
    // WebSocket connection handling
    for {
        messageType, message, err := conn.ReadMessage()
        if err != nil {
            log.Printf("Error reading message: %v", err)
            break
        }
        
        // Echo the message back
        if err := conn.WriteMessage(messageType, message); err != nil {
            log.Printf("Error writing message: %v", err)
            break
        }
    }
}
```

## Database Access

### [github.com/jmoiron/sqlx](https://github.com/jmoiron/sqlx)
An extension of Go's standard `database/sql` package with additional features that make data access simpler.

**Benefits:**
- Simplified query building
- Struct mapping for cleaner code
- Transaction management
- Named parameter support

**Installation:**
```bash
go get -u github.com/jmoiron/sqlx
```

**Example Usage:**
```go
type Coin struct {
    ID            int64     `db:"id"`
    Symbol        string    `db:"symbol"`
    PurchasePrice float64   `db:"purchase_price"`
    Quantity      float64   `db:"quantity"`
    PurchasedAt   time.Time `db:"purchased_at"`
}

func GetCoinByID(db *sqlx.DB, id int64) (*Coin, error) {
    var coin Coin
    err := db.Get(&coin, "SELECT * FROM coins WHERE id = ?", id)
    if err != nil {
        return nil, err
    }
    return &coin, nil
}

func GetAllCoins(db *sqlx.DB) ([]Coin, error) {
    var coins []Coin
    err := db.Select(&coins, "SELECT * FROM coins ORDER BY purchased_at DESC")
    if err != nil {
        return nil, err
    }
    return coins, nil
}
```

### [github.com/mattn/go-sqlite3](https://github.com/mattn/go-sqlite3)
A SQLite driver for Go's `database/sql` package.

**Benefits:**
- Easy setup with no external server required
- Good performance for most use cases
- Full SQLite feature support
- ACID compliance

**Installation:**
```bash
go get -u github.com/mattn/go-sqlite3
```

**Example Usage:**
```go
import (
    "database/sql"
    "log"
    
    _ "github.com/mattn/go-sqlite3"
)

func InitDB(dbPath string) (*sql.DB, error) {
    db, err := sql.Open("sqlite3", dbPath)
    if err != nil {
        return nil, err
    }
    
    // Test the connection
    if err := db.Ping(); err != nil {
        return nil, err
    }
    
    // Create schema if needed
    if _, err := db.Exec(`
        CREATE TABLE IF NOT EXISTS coins (
            id INTEGER PRIMARY KEY AUTOINCREMENT,
            symbol TEXT NOT NULL,
            purchase_price REAL NOT NULL,
            quantity REAL NOT NULL,
            purchased_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
        );
    `); err != nil {
        return nil, err
    }
    
    return db, nil
}
```

## Configuration Management

### [github.com/spf13/viper](https://github.com/spf13/viper)
A complete configuration solution that supports various config formats and sources.

**Benefits:**
- Multiple config formats (JSON, YAML, TOML, etc.)
- Environment variable binding
- Hot reloading of configuration
- Command line flag binding

**Installation:**
```bash
go get -u github.com/spf13/viper
```

**Example Usage:**
```go
func InitConfig() {
    viper.SetConfigName("config")     // name of config file (without extension)
    viper.SetConfigType("yaml")       // type of config file
    viper.AddConfigPath(".")          // look for config in current directory
    viper.AddConfigPath("$HOME/.app") // look for config in home directory
    
    // Set defaults
    viper.SetDefault("server.port", 8080)
    viper.SetDefault("database.path", "./data.db")
    
    // Read environment variables prefixed with APP_
    viper.SetEnvPrefix("APP")
    viper.AutomaticEnv()
    
    // Read the config file
    if err := viper.ReadInConfig(); err != nil {
        if _, ok := err.(viper.ConfigFileNotFoundError); ok {
            log.Println("Config file not found, using defaults")
        } else {
            log.Fatalf("Error reading config file: %v", err)
        }
    }
}

func GetServerConfig() ServerConfig {
    return ServerConfig{
        Port: viper.GetInt("server.port"),
        Host: viper.GetString("server.host"),
    }
}
```

## Logging

### [github.com/rs/zerolog](https://github.com/rs/zerolog)
A fast and efficient structured logger for Go.

**Benefits:**
- Zero allocation JSON logging
- Highly customizable
- Extremely fast performance
- Contextual logging support

**Installation:**
```bash
go get -u github.com/rs/zerolog/log
```

**Example Usage:**
```go
import (
    "github.com/rs/zerolog"
    "github.com/rs/zerolog/log"
    "os"
    "time"
)

func InitLogger(level string) {
    // Parse log level
    logLevel, err := zerolog.ParseLevel(level)
    if err != nil {
        logLevel = zerolog.InfoLevel
    }
    
    // Set global log level
    zerolog.SetGlobalLevel(logLevel)
    
    // Configure logger output
    output := zerolog.ConsoleWriter{
        Out:        os.Stdout,
        TimeFormat: time.RFC3339,
    }
    
    // Set logger
    log.Logger = zerolog.New(output).
        With().
        Timestamp().
        Caller().
        Logger()
}

func ExampleLogging() {
    // Simple logging
    log.Info().Msg("Hello World!")
    
    // With fields
    log.Error().
        Str("service", "crypto-bot").
        Str("exchange", "MEXC").
        Err(errors.New("connection failed")).
        Msg("Failed to connect to exchange")
    
    // With context
    ctx := log.With().Str("component", "trade-service").Logger().WithContext(context.Background())
    log.Ctx(ctx).Debug().Msg("Processing trade")
}
```

## HTTP Client

### [github.com/go-resty/resty/v2](https://github.com/go-resty/resty)
A simple HTTP and REST client for Go with elegant DSL.

**Benefits:**
- Simple and intuitive API
- Middleware support
- Automatic marshaling and unmarshaling
- Request retries and timeouts

**Installation:**
```bash
go get -u github.com/go-resty/resty/v2
```

**Example Usage for MEXC API:**
```go
import (
    "crypto/hmac"
    "crypto/sha256"
    "encoding/hex"
    "github.com/go-resty/resty/v2"
    "strconv"
    "time"
)

type MexcClient struct {
    client    *resty.Client
    baseURL   string
    apiKey    string
    apiSecret string
}

func NewMexcClient(apiKey, apiSecret string) *MexcClient {
    client := resty.New()
    client.SetTimeout(10 * time.Second)
    client.SetHeader("Content-Type", "application/json")
    
    return &MexcClient{
        client:    client,
        baseURL:   "https://api.mexc.com",
        apiKey:    apiKey,
        apiSecret: apiSecret,
    }
}

func (c *MexcClient) GetTicker(symbol string) (map[string]interface{}, error) {
    resp, err := c.client.R().
        SetQueryParam("symbol", symbol).
        SetResult(map[string]interface{}{}).
        Get(c.baseURL + "/api/v3/ticker/price")
    
    if err != nil {
        return nil, err
    }
    
    return resp.Result().(map[string]interface{}), nil
}

func (c *MexcClient) CreateOrder(symbol, side, orderType string, quantity float64) (map[string]interface{}, error) {
    timestamp := strconv.FormatInt(time.Now().UnixMilli(), 10)
    
    // Create signature
    params := "symbol=" + symbol +
        "&side=" + side +
        "&type=" + orderType +
        "&quantity=" + strconv.FormatFloat(quantity, 'f', -1, 64) +
        "&timestamp=" + timestamp
    
    h := hmac.New(sha256.New, []byte(c.apiSecret))
    h.Write([]byte(params))
    signature := hex.EncodeToString(h.Sum(nil))
    
    resp, err := c.client.R().
        SetQueryParams(map[string]string{
            "symbol":    symbol,
            "side":      side,
            "type":      orderType,
            "quantity":  strconv.FormatFloat(quantity, 'f', -1, 64),
            "timestamp": timestamp,
            "signature": signature,
        }).
        SetHeader("X-MEXC-APIKEY", c.apiKey).
        SetResult(map[string]interface{}{}).
        Post(c.baseURL + "/api/v3/order")
    
    if err != nil {
        return nil, err
    }
    
    return resp.Result().(map[string]interface{}), nil
}
```

## Testing

### [github.com/stretchr/testify](https://github.com/stretchr/testify)
A toolkit with common assertions and mocks for testing Go code.

**Benefits:**
- Rich assertion library
- Suite testing support
- Mock generation and verification
- Easy-to-read test output

**Installation:**
```bash
go get -u github.com/stretchr/testify
```

**Example Usage:**
```go
import (
    "testing"
    "github.com/stretchr/testify/assert"
    "github.com/stretchr/testify/mock"
)

// Unit test with assertions
func TestCalculateProfitLoss(t *testing.T) {
    // Test data
    coin := &Coin{
        Symbol:        "BTCUSDT",
        PurchasePrice: 50000.0,
        Quantity:      0.1,
    }
    currentPrice := 55000.0
    
    // Call function
    profit, percentage := CalculateProfitLoss(coin, currentPrice)
    
    // Assertions
    assert.Equal(t, 500.0, profit, "Profit calculation incorrect")
    assert.Equal(t, 10.0, percentage, "Profit percentage calculation incorrect")
}

// Mock service for testing
type MockMarketService struct {
    mock.Mock
}

func (m *MockMarketService) GetCurrentPrice(symbol string) (float64, error) {
    args := m.Called(symbol)
    return args.Get(0).(float64), args.Error(1)
}

// Testing with mocks
func TestTradeService_ExecutePurchase(t *testing.T) {
    // Create mocks
    mockMarket := new(MockMarketService)
    mockRepo := new(MockCoinRepository)
    
    // Set expectations
    mockMarket.On("GetCurrentPrice", "BTCUSDT").Return(50000.0, nil)
    mockRepo.On("Create", mock.Anything).Return(int64(1), nil)
    
    // Create service with mocks
    service := NewTradeService(mockRepo, mockMarket)
    
    // Call method
    result, err := service.ExecutePurchase("BTCUSDT", 0.1)
    
    // Assertions
    assert.NoError(t, err)
    assert.NotNil(t, result)
    assert.Equal(t, "BTCUSDT", result.Symbol)
    assert.Equal(t, 50000.0, result.PurchasePrice)
    assert.Equal(t, 0.1, result.Quantity)
    
    // Verify mock expectations
    mockMarket.AssertExpectations(t)
    mockRepo.AssertExpectations(t)
}
```

## Extra Utilities

### [github.com/google/uuid](https://github.com/google/uuid)
For generating UUIDs that can be used as correlation IDs in logs and requests.

**Installation:**
```bash
go get -u github.com/google/uuid
```

### [github.com/shopspring/decimal](https://github.com/shopspring/decimal)
For safe decimal math when dealing with financial calculations.

**Installation:**
```bash
go get -u github.com/shopspring/decimal
```

### [github.com/go-playground/validator/v10](https://github.com/go-playground/validator)
For struct and field validation in request models.

**Installation:**
```bash
go get -u github.com/go-playground/validator/v10
```

## Conclusion

These recommended packages provide a strong foundation for building a robust crypto trading bot. They're well-maintained, widely used in the Go ecosystem, and solve many common needs in your application. The examples provided should give you a good starting point for implementation.

For a more detailed example of how these packages work together in a complete application, see the implementation guides in this documentation series.
