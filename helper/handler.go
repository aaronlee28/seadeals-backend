package helper

import "github.com/gin-gonic/gin"

func GetQuery(c *gin.Context, key, defaultVal string) string {
	if s, ok := c.GetQuery(key); ok {
		return s
	}
	return defaultVal
}
