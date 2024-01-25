package routes

import(
	"github.com/gorilla/mux"
	"github.com/EmilioCliff/learn-go/go-songapi-with-db/pkg/song-controller"
)

var SetRoutes = func(r *mux.Router){
	r.HandleFunc("/songs", controller.GetSongs).Methods("GET")
	r.HandleFunc("/songs", controller.CreateSong).Methods("POST")
	r.HandleFunc("/songs/{id}", controller.GetSongById).Methods("GET")
	r.HandleFunc("/songs/{id}", controller.UpdateSong).Methods("PUT")
	r.HandleFunc("/songs/{id}", controller.DeleteSong).Methods("DELETE")
}