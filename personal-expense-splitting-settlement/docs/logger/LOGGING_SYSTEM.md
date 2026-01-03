# Server Logging System - Complete Documentation

## 📋 Overview

The Personal Expense Splitting application implements a comprehensive logging system using **Zap logger** with dual output (console + file) and detailed request/response tracking.

---

## 🎯 Features

### 1. **Dual Output Logging**
- **Console Output**: Human-readable format for development
- **File Output**: JSON format in `logs/app.log` for production analysis

### 2. **Request Tracking**
- Unique Request ID for each HTTP request
- Request/Response body logging
- Duration tracking
- IP address and User-Agent capture
- Error tracking and aggregation

### 3. **Log Levels**
- **INFO**: Successful operations (2xx responses)
- **WARN**: Client errors (4xx responses)
- **ERROR**: Server errors (5xx responses)
- **DEBUG**: Detailed debugging information

### 4. **Automatic Log Rotation**
- Logs automatically rotate based on size
- Old logs are preserved for analysis
- Configurable retention policy

---

## 📁 Log Files Location

```
project-root/
└── logs/
    ├── app.log              # Current log file (JSON format)
    ├── app.log.2026-01-03   # Rotated logs (if rotation enabled)
    └── ...
```

**Note**: The `logs/` directory is automatically created on first run and is excluded from Git (via `.gitignore`).

---

## 🔧 Configuration

### Environment Variables

```bash
# Set log level (debug, info, warn, error)
ENVIRONMENT=development  # Uses debug level
ENVIRONMENT=production   # Uses info level
```

### Logger Initialization

Located in `pkg/logger/logger.go`:

```go
// InitLogger initializes the Zap logger with file and console output
func InitLogger() (*zap.SugaredLogger, error)
```

**Features**:
- Creates `logs/` directory if it doesn't exist
- Configures JSON encoder for file output
- Configures console encoder for stdout
- Sets log level based on environment
- Enables caller information (file and line number)

---

## 📊 Log Entry Structure

### Request Log Entry (JSON Format)

```json
{
  "level": "info",
  "ts": "2026-01-03T12:45:30.123+0530",
  "caller": "middleware/logger_middleware.go:45",
  "msg": "Incoming request",
  "request_id": "a1b2c3d4-e5f6-7890-1234-567890abcdef",
  "method": "POST",
  "path": "/api/v1/friends/request",
  "query": "",
  "ip": "::1",
  "user_agent": "curl/8.14.1",
  "content_type": "application/json",
  "request_body": "{\"friend_email\":\"bob@test.com\"}"
}
```

### Response Log Entry (JSON Format)

```json
{
  "level": "info",
  "ts": "2026-01-03T12:45:30.456+0530",
  "caller": "middleware/logger_middleware.go:75",
  "msg": "Request completed",
  "request_id": "a1b2c3d4-e5f6-7890-1234-567890abcdef",
  "method": "POST",
  "path": "/api/v1/friends/request",
  "status": 201,
  "duration_ms": 333,
  "duration": "333.456789ms",
  "response_size": 78,
  "response_body": "{\"success\":true,\"message\":\"Friend request sent successfully\",\"data\":null}",
  "ip": "::1"
}
```

### Error Log Entry

```json
{
  "level": "error",
  "ts": "2026-01-03T12:45:30.789+0530",
  "caller": "middleware/logger_middleware.go:90",
  "msg": "Request error",
  "request_id": "a1b2c3d4-e5f6-7890-1234-567890abcdef",
  "error": "user not found",
  "type": "private"
}
```

---

## 🛠️ Implementation Details

### 1. Logger Middleware

**File**: `internal/middleware/logger_middleware.go`

**Key Functions**:
- Generates unique Request ID using UUID
- Captures request body (and restores it for handlers)
- Wraps response writer to capture response
- Logs incoming request details
- Logs outgoing response details
- Tracks request duration
- Handles error logging

**Usage**:
```go
router.Use(middleware.LoggerMiddleware(sugar))
```

### 2. Request ID Propagation

The Request ID is stored in Gin context and can be accessed in handlers:

```go
func (h *Handler) SomeEndpoint(c *gin.Context) {
    requestID, exists := c.Get("request_id")
    if exists {
        // Use request ID for tracking
    }
}
```

### 3. Response Body Capture

Custom `responseWriter` implementation:
- Extends `gin.ResponseWriter`
- Captures response body in buffer
- Limits response body logging to 1000 characters (prevents huge logs)
- Maintains original response behavior

---

## 📝 Log Fields Reference

### Request Fields

| Field | Type | Description | Example |
|-------|------|-------------|---------|
| `request_id` | UUID | Unique identifier for request | `a1b2c3d4-e5f6-7890-1234-567890abcdef` |
| `method` | String | HTTP method | `POST`, `GET`, `PUT`, `DELETE` |
| `path` | String | Request path | `/api/v1/friends/request` |
| `query` | String | Query parameters | `?page=1&limit=10` |
| `ip` | String | Client IP address | `::1`, `192.168.1.1` |
| `user_agent` | String | Client user agent | `curl/8.14.1`, `Mozilla/5.0...` |
| `content_type` | String | Request content type | `application/json` |
| `request_body` | String | Request payload | `{"email":"test@test.com"}` |

### Response Fields

| Field | Type | Description | Example |
|-------|------|-------------|---------|
| `status` | Integer | HTTP status code | `200`, `201`, `400`, `500` |
| `duration_ms` | Integer | Request duration in milliseconds | `333` |
| `duration` | String | Formatted duration | `333.456789ms` |
| `response_size` | Integer | Response size in bytes | `78` |
| `response_body` | String | Response payload (max 1000 chars) | `{"success":true,...}` |

### Error Fields

| Field | Type | Description | Example |
|-------|------|-------------|---------|
| `error` | String | Error message | `user not found` |
| `type` | String | Error type | `private`, `public` |

---

## 🔍 Querying Logs

### Using jq (JSON processor)

**Filter by Request ID**:
```bash
cat logs/app.log | jq 'select(.request_id == "a1b2c3d4-e5f6-7890-1234-567890abcdef")'
```

**Filter by Status Code**:
```bash
cat logs/app.log | jq 'select(.status >= 400)'
```

**Filter by Path**:
```bash
cat logs/app.log | jq 'select(.path == "/api/v1/friends/request")'
```

**Get Average Response Time**:
```bash
cat logs/app.log | jq -s 'map(select(.duration_ms)) | map(.duration_ms) | add / length'
```

**Count Errors**:
```bash
cat logs/app.log | jq 'select(.level == "error")' | wc -l
```

**Find Slow Requests (> 1 second)**:
```bash
cat logs/app.log | jq 'select(.duration_ms > 1000)'
```

---

## 📈 Monitoring Examples

### Real-time Log Monitoring

**Watch all logs**:
```bash
tail -f logs/app.log | jq .
```

**Watch errors only**:
```bash
tail -f logs/app.log | jq 'select(.level == "error")'
```

**Watch specific endpoint**:
```bash
tail -f logs/app.log | jq 'select(.path == "/api/v1/friends/request")'
```

### Log Analysis

**Top 10 slowest endpoints**:
```bash
cat logs/app.log | jq -s 'map(select(.duration_ms)) | sort_by(.duration_ms) | reverse | .[0:10]'
```

**Requests by status code**:
```bash
cat logs/app.log | jq -s 'group_by(.status) | map({status: .[0].status, count: length})'
```

**Requests by IP address**:
```bash
cat logs/app.log | jq -s 'group_by(.ip) | map({ip: .[0].ip, count: length})'
```

---

## 🚨 Error Tracking

### Error Log Levels

1. **WARN (400-499)**: Client errors
   - Invalid request format
   - Unauthorized access
   - Validation failures
   - Resource not found

2. **ERROR (500-599)**: Server errors
   - Database connection failures
   - Internal server errors
   - Unhandled exceptions

### Example Error Investigation

**Step 1**: Find error in logs
```bash
cat logs/app.log | jq 'select(.level == "error") | {request_id, error, path}'
```

**Step 2**: Get full request/response for that request ID
```bash
cat logs/app.log | jq 'select(.request_id == "YOUR_REQUEST_ID")'
```

**Step 3**: Check database logs if needed
```bash
cat logs/app.log | jq 'select(.msg | contains("database"))'
```

---

## 🔐 Security Considerations

### Sensitive Data Handling

**Current Implementation**:
- ✅ Request/Response bodies are logged
- ⚠️ **Warning**: Passwords, tokens, and sensitive data will be visible in logs

**Recommendations for Production**:

1. **Filter Sensitive Fields**:
```go
// Add to logger middleware
func sanitizeRequestBody(body string) string {
    // Remove password, token, etc.
    // Implementation example:
    var data map[string]interface{}
    json.Unmarshal([]byte(body), &data)
    
    sensitiveFields := []string{"password", "token", "secret", "api_key"}
    for _, field := range sensitiveFields {
        if _, exists := data[field]; exists {
            data[field] = "***REDACTED***"
        }
    }
    
    sanitized, _ := json.Marshal(data)
    return string(sanitized)
}
```

2. **Disable Body Logging for Sensitive Endpoints**:
```go
// Skip logging for auth endpoints
if strings.Contains(c.Request.URL.Path, "/auth/") {
    requestBody = []byte("***REDACTED***")
}
```

3. **Use Environment-based Logging**:
```go
// Only log request/response bodies in development
if os.Getenv("ENVIRONMENT") != "production" {
    logger.Infow("...", "request_body", string(requestBody))
}
```

---

## 📦 File Structure

```
project-root/
├── cmd/
│   └── api/
│       └── main.go                    # Logger initialization
├── internal/
│   ├── middleware/
│   │   └── logger_middleware.go       # Request/Response logging
│   └── router/
│       └── router.go                  # Logger middleware registration
├── pkg/
│   └── logger/
│       └── logger.go                  # Core logger configuration
└── logs/
    └── app.log                        # Log output file
```

---

## 🎯 Best Practices

### 1. Request ID Usage
Always include request ID when logging within handlers:
```go
requestID, _ := c.Get("request_id")
sugar.Infow("Processing request", "request_id", requestID, "user_id", userID)
```

### 2. Structured Logging
Use key-value pairs instead of formatted strings:
```go
// ✅ Good
sugar.Infow("User created", "user_id", userID, "email", email)

// ❌ Bad
sugar.Infof("User created: %s (%s)", userID, email)
```

### 3. Log Levels
- `Debug`: Development/debugging information
- `Info`: Normal operations, successful requests
- `Warn`: Deprecations, client errors (4xx)
- `Error`: Server errors (5xx), failed operations

### 4. Performance
- Response body is limited to 1000 characters
- Consider disabling body logging in high-traffic scenarios
- Use sampling for very high volume endpoints

---

## 🧪 Testing Logs

### Example Request Test

```bash
# Send request
curl -X POST http://localhost:8080/api/v1/friends/request \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $TOKEN" \
  -d '{"friend_email":"test@test.com"}'

# Check logs
tail -1 logs/app.log | jq .
```

**Expected Output**:
```json
{
  "level": "info",
  "ts": "2026-01-03T12:45:30.456+0530",
  "caller": "middleware/logger_middleware.go:75",
  "msg": "Request completed",
  "request_id": "a1b2c3d4-e5f6-7890-1234-567890abcdef",
  "method": "POST",
  "path": "/api/v1/friends/request",
  "status": 201,
  "duration_ms": 333,
  "response_body": "{\"success\":true,\"message\":\"Friend request sent successfully\",\"data\":null}"
}
```

---

## 🔄 Log Rotation (Optional)

To enable automatic log rotation, you can use **lumberjack**:

```go
// In pkg/logger/logger.go
import "gopkg.in/natefinch/lumberjack.v2"

func InitLogger() (*zap.SugaredLogger, error) {
    // ... existing code ...
    
    // Add log rotation
    logWriter := &lumberjack.Logger{
        Filename:   "logs/app.log",
        MaxSize:    100, // megabytes
        MaxBackups: 3,
        MaxAge:     28,  // days
        Compress:   true,
    }
    
    // Use logWriter instead of file
    // ... rest of implementation
}
```

---

## 📊 Metrics & Analytics

### Key Metrics to Track

1. **Request Volume**
   ```bash
   cat logs/app.log | jq -s 'map(select(.msg == "Request completed")) | length'
   ```

2. **Average Response Time**
   ```bash
   cat logs/app.log | jq -s 'map(select(.duration_ms)) | map(.duration_ms) | add / length'
   ```

3. **Error Rate**
   ```bash
   errors=$(cat logs/app.log | jq 'select(.status >= 400)' | wc -l)
   total=$(cat logs/app.log | jq 'select(.status)' | wc -l)
   echo "Error rate: $(echo "scale=2; $errors / $total * 100" | bc)%"
   ```

4. **Top Endpoints by Volume**
   ```bash
   cat logs/app.log | jq -s 'group_by(.path) | map({path: .[0].path, count: length}) | sort_by(.count) | reverse | .[0:10]'
   ```

---

## ✅ Summary

**Logging System Features**:
- ✅ Dual output (console + file)
- ✅ Unique Request ID per request
- ✅ Request/Response body capture
- ✅ Duration tracking
- ✅ Error tracking with context
- ✅ Structured JSON logging
- ✅ Environment-based log levels
- ✅ Automatic directory creation
- ✅ IP and User-Agent tracking

**Log File Location**: `logs/app.log`

**Query Logs**: Use `jq` for JSON log analysis

**Monitor Logs**: `tail -f logs/app.log | jq .`

---

## 🎉 Next Steps

1. **Production Deployment**: Implement sensitive data filtering
2. **Log Aggregation**: Consider ELK stack (Elasticsearch, Logstash, Kibana)
3. **Alerting**: Set up alerts for error rate thresholds
4. **Metrics**: Export metrics to Prometheus/Grafana
5. **Log Retention**: Implement automatic cleanup of old logs

---

**For questions or improvements, please refer to the codebase or update this documentation.**
