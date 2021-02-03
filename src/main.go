package main

import (
	"net/http"

	"time"
	// "encoding/json"
	"fmt"
	"gorm.io/gorm"
	"github.com/gin-gonic/gin"
	"github.com/gomodule/redigo/redis"
	"gorm.io/driver/mysql"

)

type JSONTime time.Time

func (t JSONTime)MarshalJSON() ([]byte, error) {
    //do your serializing here
    stamp := fmt.Sprintf("\"%s\"", time.Time(t).Format("Mon Jan _2"))
    return []byte(stamp), nil
}

func CheckError(err error) {
    if err != nil {
        panic(err)
    }
}

func SetErrCode(c *gin.Context, err error) {
	switch err {
		case gorm.ErrRecordNotFound:
			c.JSON(http.StatusNotFound, gin.H{
				"error": "not found",
			})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "unknown",
			})
	}
}

var db *gorm.DB
var redisConn redis.Conn

func main() {

	db, _ = gorm.Open(mysql.Open(connect), &gorm.Config{})
	redisConn, _ = redis.Dial("tcp", "127.0.0.1:6379")
	r := gin.Default()
	r.Use(SetHeader())
	r.Use(CheckCache())

	go CronJob()

	// Simple group: v2
	v2 := r.Group("/v2")
	{
		v2.GET("/events", IndexEvents)
		v2.GET("/events/:id", ShowEvent)

		v2.GET("/gachas", IndexGachas)
		v2.GET("/gachas/:id", ShowGacha)

		v2.GET("/idols", IndexIdols)
		v2.GET("/idols/:id", ShowIdol)

		v2.GET("/songs", IndexSongs)
		v2.GET("/songs/:id", ShowSong)

		v2.GET("/cards", IndexCards)
		v2.GET("/cards/:id", ShowCard)

	}
	r.Run(":8002")
}

