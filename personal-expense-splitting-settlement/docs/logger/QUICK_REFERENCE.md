# Quick Reference - Logging System

## 🚀 Quick Start

### View Logs in Real-time
```bash
tail -f logs/app.log | jq .
```

### Common Log Queries

**Find Errors**:
```bash
cat logs/app.log | jq 'select(.level == "error")'
```

**Track Specific Request**:
```bash
cat logs/app.log | jq 'select(.request_id == "YOUR_REQUEST_ID")'
```

**Find Slow Requests (> 1s)**:
```bash
cat logs/app.log | jq 'select(.duration_ms > 1000)'
```

**Filter by Endpoint**:
```bash
cat logs/app.log | jq 'select(.path == "/api/v1/friends/request")'
```

**Count Requests by Status**:
```bash
cat logs/app.log | jq -s 'group_by(.status) | map({status: .[0].status, count: length})'
```

---

## 📋 Log Entry Example

```json
{
  "level": "info",
  "ts": "2026-01-03T12:45:30.456+0530",
  "msg": "Request completed",
  "request_id": "a1b2c3d4-e5f6-7890-1234-567890abcdef",
  "method": "POST",
  "path": "/api/v1/friends/request",
  "status": 201,
  "duration_ms": 333,
  "duration": "333.456789ms",
  "ip": "::1"
}
```

---

## 🔑 Key Fields

| Field | Description |
|-------|-------------|
| `request_id` | Unique request identifier |
| `method` | HTTP method (GET, POST, etc.) |
| `path` | Request endpoint |
| `status` | HTTP status code |
| `duration_ms` | Request duration in milliseconds |
| `ip` | Client IP address |
| `request_body` | Request payload |
| `response_body` | Response payload (max 1000 chars) |

---

## 📊 Useful Metrics

**Average Response Time**:
```bash
cat logs/app.log | jq -s 'map(select(.duration_ms)) | map(.duration_ms) | add / length'
```

**Error Rate**:
```bash
total=$(cat logs/app.log | jq 'select(.status)' | wc -l)
errors=$(cat logs/app.log | jq 'select(.status >= 400)' | wc -l)
echo "Error Rate: $(echo "scale=2; $errors / $total * 100" | bc)%"
```

**Top 10 Slowest Requests**:
```bash
cat logs/app.log | jq -s 'map(select(.duration_ms)) | sort_by(.duration_ms) | reverse | .[0:10]'
```

---

## 🎯 Development Tips

1. **Monitor logs during testing**: `tail -f logs/app.log | jq .`
2. **Check last error**: `cat logs/app.log | jq 'select(.level == "error")' | tail -1`
3. **Filter sensitive endpoints**: Look for auth/login requests
4. **Track database queries**: Search for "database" in msg field

---

For complete documentation, see [LOGGING_SYSTEM.md](LOGGING_SYSTEM.md)
