package helper

import (
	"github.com/gin-gonic/gin"
	"strconv"
)

func GetQuery(c *gin.Context, key, defaultVal string) string {
	if s, ok := c.GetQuery(key); ok {
		return s
	}
	return defaultVal
}

func GetQueryToUint(c *gin.Context, key string, defaultVal uint) uint {
	if s, ok := c.GetQuery(key); ok {
		value, _ := strconv.ParseUint(s, 10, 64)
		if value == 0 {
			return defaultVal
		}
		return uint(value)
	}
	return defaultVal
}

func GetQueryToInt(c *gin.Context, key string, defaultVal int) int {
	if s, ok := c.GetQuery(key); ok {
		value, _ := strconv.Atoi(s)
		if value == 0 {
			return defaultVal
		}
		return value
	}
	return defaultVal
}

func GetQueryToFloat64(c *gin.Context, key string, defaultVal float64) float64 {
	if s, ok := c.GetQuery(key); ok {
		value, _ := strconv.ParseFloat(s, 64)
		if value == 0 {
			return defaultVal
		}
		return value
	}
	return defaultVal
}

func GetQueryToBool(c *gin.Context, key string, defaultVal bool) bool {
	if s, ok := c.GetQuery(key); ok {
		if s == "true" {
			return true
		}
		return false
	}
	return defaultVal
}
