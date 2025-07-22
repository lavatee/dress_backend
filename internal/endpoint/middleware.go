package endpoint

import (
	"errors"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

func (e *Endpoint) Middleware(c *gin.Context) {
	header := c.GetHeader("Authorization")
	sliceOfHeaders := strings.Split(header, " ")
	if len(sliceOfHeaders) != 2 || sliceOfHeaders[0] != "Bearer" {
		c.AbortWithStatusJSON(http.StatusUnauthorized, map[string]string{
			"message": "middleware: invalid header",
		})
		return
	}
	claims, err := e.services.Auth.ParseToken(sliceOfHeaders[1])
	if err != nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, map[string]string{
			"message": "middleware: " + err.Error(),
		})
		return
	}
	c.Set("user_id", claims["id"])
	return
}

func (e *Endpoint) GetUserId(c *gin.Context) (int, error) {
	userId, exists := c.Get("user_id")
	if !exists {
		return 0, errors.New("middleware: user_id not found")
	}

	if id, ok := userId.(float64); ok {
		return int(id), nil
	}

	return 0, errors.New("middleware: invalid user_id type")
}
