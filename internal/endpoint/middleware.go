package endpoint

import (
	"errors"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func (e *Endpoint) Middleware(c *gin.Context) {
	header := c.GetHeader("Authorization")
	sliceOfHeaders := strings.Split(header, " ")
	if len(sliceOfHeaders) != 2 || sliceOfHeaders[0] != "Bearer" {
		c.Set("user_id", float64(0))
		return
	}
	claims, err := e.services.Auth.ParseToken(sliceOfHeaders[1])
	if err != nil {
		logrus.Errorf("Middleware error: " + err.Error())
		if err.Error() == "expired" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Need to refresh token"})
		} else {
			c.Set("user_id", float64(0))
		}
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
