package main

import (
	"music-sharing/handlers"
	_ "music-sharing/models/track"
	"net/http"
)

func main() {
	//user := track.User{Username: "Jinzhu", FirstName: "Delilah", LastName: "Dessalegn"}
	//database.NewRecord(user) // => returns `true` as primary key is blank
	r := handlers.NewRouter()
	http.ListenAndServe(":8080", r)
	//database.Create(&user)

	//database.NewRecord(user)
}
