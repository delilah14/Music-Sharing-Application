package handlers

import (
	"fmt"
	_ "fmt"
	"html/template"
	_ "html/template"
	"io/ioutil"
	"log"
	_ "log"

	_ "music-sharing/models/db"
	"music-sharing/models/track"
	_ "music-sharing/models/track"
	"music-sharing/utilities"
	_ "music-sharing/utilities"
	"net/http"
	_ "os"
	"time"
)

type PageVars struct {
	Username        string
	Password        string
	Tracks          []TrackList
	DisplayTracks   []track.Track
	DisplayPlaylist []track.Playlist
	SearchValue     string
	NotAvailable    bool
	PlaylistName    string
	NumberOfSongs   int
}

type TrackList struct {
	TrackName     string
	TrackDuration *time.Time
	ArtistName    string
	GenreName     string
	TrackPath     string
}

func SearchMusicHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	lists := []track.Track{}
	value := r.FormValue("searchvalue")
	track := track.Track{}
	track.Name = value
	println(value)
	database.Where("name = ?", value+".mp3").Find(&lists)
	println(len(lists))
	pageVars := PageVars{}
	pageVars.DisplayTracks = lists
	pageVars.SearchValue = value
	pageVars.NotAvailable = false
	if len(lists) == 0 {
		pageVars.NotAvailable = true
	}
	tmpl, err := template.ParseFiles(utilities.TemplatePath + "/assets/search.html")
	if err != nil { // if there is an error
		log.Print("template parsing error: ", err) // log it
	}
	err = tmpl.Execute(w, pageVars) //execute the template and pass it the HomePageVars struct to fill in the gaps
	if err != nil {                 // if there is an error
		log.Print("template executing error: ", err) //log it
	}

}

func GetAllMusicHandler(w http.ResponseWriter, r *http.Request) {

	lists := []track.Track{}
	plylists := []track.Playlist{}
	database.Find(&lists)
	database.Find(&plylists)
	p := PageVars{DisplayTracks: lists}
	p.DisplayPlaylist = plylists
	//println("Track path:", lists[0].TrackPath)
	tmpl, err := template.ParseFiles(utilities.TemplatePath + "/assets/index.html")

	if err != nil { // if there is an error
		log.Print("template parsing error: ", err) // log it
	}
	err = tmpl.Execute(w, p) //execute the template and pass it the HomePageVars struct to fill in the gaps
	if err != nil {          // if there is an error
		log.Print("template executing error: ", err) //log it
	}

}
func GetAddMusicPage(w http.ResponseWriter, r *http.Request) {

	lists := []track.Track{}
	database.Find(&lists)
	println("Tracks:", lists)

	p := PageVars{DisplayTracks: lists}
	tmpl, err := template.ParseFiles(utilities.TemplatePath + "/assets/index.html")

	if err != nil { // if there is an error
		log.Print("template parsing error: ", err) // log it
	}
	err = tmpl.Execute(w, p) //execute the template and pass it the HomePageVars struct to fill in the gaps
	if err != nil {          // if there is an error
		log.Print("template executing error: ", err) //log it
	}
}
func AddMusicHandler(w http.ResponseWriter, r *http.Request) {
	t := track.Track{}

	r.ParseMultipartForm(10 << 20)

	file, handler, err := r.FormFile("myFile")
	if err != nil {
		fmt.Println("Error Retrieving the File")
		fmt.Println(err)
		return
	}
	defer file.Close()
	fmt.Printf("Uploaded File: %+v\n", handler.Filename)
	fmt.Printf("File Size: %+v\n", handler.Size)
	fmt.Printf("MIME Header: %+v\n", handler.Header)

	fileBytes, err := ioutil.ReadAll(file)
	ioutil.WriteFile(utilities.TemplatePath+"/assets/img/"+handler.Filename, fileBytes, 777)
	if err != nil {
		fmt.Println(err)
	}

	t.Name = handler.Filename
	t.TrackPath = "/static/img/" + handler.Filename
	t.Size = byte(handler.Size)

	database.Create(&t)

	lists := []track.Track{}

	database.Find(&lists)

	p := PageVars{DisplayTracks: lists}

	tmpl, err := template.ParseFiles(utilities.TemplatePath + "/assets/index.html")

	if err != nil { // if there is an error
		log.Print("template parsing error: ", err) // log it
	}
	err = tmpl.Execute(w, p) //execute the template and pass it the HomePageVars struct to fill in the gaps
	if err != nil {          // if there is an error
		log.Print("template executing error: ", err) //log it
	}
	// return that we have successfully uploaded our file!

}
