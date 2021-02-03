package main

import (
	"net/http"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/gomodule/redigo/redis"
)

type Event struct {
	ID string `json:"ID"`
	Name_jp string `json:"NameJP"`
	Name_tw string `json:"NameTW"`
	StartDate JSONTime `json:"StartDate"`
	BoostDate JSONTime `json:"BoostDate"`
	EndDate JSONTime `json:"EndDate"`
	Ori_url string `json:"Ori_url"`
	Type string `json:"Type"`
	Image string `json:"Image"`
	Cards []Card `gorm:"polymorphic:GetCard;polymorphicValue:Event"`
	Api_ID string `json:"Api_ID`

}

// func IndexEvents(c *gin.Context) {
// 	var events []Event
// 	temp := db.Preload("Cards")
// 	if filter["year"] != nil {
// 		temp.Where("(start_date BETWEEN ? AND ?) OR (end_date BETWEEN ? AND ?)", 
// 			filter["year"][0] + "-01-01", filter["year"][0] + "-12-31", filter["year"][0] + "-01-01", filter["year"][0] + "-12-31")
// 	}
// 	res := temp.Find(&events)
// 	CheckError(res.Error)
// }

// func (event *Event) Show(c *gin.Context) {
// 	res := db.Preload("Cards").Where("id = ?", id).First(&event)
// 	CheckError(res.Error)

// }

func IndexEvents(c *gin.Context) {
	var events []Event

	if c.Keys["Cached"] != false {
		v, _ := redis.String(redisConn.Do("GET", c.FullPath())) // reply from GET
		json.Unmarshal([]byte(v), &events)
	} else {
		res := db.Preload("Cards").Find(&events)
		if res.Error != nil {
			SetErrCode(c, res.Error)
			return
		}

		json_raw, _ := json.Marshal(events)

		_, err := redisConn.Do("SET", c.FullPath(), string(json_raw), "EX", "86400")
		CheckError(err)
	}

	c.JSON(http.StatusOK, events)

}

func ShowEvent(c *gin.Context) {
	var event Event
	res := db.Preload("Cards").Where("id = ?", c.Param("id")).First(&event)
	if res.Error != nil {
		SetErrCode(c, res.Error)
		return
	}

	c.JSON(http.StatusOK, event)
}