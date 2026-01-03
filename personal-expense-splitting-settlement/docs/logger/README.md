# Logger Documentation

This folder contains comprehensive documentation for the application's logging system.

## 📚 Documentation Files

### 1. [LOGGING_SYSTEM.md](LOGGING_SYSTEM.md)
**Complete logging system documentation**
- Overview and features
- Configuration details
- Log entry structure
- Implementation details
- Querying and monitoring examples
- Security considerations
- Best practices

**Read this for**: In-depth understanding of the logging system

---

### 2. [QUICK_REFERENCE.md](QUICK_REFERENCE.md)
**Quick reference guide**
- Common log queries
- Key fields reference
- Useful metrics commands
- Development tips

**Read this for**: Day-to-day logging operations

---

## 🎯 Quick Start

### View Logs in Real-time
```bash
tail -f logs/app.log | jq .
```

### Test a Request
```bash
# Send a request
curl -X GET http://localhost:8080/api/v1/friends/pending \
  -H "Authorization: Bearer $TOKEN"

# Check the logs
tail -5 logs/app.log | jq .
```

### Find Errors
```bash
cat logs/app.log | jq 'select(.level == "error")'
```

---

## 📊 What Gets Logged

✅ **Every HTTP Request**:
- Request ID (unique UUID)
- Method, Path, Query params
- Client IP and User-Agent
- Request body (JSON)
- Headers (Content-Type)

✅ **Every HTTP Response**:
- Status code
- Duration (milliseconds)
- Response size
- Response body (max 1000 chars)

✅ **All Errors**:
- Error message
- Stack trace
- Request context

✅ **Database Operations**:
- SQL queries
- Execution time
- Affected rows

---

## 🔧 Key Features

1. **Dual Output**: Console (human-readable) + File (JSON)
2. **Request Tracking**: Unique ID per request
3. **Performance Metrics**: Duration tracking for every request
4. **Error Context**: Full request context for debugging
5. **Searchable**: JSON format for easy querying with jq

---

## 📁 Log Files

```
project-root/
└── logs/
    └── app.log     # All application logs in JSON format
```

**Note**: Logs directory is auto-created and excluded from Git

---

## 🔍 Example Log Entry

**Request Log**:
```json
{
  "level": "info",
  "ts": "2026-01-03T13:14:30.123+0530",
  "msg": "Incoming request",
  "request_id": "a1b2c3d4-e5f6-7890-1234-567890abcdef",
  "method": "POST",
  "path": "/api/v1/friends/request",
  "ip": "::1",
  "user_agent": "curl/8.14.1",
  "request_body": "{\"friend_email\":\"bob@test.com\"}"
}
```

**Response Log**:
```json
{
  "level": "info",
  "ts": "2026-01-03T13:14:30.456+0530",
  "msg": "Request completed",
  "request_id": "a1b2c3d4-e5f6-7890-1234-567890abcdef",
  "method": "POST",
  "path": "/api/v1/friends/request",
  "status": 201,
  "duration_ms": 333,
  "response_body": "{\"success\":true,\"message\":\"Friend request sent successfully\"}"
}
```

---

## 🛠️ Common Tasks

### Monitor Logs During Development
```bash
tail -f logs/app.log | jq .
```

### Find Slow Requests (> 1 second)
```bash
cat logs/app.log | jq 'select(.duration_ms > 1000)'
```

### Track Specific Request by ID
```bash
cat logs/app.log | jq 'select(.request_id == "YOUR_REQUEST_ID")'
```

### Count Requests by Status
```bash
cat logs/app.log | jq -s 'group_by(.status) | map({status: .[0].status, count: length})'
```

### Calculate Average Response Time
```bash
cat logs/app.log | jq -s 'map(select(.duration_ms)) | map(.duration_ms) | add / length'
```

---

## 📦 Implementation

The logging system consists of:

1. **Logger Configuration** (`pkg/logger/logger.go`)
   - Initializes Zap logger
   - Configures dual output (file + console)
   - Sets log levels based on environment

2. **Logger Middleware** (`internal/middleware/logger_middleware.go`)
   - Generates request IDs
   - Captures request/response details
   - Logs incoming requests and responses
   - Tracks duration and errors

3. **Router Integration** (`internal/router/router.go`)
   - Registers logger middleware globally
   - Applied to all routes

---

## 🔐 Security Note

⚠️ **Current Implementation**: Request and response bodies are logged in full

**For Production**: Implement filtering for sensitive data (passwords, tokens, API keys)

See [LOGGING_SYSTEM.md](LOGGING_SYSTEM.md#security-considerations) for implementation details.

---

## 📈 Monitoring & Analytics

The structured JSON format enables:
- Performance analysis
- Error rate tracking
- Endpoint usage statistics
- Response time trends
- Client behavior analysis

Tools you can use:
- **jq**: Command-line JSON processor
- **ELK Stack**: Elasticsearch, Logstash, Kibana
- **Grafana**: Visualization and dashboards
- **Prometheus**: Metrics and alerting

---

## ✅ Next Steps

1. **Development**: Monitor logs while testing with `tail -f logs/app.log | jq .`
2. **Production**: Implement sensitive data filtering
3. **Scaling**: Set up log aggregation (ELK, Splunk, etc.)
4. **Monitoring**: Create dashboards for key metrics
5. **Alerting**: Set up alerts for error rate thresholds

---

## 🆘 Troubleshooting

**Logs not appearing?**
- Check if `logs/` directory exists
- Verify file permissions
- Check `ENVIRONMENT` variable

**Too much log data?**
- Response bodies are limited to 1000 chars
- Consider filtering sensitive endpoints
- Implement log sampling for high-traffic routes

**Can't read logs?**
- Install jq: `sudo apt install jq` or `brew install jq`
- Use `cat logs/app.log | jq .` for pretty-printing

---

For detailed information, see [LOGGING_SYSTEM.md](LOGGING_SYSTEM.md)

**Questions or improvements?** Update this documentation or refer to the codebase.
