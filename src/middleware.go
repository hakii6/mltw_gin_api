package main

import (

	"github.com/gin-gonic/gin"
	"github.com/gomodule/redigo/redis"

)

func SetHeader() gin.HandlerFunc {
    return func(c *gin.Context) {

		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Content-Type", "application/json; charset=utf-8")
	    c.Next()
    }
}

func CheckCache() gin.HandlerFunc {
	return func(c *gin.Context) {
				
		exists, err := redis.Bool(redisConn.Do("Exists", c.FullPath()))
		CheckError(err)

		// Set example variable
		c.Set("Cached", exists)

		// before request
		c.Next()
	}
}
