---
name: go-http-transport
description: Use when building or auditing the HTTP layer: chi router, middleware stack (Content-Type, Timeout, CORS, RateLimit), per-IP rate limiting via in-memory cache, context_key constants, unified ErrorResponse.
category: methodology
go-version-min: "1.21"
sources:
  - go-old-pattern/golang-service-methodology.md (§8)
related:
  - skills/methodology/00-canonical-full.md
---

## 8. HTTP transport

### Router

`go-chi/chi/v5` — lightweight, idiomatic, fully compatible with `net/http`.

### Standard middleware stack

```go
r.Use(
    middleware.JSON(),                       // Content-Type
    chiMiddleware.Timeout(20*time.Second),   // global timeout
    middleware.Cors(allowedCORSOrigins),     // CORS
    middleware.RateLimit(guard),             // per-IP rate limiting
    // tracing.Middleware(),                 // optional
)
```

### Custom middleware example

```go
func JSON() func(http.Handler) http.Handler {
    return func(next http.Handler) http.Handler {
        return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
            w.Header().Set("Content-Type", "application/json; charset=utf-8")
            next.ServeHTTP(w, r)
        })
    }
}
```

### Rate limiting (in-memory cache)

Approach: per-IP counters in 1s / 1m / 1h windows. When a window exceeds its limit, write a block key with a TTL.

```go
func RateLimit(c *cache.Cache) func(http.Handler) http.Handler {
    return func(next http.Handler) http.Handler {
        return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
            ip := getClientIP(r)
            if _, blocked := c.Get(keyBlocked + ip); blocked {
                http.Error(w, `{"error":"rate limit exceeded"}`, http.StatusTooManyRequests)
                return
            }
            if incr(c, keySec+ip, time.Second) > maxPerSecond {
                c.Set(keyBlocked+ip, true, blockTTL)
                http.Error(w, `{"error":"rate limit exceeded"}`, http.StatusTooManyRequests)
                return
            }
            next.ServeHTTP(w, r)
        })
    }
}
```

### Passing data via context

Key constants live in `pkg/common/middleware/context_key/`:

```go
type contextKey string
const ClientIP contextKey = "client_ip"

// in middleware
ctx := context.WithValue(r.Context(), context_key.ClientIP, ip)
next.ServeHTTP(w, r.WithContext(ctx))

// in handler
ip, _ := r.Context().Value(context_key.ClientIP).(string)
```

### Unified error response

```go
// pkg/common/error.go
type ErrorResponse struct {
    StatusCode int    `json:"-"`
    StatusText string `json:"status"`
    ErrorText  string `json:"error"`
}

func NewErrorResponse(msg string, code int) render.Renderer {
    return &ErrorResponse{
        StatusCode: code,
        StatusText: http.StatusText(code),
        ErrorText:  msg,
    }
}

func (e *ErrorResponse) Render(_ http.ResponseWriter, r *http.Request) error {
    render.Status(r, e.StatusCode)
    return nil
}
```
