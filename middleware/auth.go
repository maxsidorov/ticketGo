package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/gin-contrib/sessions"
	"net/http"
)

// AuthRequired middleware проверяет, авторизован ли пользователь
func AuthRequired() gin.HandlerFunc {
	return func(c *gin.Context) {
		session := sessions.Default(c)
		userID := session.Get("user_id")
		if userID == nil {
			c.Redirect(http.StatusFound, "/login")
			c.Abort()
			return
		}
		c.Next()
	}
}

func SetAuthStatus() gin.HandlerFunc {
	return func(c *gin.Context) {
		session := sessions.Default(c)
		userID := session.Get("user_id")
		username := session.Get("username")
		
		c.Set("is_authenticated", userID != nil)
		if username != nil {
			c.Set("username", username)
		}
		if userID != nil {
			if userIDInt, ok := userID.(int); ok {
				c.Set("user_id", uint(userIDInt))
			}
		}
		
		c.Next()
	}
} 