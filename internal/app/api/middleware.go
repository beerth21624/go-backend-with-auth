package api

import (
	"crypto/rand"
	"encoding/hex"
	"net/http"
	"strings"
	"time"

	"errors"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/rs/zerolog/log"

	"venturex-backend/internal/app/service"
)

const (
	RequestIDHeader = "X-Request-ID"
)

func RequestID() gin.HandlerFunc {
	return gin.HandlerFunc(func(c *gin.Context) {
		requestID := c.GetHeader("X-Request-ID")
		if requestID == "" {
			requestID = generateRequestID()
		}

		c.Set("request_id", requestID)

		c.Header("X-Request-ID", requestID)

		c.Next()
	})
}

func generateRequestID() string {
	bytes := make([]byte, 16)
	rand.Read(bytes)
	return hex.EncodeToString(bytes)
}

func LoggerMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()

		requestID := c.Request.Header.Get(RequestIDHeader)
		if requestID == "" {
			requestID = uuid.New().String()
		}
		c.Set("request_id", requestID)
		c.Writer.Header().Set(RequestIDHeader, requestID)

		l := log.With().
			Str("request_id", requestID).
			Str("method", c.Request.Method).
			Str("path", c.Request.URL.Path).
			Str("ip", c.ClientIP()).
			Logger()

		c.Set("logger", l)

		c.Next()

		latency := time.Since(start)
		l.Info().
			Int("status", c.Writer.Status()).
			Dur("latency", latency).
			Msg("Request completed")
	}
}

func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Origin, Content-Type, Authorization, X-Request-ID")
		c.Header("Access-Control-Expose-Headers", "X-Request-ID")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}

		c.Next()
	}
}

func AuthMiddleware(authService service.AuthService) gin.HandlerFunc {
	return gin.HandlerFunc(func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			AbortWithError(c, NewUnauthorizedError("Authorization header required"))
			return
		}

		const bearerPrefix = "Bearer "
		if !strings.HasPrefix(authHeader, bearerPrefix) {
			AbortWithError(c, NewBadRequestError("Invalid authorization header format"))
			return
		}

		token := strings.TrimPrefix(authHeader, bearerPrefix)
		if token == "" {
			AbortWithError(c, NewUnauthorizedError("Access token required"))
			return
		}

		claims, err := authService.ValidateToken(c.Request.Context(), token)
		if err != nil {
			AbortWithError(c, err)
			return
		}

		c.Set("user_id", claims.UserID)
		c.Set("session_id", claims.SessionID)
		c.Set("username", claims.Username)
		c.Set("email", claims.Email)
		c.Set("role", claims.Role)
		c.Set("token_claims", claims)

		c.Set("user_uuid", claims.UserUUID)
		c.Set("session_uuid", claims.SessionUUID)

		c.Next()
	})
}

func RequireRole(allowedRoles ...string) gin.HandlerFunc {
	return gin.HandlerFunc(func(c *gin.Context) {
		role, exists := c.Get("role")
		if !exists {
			AbortWithError(c, NewUnauthorizedError("Authentication required"))
			return
		}

		userRole, ok := role.(string)
		if !ok {
			AbortWithError(c, NewInternalError(errors.New("Invalid role format")))
			return
		}

		for _, allowedRole := range allowedRoles {
			if userRole == allowedRole {
				c.Next()
				return
			}
		}

		AbortWithError(c, NewForbiddenError("Insufficient permissions"))
	})
}

func AdminMiddleware() gin.HandlerFunc {
	return RequireRole("admin")
}

func RateLimitMiddleware(maxRequests int, window string) gin.HandlerFunc {
	// TODO: Implement rate limiting using Redis
	return gin.HandlerFunc(func(c *gin.Context) {
		c.Next()
	})
}

func ValidateJSONMiddleware() gin.HandlerFunc {
	return gin.HandlerFunc(func(c *gin.Context) {
		if c.Request.Method == "POST" || c.Request.Method == "PUT" || c.Request.Method == "PATCH" {
			contentType := c.GetHeader("Content-Type")
			if !strings.Contains(contentType, "application/json") {
				AbortWithError(c, NewBadRequestError("Content-Type must be application/json"))
				return
			}
		}
		c.Next()
	})
}

func SecurityHeaders() gin.HandlerFunc {
	return gin.HandlerFunc(func(c *gin.Context) {
		c.Header("X-Frame-Options", "DENY")
		c.Header("X-Content-Type-Options", "nosniff")
		c.Header("X-XSS-Protection", "1; mode=block")
		c.Header("Referrer-Policy", "strict-origin-when-cross-origin")
		c.Header("Content-Security-Policy", "default-src 'self'")
		c.Next()
	})
}

func GetUserID(c *gin.Context) (int64, bool) {
	userID, exists := c.Get("user_id")
	if !exists {
		return 0, false
	}
	id, ok := userID.(int64)
	return id, ok
}

func GetUsername(c *gin.Context) (string, bool) {
	username, exists := c.Get("username")
	if !exists {
		return "", false
	}
	name, ok := username.(string)
	return name, ok
}

func GetUserEmail(c *gin.Context) (string, bool) {
	email, exists := c.Get("email")
	if !exists {
		return "", false
	}
	userEmail, ok := email.(string)
	return userEmail, ok
}

func GetSessionID(c *gin.Context) (int64, bool) {
	sessionID, exists := c.Get("session_id")
	if !exists {
		return 0, false
	}
	id, ok := sessionID.(int64)
	return id, ok
}

func GetUserRole(c *gin.Context) (string, bool) {
	role, exists := c.Get("role")
	if !exists {
		return "", false
	}
	userRole, ok := role.(string)
	return userRole, ok
}

func GetTokenClaims(c *gin.Context) (*service.AuthClaims, bool) {
	claims, exists := c.Get("token_claims")
	if !exists {
		return nil, false
	}
	tokenClaims, ok := claims.(*service.AuthClaims)
	return tokenClaims, ok
}

func RequireAuth(c *gin.Context) error {
	if _, exists := c.Get("user_id"); !exists {
		return NewUnauthorizedError("Authentication required")
	}
	return nil
}

func GetClientIP(c *gin.Context) string {
	if xff := c.GetHeader("X-Forwarded-For"); xff != "" {
		if ips := strings.Split(xff, ","); len(ips) > 0 {
			return strings.TrimSpace(ips[0])
		}
	}

	if xri := c.GetHeader("X-Real-IP"); xri != "" {
		return xri
	}

	return c.ClientIP()
}

func GetUserAgent(c *gin.Context) string {
	return c.GetHeader("User-Agent")
}

func GetUserUUID(c *gin.Context) (string, bool) {
	userUUID, exists := c.Get("user_uuid")
	if !exists {
		return "", false
	}
	uuid, ok := userUUID.(string)
	return uuid, ok
}

func GetSessionUUID(c *gin.Context) (string, bool) {
	sessionUUID, exists := c.Get("session_uuid")
	if !exists {
		return "", false
	}
	uuid, ok := sessionUUID.(string)
	return uuid, ok
}
