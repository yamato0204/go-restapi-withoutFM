package main

import (
	"encoding/json"
	"fmt"

	"net/http"
)


type song struct{
	ID     string `json:"id"`
	Title  string `json :"title"`
	Singer string `json :"singer"`
}

type allSongs []song



var songs = allSongs{
	{ID: "1", Title: "Have a Nice Day", Singer: "Bon Jovi"},
	{ID: "2", Title: "The Nights", Singer: "Avicii"},
	{ID: "3", Title: "One Way Ticket", Singer: "ONE OK ROCK"},	
}



func getSong(w http.ResponseWriter, r *http.Request){
	urlID := r.URL.Path
	fmt.Println(urlID)
	for _, song := range songs{
		if song.ID == urlID {
			json.NewEncoder(w).Encode(song)
			break
		}
	}
}


func getSongs(w http.ResponseWriter, r *http.Request){
	json.NewEncoder(w).Encode(songs)
}


func createSong(w http.ResponseWriter, r *http.Request){

	var song song
	json.NewDecoder(r.Body).Decode(&song)
	songs = append(songs, song)

	json.NewEncoder(w).Encode(songs)

}

func updateSong(w http.ResponseWriter, r *http.Request){
	urlID := r.URL.Path

	var updateSong song

	decode := json.NewDecoder(r.Body) 
		
	if err := decode.Decode(&updateSong); err != nil{
		fmt.Println(err)
		return 
	} 

	for i, singleSong := range songs{
		if singleSong.ID == urlID {
			songs[i] = updateSong
			json.NewEncoder(w).Encode(updateSong)
			return

		}
	}	
}

func deleteSong(w http.ResponseWriter, r *http.Request){

	urlID := r.URL.Path
	for i, song := range songs{
		fmt.Println("songs: ", songs)
		if urlID == song.ID {
			songs = append(songs[:i], songs[i+1:]...)
			fmt.Println("deleted",songs)
		}
	}
}



func SingleSongHandler(w http.ResponseWriter, r *http.Request  ){
	switch r.Method {
	case "GET":
		getSong(w,r)
    case "DELETE":
		deleteSong(w,r)
	case "PUT":
		updateSong(w,r)
	default:
            w.WriteHeader(http.StatusMethodNotAllowed)
	}

}


func AllSongHandler(w http.ResponseWriter, r *http.Request){
	switch r.Method {
	case "GET":
		getSongs(w,r)
	case "POST":
		createSong(w,r)
	default:
            w.WriteHeader(http.StatusMethodNotAllowed)
	} 

}


func main(){
	http.HandleFunc("/songs", AllSongHandler)
	//http.HandleFunc("/songs/", SingleSongHandler) この場合　pathは　/songs/1
	http.Handle("/songs/", http.StripPrefix("/songs/", http.HandlerFunc(SingleSongHandler)))
	//http.HandlerFunc(SingleSongHandler)により、/songs/を取り除き、　pathは１
	http.ListenAndServe(":8080", nil)
}