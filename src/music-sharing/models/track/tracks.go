package track

import (
	"time"

	"github.com/jinzhu/gorm"
)

type Track struct {
	gorm.Model
	Name      string
	Duration  *time.Time
	ArtistId  int
	GenreId   int
	TrackPath string `sql:"type:text"`
	Size      byte
	Playlists []Playlist `gorm:"many2many:track_playlist;"`
}

type TrackUploaded struct {
	Id         int `gorm:"primary_key"`
	UploadedBy string
}

type Playlist struct {
	gorm.Model
	Name       string
	Created_by string
	Tracks     []Track `gorm:"many2many:track_playlist;"`
}

type Mood struct {
	gorm.Model
	Name string
}

type Genre struct {
	gorm.Model
	Name string
}

type User struct {
	gorm.Model
	Username  string
	FirstName string
	LastName  string
	Email     string
	Password  string `sql:"type:text"`
	ImageUrl  string `sql:"type:text"`
}
