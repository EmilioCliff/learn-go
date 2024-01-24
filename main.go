package main

import (
	"log"
	"github.com/gorilla/mux"
	"net/http"
	"encoding/json"
	"strconv"
	"math/rand"
	"fmt"
)

type Song struct{
	ID string `json: "id"`
	Name string	`json: "name"`
	Duration string `json: "duration"`
	Album string `json: "album"`
	Artist *Artist `json: "artist"`
}
type Artist struct{
	Firstname string `json: "firstname"`
	Lastname string `json: "lastname"`
	Country string `json: "country"`
}

var songs []Song

func getSongs(w http.ResponseWriter, r *http.Request){
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(songs)
}

func createSong(w http.ResponseWriter, r *http.Request){
	w.Header().Set("Content-Type", "application/json")
	var song Song
	json.NewDecoder(r.Body).Decode(&song)
	fmt.Println(song)
	song.ID = strconv.Itoa(rand.Intn(1000))
	songs = append(songs, song)
	json.NewEncoder(w).Encode(songs)
}

func deleteSong(w http.ResponseWriter, r *http.Request){
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	for i, v := range songs{
		if v.ID == params["id"]{
			songs = append(songs[:i], songs[i+1:]...)
			json.NewEncoder(w).Encode(songs)
			break
		}
	}
}

func getSongById(w http.ResponseWriter, r *http.Request){
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	for _,v := range songs{
		if v.ID == params["id"]{
		json.NewEncoder(w).Encode(v)
	}
	}
}

func updateSong(w http.ResponseWriter, r *http.Request){
	w.Header().Set("Content-Type", "application/json")
	var song Song
	json.NewDecoder(r.Body).Decode(&song)
	params := mux.Vars(r)
	for i,v := range songs{
		if v.ID == params["id"]{
			song.ID = v.ID
			songs = append(songs[:i], songs[i+1:]...)
			songs = append(songs, song)
			json.NewEncoder(w).Encode(songs)
		}
	}
}

func main(){
	songs = append(songs, Song{ID: "1",Name:"song 1", Duration:"2.30s", Album: "album 1", Artist: &Artist{Firstname: "John", Lastname: "Doe", Country: "Kenya"}})
	r := mux.NewRouter()
	r.HandleFunc("/songs", getSongs).Methods("GET")
	r.HandleFunc("/songs", createSong).Methods("POST")
	r.HandleFunc("/songs/{id}", getSongById).Methods("GET")
	r.HandleFunc("/songs/{id}", updateSong).Methods("PUT")
	r.HandleFunc("/songs/{id}", deleteSong).Methods("DELETE")

	log.Fatal(http.ListenAndServe(":8080", r))
}
