package main

import (
	"net/http"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/gomodule/redigo/redis"
)

type Idol struct {
	ID string `json:"ID"`
	Name_jp string `json:"NameJP"`
	Name_tw string `json:"NameTW"`
	Type string `json:"Type"`
	Thumbnail string `json:"Thumbnail"`
	Intro string `json:"Intro"`
	Songs []Song `gorm:"many2many:idol_song"`
	Cards []Card `gorm:"foreignKey:IdolID"`
}

// func IndexIdols(c *gin.Context) {
// 	var idols []Idol
// 	res := db.Select("id", "name_jp", "name_tw", "type", "thumbnail", "intro").Find(&idols)
// 	CheckError(res.Error)
// }

// func (idol *Idol) Show(c *gin.Context) {
// 	res := db.Preload("Songs").Preload("Cards").Where("id = ?", id).First(&idol)
// 	CheckError(res.Error)

// }

func IndexIdols(c *gin.Context) {
	var idols []Idol

	if c.Keys["Cached"] != false {
		v, _ := redis.String(redisConn.Do("GET", c.FullPath())) // reply from GET
		json.Unmarshal([]byte(v), &idols)
	} else {
		res := db.Find(&idols)
		if res.Error != nil {
			SetErrCode(c, res.Error)
			return
		}

		json_raw, _ := json.Marshal(idols)

		_, err := redisConn.Do("SET", c.FullPath(), string(json_raw), "EX", "86400")
		CheckError(err)
	}

	c.JSON(http.StatusOK, idols)

}

func ShowIdol(c *gin.Context) {
	var idol Idol
	res := db.Preload("Cards").Preload("Songs").Where("id = ?", c.Param("id")).First(&idol)
	if res.Error != nil {
		SetErrCode(c, res.Error)
		return
	}

	c.JSON(http.StatusOK, idol)
}