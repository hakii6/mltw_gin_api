package main

import (
	"net/http"
	"os"
	"log"
	"time"
	// "encoding/json"
	"fmt"
	"gorm.io/gorm"
	"github.com/gin-gonic/gin"
	"github.com/gomodule/redigo/redis"
	"gorm.io/driver/mysql"
	"github.com/kardianos/service"

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
var logger service.Logger

type program struct{}

func (p *program) Start(s service.Service) error {
	// Start should not block. Do the actual work async.
	go p.run()
	return nil
}

func (p *program) run() {

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

func (p *program) Stop(s service.Service) error {
	// Stop should not block. Return with a few seconds.
	return nil
}

func main() {
	svcConfig := &service.Config{
		Name:        "GinAPI",
		DisplayName: "Gin API",
		Description: "Mltw Gin API",
	}
	prg := &program{}
	s, err := service.New(prg, svcConfig)
	if err != nil {
		log.Fatal(err)
	}

	// if in windows install and start uself
	if len(os.Args) > 1 {
		err = service.Control(s, os.Args[1])
		if err != nil {
			log.Fatal(err)
		}
		if os.Args[1] != "start" {
			return
		}
	}

	logger, err = s.Logger(nil)
	if err != nil {
		log.Fatal(err)
	}

	err = s.Run()
	if err != nil {
		logger.Error(err)
	}


}