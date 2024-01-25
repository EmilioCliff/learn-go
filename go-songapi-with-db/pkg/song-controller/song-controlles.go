package controller

import (
	"fmt"
	"encoding/json"
	"github.com/gorilla/mux"
	"net/http"
	"strconv"
	"github.com/EmilioCliff/learn-go/go-songapi-with-db/pkg/utils"
	"github.com/EmilioCliff/learn-go/go-songapi-with-db/pkg/models"
)

var NewSong models.Song

func GetSongs(w http.ResponseWriter, r *http.Request){
	NewSong := models.GetAllSongs()
	res, _ := json.Marshal(NewSong)
	w.Header().Set("Content-Type", "pkglication/json")
	w.WriteHeader(http.StatusOK)
	w.Write(res)
}

func GetSongById(w http.ResponseWriter, r *http.Request){
	params := mux.Vars(r)
	songId := params["id"]
	id, err := strconv.ParseInt(songId, 0,0)
	if err != nil{
		fmt.Println("Error while parsing")
	}
	song, _ := models.GetSongById(id)
	res, _ := json.Marshal(song)
	w.Header().Set("Content-Type", "pkglication/json")
	w.WriteHeader(http.StatusOK)
	w.Write(res)
}

func CreateSong(w http.ResponseWriter, r *http.Request){
	createSong := &models.Song{}
	utils.ParseBody(r, createSong)
	s := createSong.CreateSong()
	res, _ := json.Marshal(s)
	// w.Header().Set("Content-Type", "pkglication")
	w.WriteHeader(http.StatusOK)
	w.Write(res)
}

func DeleteSong(w http.ResponseWriter, r *http.Request){
	params := mux.Vars(r)
	songId, err := strconv.ParseInt(params["id"], 0,0)
	if err != nil{
		fmt.Println("Error while parsing")
	}
	song := models.DeleteSong(songId)
	res, _ := json.Marshal(song)
	w.Header().Set("Content-Type", "pkglication/json")
	w.WriteHeader(http.StatusOK)
	w.Write(res)
}

func UpdateSong(w http.ResponseWriter, r *http.Request){
	updateSong := &models.Song{}
	utils.ParseBody(r, updateSong)
	params := mux.Vars(r)
	songId, err := strconv.ParseInt(params["id"], 0,0)
	if err != nil{
		fmt.Println("Error while parsing")
	}
	song, db := models.GetSongById(songId)
	if updateSong.Name != ""{
		song.Name = updateSong.Name
	}
	if updateSong.Duration != ""{
		song.Duration = updateSong.Duration
	}
	if updateSong.Album != ""{
		song.Album = updateSong.Album
	}
	db.Save(&song)
	res, _ := json.Marshal(song)
	w.Header().Set("Content-Type", "pkglication/json")
	w.WriteHeader(http.StatusOK)
	w.Write(res)
}