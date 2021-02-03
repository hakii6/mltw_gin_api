package main

import (
	"net/http"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/gomodule/redigo/redis"
)

type Gacha struct {
	ID string `json:"ID"`
	Name_jp string `json:"NameJP"`
	Name_tw string `json:"NameTW"`
	StartDate JSONTime `json:"StartDate"`
	EndDate JSONTime `json:"EndDate"`
	Ori_url string `json:"Ori_url"`
	Image string `json:"Image"`
	Cards []Card `gorm:"polymorphic:GetCard;polymorphicValue:Gacha"`
}

// func IndexGachas(c *gin.Context) {
// 	var gachas []Gacha
// 	temp := db.Preload("Cards")
// 	if filter["year"] != nil {
// 		temp.Where("(start_date BETWEEN ? AND ?) OR (end_date BETWEEN ? AND ?)", 
// 			filter["year"][0] + "-01-01", filter["year"][0] + "-12-31", filter["year"][0] + "-01-01", filter["year"][0] + "-12-31")
// 	}
// 	res := temp.Find(&gachas)
// 	CheckError(res.Error)
// }

// func (gacha *Gacha) Show(c *gin.Context) {
// 	res := db.Where("id = ?", id).First(&gacha)
// 	CheckError(res.Error)
// }



func IndexGachas(c *gin.Context) {
	var gachas []Gacha

	if c.Keys["Cached"] != false {
		v, _ := redis.String(redisConn.Do("GET", c.FullPath())) // reply from GET
		json.Unmarshal([]byte(v), &gachas)
	} else {
		res := db.Preload("Cards").Find(&gachas)
		if res.Error != nil {
			SetErrCode(c, res.Error)
			return
		}

		json_raw, _ := json.Marshal(gachas)

		_, err := redisConn.Do("SET", c.FullPath(), string(json_raw), "EX", "86400")
		CheckError(err)
	}

	c.JSON(http.StatusOK, gachas)

}

func ShowGacha(c *gin.Context) {
	var gacha Gacha
	res := db.Preload("Cards").Where("id = ?", c.Param("id")).First(&gacha)
	if res.Error != nil {
		SetErrCode(c, res.Error)
		return
	}

	c.JSON(http.StatusOK, gacha)
}