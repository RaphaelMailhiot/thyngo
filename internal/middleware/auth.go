package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"thyngo/internal/auth"
)

const (
	UserIDKey = "user_id"
	RoleKey   = "user_role"
)

// RequireAuth is a Gin middleware that validates the JWT token.
// It extracts the token from the Authorization header (Bearer token)
// and sets the user_id and role inside the Gin context.
func RequireAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Authorization header missing"})
			return
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid authorization format. Expected 'Bearer <token>'"})
			return
		}

		tokenString := parts[1]
		claims, err := auth.ValidateToken(tokenString)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid or expired token"})
			return
		}

		// Set the user information in the context so subsequent handlers can access it
		c.Set(UserIDKey, claims.UserID)
		c.Set(RoleKey, claims.Role)

		c.Next()
	}
}

// RequireRole is a Gin middleware that ensures the user has a specific role.
// It must be used AFTER RequireAuth.
func RequireRole(requiredRoles ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		roleVal, exists := c.Get(RoleKey)
		if !exists {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "User role not found in context"})
			return
		}

		userRole, ok := roleVal.(string)
		if !ok {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Invalid role format in context"})
			return
		}

		hasRole := false
		for _, requiredRole := range requiredRoles {
			if userRole == requiredRole {
				hasRole = true
				break
			}
		}

		if !hasRole {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "You don't have permission to perform this action"})
			return
		}

		c.Next()
	}
}
