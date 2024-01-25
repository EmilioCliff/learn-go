package models

import (
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"github.com/EmilioCliff/learn-go/go-songapi-with-db/pkg/config"
)

var db *gorm.DB


type Song struct{
	gorm.Model
	Name string `gorm:""json:"name"`
	Duration string `json:"song"`
	Album string `json:"album"`
}

func init(){
	config.Connect()
	db := config.GetDB()
	db.AutoMigrate(&Song{})
}

func (s *Song) CreateSong() *Song{
	db.NewRecord(s)
	db.Create(&s)
	return s
}

func GetAllSongs() []Song{
	var songs []Song
	db.Find(&songs)
	return songs
}

func GetSongById(Id int64) (Song, *gorm.DB) {
	var getsong Song
	db := db.Where("ID=?", Id).Find(&getsong)
	return getsong, db
}

func DeleteSong(Id int64) Song{
	var song Song
	db.Where("ID=?", Id).Delete(song)
	return song
}