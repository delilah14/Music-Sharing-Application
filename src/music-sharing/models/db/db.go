package db

import (
	"music-sharing/models/track"
	_ "music-sharing/models/track"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

func InitDatabase() *gorm.DB {

	db, err := gorm.Open("mysql", "root:user@(localhost:3306)/music_streaming?charset=utf8&parseTime=True&loc=Local")
	if err != nil {
		panic("failed to connect database")
	}

	//db.AutoMigrate(&track.Track{})
	//db.Set("gorm:table_options", "ENGINE=InnoDB").AutoMigrate(&track.Track{})

	db.AutoMigrate(&track.Track{})
	db.AutoMigrate(&track.Playlist{})

	//db.AutoMigrate(&track.TrackUploaded{})

	//db.AutoMigrate(&track.User{})
	return db
}
