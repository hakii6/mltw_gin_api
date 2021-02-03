package main

import (
	"net/http"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/gomodule/redigo/redis"
)

type Song struct {
	ID string `json:"ID"`
	NameJP string `json:"NameJP"`
	NameTW string `json:"NameTW"`
	BPM string `json:"BPM"`
	Length string `json:"Length"`
	Date JSONTime `json:"Date"`
	Image string `json:"Image"`
	Type string `json:"Type"`

	EzLv int `json:"EzLv"`
	NmLv int `json:"NmLv"`
	HrLv int `json:"HrLv"`
	Hr2Lv int `json:"Hr2Lv"`
	ExLv int `json:"ExLv"`

	EzNotes int `json:"EzNotes"`
	NmNotes int `json:"NmNotes"`
	HrNotes int `json:"HrNotes"`
	Hr2Notes int `json:"Hr2Notes"`
	ExNotes int `json:"ExNotes"`
}

// func IndexSongs(c *gin.Context) []Song {
// 	var songs []Song
// 	temp := db.Select("id", "name_jp", "name_tw", "BPM", "length", "image", "date", "type")
// 	if filter["type"] != nil {
// 		temp.Where("type = ?", filter["type"][0])
// 	}
// 	if filter["year"] != nil {
// 		temp.Where("date BETWEEN ? AND ?", filter["year"][0] + "-01-01", filter["year"][0] + "-12-31")
// 	}
// 	// for key, _ := range filter {
// 	// 	temp.Where("type = ?", key)
// 	// }
// 	res := temp.Find(&songs)
// 	CheckError(res.Error)

// 	return songs
// }

func IndexSongs(c *gin.Context) {
	var songs []Song

	if c.Keys["Cached"] != false {
		v, _ := redis.String(redisConn.Do("GET", c.FullPath())) // reply from GET
		json.Unmarshal([]byte(v), &songs)
	} else {
		res := db.Find(&songs)
		if res.Error != nil {
			SetErrCode(c, res.Error)
			return
		}

		json_raw, _ := json.Marshal(songs)

		_, err := redisConn.Do("SET", c.FullPath(), string(json_raw), "EX", "86400")
		CheckError(err)
	}

	c.JSON(http.StatusOK, songs)

}

func ShowSong(c *gin.Context) {
	var song Song
	res := db.Where("id = ?", c.Param("id")).First(&song)
	if res.Error != nil {
		SetErrCode(c, res.Error)
		return
	}

	c.JSON(http.StatusOK, song)
}

// func (song *Song) Show(c *gin.Context) *Song{
// 	res := db.Where("id = ?", id).First(&song)
// 	CheckError(res.Error)

// 	return song
// }