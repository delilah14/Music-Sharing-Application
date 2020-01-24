package handlers

import (
	_ "fmt"
	"html/template"
	_ "html/template"
	"log"
	_ "log"
	"net/http"

	_ "music-sharing/models/db"
	"music-sharing/models/track"
	_ "music-sharing/models/track"
	"music-sharing/utilities"
	_ "music-sharing/utilities"
	_ "os"
)

func GetAllPlaylistHandler(w http.ResponseWriter, r *http.Request) {
	/*play := track.Playlist{}
	play.Name = "New playlist"
	play.Created_by = "delilah14"
	database.Create(&play)
	*/

	lists := []track.Playlist{}
	database.Find(&lists)

	p := PageVars{DisplayPlaylist: lists}
	p.NotAvailable = false
	if len(lists) == 0 {
		p.NotAvailable = true
	}
	//println("Track path:", lists[0].Name)
	tmpl, err := template.ParseFiles(utilities.TemplatePath + "/assets/playlist.html")

	if err != nil { // if there is an error
		log.Print("template parsing error: ", err) // log it
	}
	err = tmpl.Execute(w, p) //execute the template and pass it the HomePageVars struct to fill in the gaps
	if err != nil {          // if there is an error
		log.Print("template executing error: ", err) //log it
	}

}

func AddPlaylistHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	play := track.Playlist{}
	play.Name = r.FormValue("addplaylistinput")
	play.Created_by = "delilah14"
	database.Create(&play)
	lists := []track.Playlist{}
	database.Find(&lists)
	p := PageVars{DisplayPlaylist: lists}
	tmpl, err := template.ParseFiles(utilities.TemplatePath + "/assets/playlist.html")

	if err != nil { // if there is an error
		log.Print("template parsing error: ", err) // log it
	}
	err = tmpl.Execute(w, p) //execute the template and pass it the HomePageVars struct to fill in the gaps
	if err != nil {          // if there is an error
		log.Print("template executing error: ", err) //log it
	}
}

func ViewSinglePlaylistHandler(w http.ResponseWriter, r *http.Request) {
	//playlist name and tracks with that playlist
	//get the playlist id with it's name
	//get list of tracks with the playlist id
	//set it on pagevars
	r.ParseForm()
	tracks := []track.Track{}
	playlist := track.Playlist{}
	database.Where("name = ?", r.FormValue("playlistname")).First(&playlist)
	database.Model(&playlist).Related(&tracks, "tracks")
	p := PageVars{DisplayTracks: tracks}
	p.PlaylistName = r.FormValue("playlistname")
	tmpl, err := template.ParseFiles(utilities.TemplatePath + "/assets/single_playlist.html")

	if err != nil { // if there is an error
		log.Print("template parsing error: ", err) // log it
	}
	err = tmpl.Execute(w, p) //execute the template and pass it the HomePageVars struct to fill in the gaps
	if err != nil {          // if there is an error
		log.Print("template executing error: ", err) //log it
	}

}

func AddMusicToPlaylistHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	trackInPlaylist := track.Track{}
	playlistOfTrack := track.Playlist{}
	list_playlist := []track.Playlist{}
	list_tracks := []track.Track{}
	database.Where("name = ?", r.FormValue("playid")).Find(&list_playlist)
	database.Where("name = ?", r.FormValue("trackid")).Find(&list_tracks)
	database.Where("name = ?", r.FormValue("trackid")).First(&trackInPlaylist)
	database.Where("name = ?", r.FormValue("playid")).First(&playlistOfTrack)
	trackInPlaylist.Playlists = list_playlist
	playlistOfTrack.Tracks = list_tracks
	database.Save(&trackInPlaylist)
	database.Save(&playlistOfTrack)
	GetAllMusicHandler(w, r)

}
