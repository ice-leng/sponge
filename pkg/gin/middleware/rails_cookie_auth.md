## Rails cookie authorization middleware

Validate and decrypt a Rails 7.1+ encrypted session cookie using `secret_key_base`, then attach the decoded session payload to Gin context under key `rails_session`.

This middleware is useful when a Go service needs to trust Rails session cookies issued by an existing Rails app (SSO between Rails and Go services).

<br>

### Requirements

- Rails 7.1+ encrypted cookies (AES-256-GCM, PBKDF2-HMAC-SHA256, purpose: `cookie.<cookie_name>`)
- Rails `secret_key_base`
- The Rails session cookie name (from `config/initializers/session_store.rb`)

Note: Rails < 7.1 used a different encryption scheme (AES-CBC + HMAC) and is not supported by this middleware.

<br>

### Configuration example

If you keep app configuration in YAML, you can add a section like the following (example from a reference service):

```yaml
rails:
  secretKeyBase: "change-me"                           # run: rails credentials:show
  cookieName:    "_coreui_pro_rails_starter_session"   # from session_store.rb or browser devtools
  userID:        1137                                    # optional: restrict to a specific logged-in user id
```

<br>

### Usage

Attach the middleware at the group or router level where you want Rails cookie authentication enforced:

```go
import (
    "github.com/gin-gonic/gin"
    "github.com/go-dev-frame/sponge/pkg/gin/middleware"
)

func NewRouter(secretKeyBase, cookieName string) *gin.Engine {
    r := gin.Default()

    g := r.Group("/api/v1")
    g.Use(middleware.RailsCookieAuthMiddleware(secretKeyBase, cookieName))

    // routes protected by Rails cookie auth
    g.GET("/me", func(c *gin.Context) { /* ... */ })

    return r
}
```

The middleware:

- Reads the Rails session cookie by `cookieName`
- Verifies and decrypts it using `secretKeyBase`
- Puts the decoded session map into Gin context at key `"rails_session"`

<br>

### Verifying the logged-in user (optional)

If you also want to ensure the request belongs to a specific user id (as stored by Devise/Warden), you can add a follow-up middleware that extracts the user id from the Rails session and compares it. The helper below shows the pattern:

```go
import (
    "strconv"
    "github.com/gin-gonic/gin"
    "github.com/go-dev-frame/sponge/pkg/rails"
)

func VerifyRailsSessionUserIdIs(userID int64) gin.HandlerFunc {
    return func(c *gin.Context) {
        v, ok := c.Get("rails_session")
        if !ok {
            c.AbortWithStatusJSON(401, gin.H{"error": "rails_session missing"})
            return
        }
        session, ok := v.(map[string]any)
        if !ok {
            c.AbortWithStatusJSON(401, gin.H{"error": "invalid rails_session"})
            return
        }
        uidVal, ok := rails.UserIDFromSession(session) // reads warden.user.user.key -> [[id], ...]
        if !ok {
            c.AbortWithStatusJSON(401, gin.H{"error": "user id not found in session"})
            return
        }
        var uid int64
        switch vv := uidVal.(type) {
        case int64:
            uid = vv
        case int:
            uid = int64(vv)
        case float64:
            uid = int64(vv)
        case string:
            parsed, err := strconv.ParseInt(vv, 10, 64)
            if err != nil {
                c.AbortWithStatusJSON(401, gin.H{"error": "invalid user id in session"})
                return
            }
            uid = parsed
        default:
            c.AbortWithStatusJSON(401, gin.H{"error": "invalid user id type in session"})
            return
        }
        if uid != userID {
            c.AbortWithStatusJSON(403, gin.H{"error": "forbidden"})
            return
        }
        c.Next()
    }
}
```

Then chain it after `RailsCookieAuthMiddleware`:

```go
g.Use(middleware.RailsCookieAuthMiddleware(secretKeyBase, cookieName))
g.Use(VerifyRailsSessionUserIdIs(1137))
```

<br>

### Accessing the Rails session in handlers

```go
v, ok := c.Get("rails_session")
if !ok {
    // not present
}
session := v.(map[string]any)

// Example: extract the logged-in user id saved by Warden/Devise
uid, ok := rails.UserIDFromSession(session)
// uid can be number or string depending on Rails serialization
```

Common keys you may see in the session map include `"_csrf_token"` and `"warden.user.user.key"`.

<br>

### Testing with curl

```bash
curl -i \
  -H "Cookie: <cookie_name>=<encrypted_cookie_value>" \
  http://localhost:8080/api/v1/me
```

Replace `<cookie_name>` with your Rails session cookie name (e.g. `_coreui_pro_rails_starter_session`).

<br>

### Troubleshooting

- 401 Missing cookie: ensure the browser sends the Rails session cookie to this domain/path
- 401 Invalid cookie: check `secret_key_base`, cookie name, and that the cookie was created by the same Rails environment; requires Rails 7.1+
- 401 user id not found: verify your app uses Devise/Warden and the session contains `warden.user.user.key`
- 403 forbidden: your additional user id check failed

Security tips: never log the cookie value; use HTTPS in production so cookies are transmitted securely.


