package main

import (
	"net/http"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/gomodule/redigo/redis"
)

type Card struct {
	ID string `json:"ID"`
	Idol Idol
	IdolID string `json:"IdolID"`

	Name_jp string `json:"NameJP"`
	Name_tw string `json:"NameTW"`
	Rarity string `json:"Rarity"`
	Total int `json:"Total"`
	Vocal int `json:"Vocal"`
	Dance int `json:"Dance"`
	Visual int `json:"Visual"`
	Limited string `json:"Limited"`
	Date string `json:"Date"`

	ImageA string `json:"ImageA"`
	ImageB string `json:"ImageB"`

	GetCardID string `json:"GetCardID"`
	GetCardType string `json:"GetCardType"`
}

// func IndexCards(c *gin.Context) []Card {
// 	var cards []Card
// 	temp := db.Select("id", "name_jp", "name_tw", "rarity", "total", "vocal", "dance", "visual", "limited", "date")
// 	if filter["type"] != nil {
// 		temp.Where("type = ?", filter["type"][0])
// 	}
// 	if filter["rarity"] != nil {
// 		temp.Where("rarity = ?", filter["rarity"][0])
// 	}
// 	if filter["year"] != nil {
// 		temp.Where("date BETWEEN ? AND ?", filter["year"][0] + "-01-01", filter["year"][0] + "-12-31")
// 	}
// 	res := temp.Find(&cards)
// 	CheckError(res.Error)
// 	return cards
// }

// func IndexCards(c *gin.Context) {
// 	var cards []Card
// 	temp := db.Select("id", "name_jp", "name_tw", "rarity", "total", "vocal", "dance", "visual", "limited", "date")
// 	if filter["type"] != nil {
// 		temp.Where("type = ?", filter["type"][0])
// 	}
// 	if filter["rarity"] != nil {
// 		temp.Where("rarity = ?", filter["rarity"][0])
// 	}
// 	if filter["year"] != nil {
// 		temp.Where("date BETWEEN ? AND ?", filter["year"][0] + "-01-01", filter["year"][0] + "-12-31")
// 	}
// 	res := temp.Find(&cards)
// 	CheckError(res.Error)
// }


// func (card *Card) Show(c *gin.Context) {
// 	res := db.Preload("Idol").Where("id = ?", id).First(&card)
// 	CheckError(res.Error)

// }

func IndexCards(c *gin.Context) {
	var cards []Card

	if c.Keys["Cached"] != false {
		v, _ := redis.String(redisConn.Do("GET", c.FullPath())) // reply from GET
		json.Unmarshal([]byte(v), &cards)
	} else {
		res := db.Find(&cards)
		if res.Error != nil {
			SetErrCode(c, res.Error)
			return
		}

		json_raw, _ := json.Marshal(cards)

		_, err := redisConn.Do("SET", c.FullPath(), string(json_raw), "EX", "86400")
		CheckError(err)
	}

	c.JSON(http.StatusOK, cards)

}

func ShowCard(c *gin.Context) {
	var card Card
	res := db.Preload("Idols").Where("id = ?", c.Param("id")).First(&card)
	if res.Error != nil {
		SetErrCode(c, res.Error)
		return
	}

	c.JSON(http.StatusOK, card)
}