package model

import "first/service/video/model/db"

// InitVideoDB init model
func InitVideoDB() {
	db.InitVideo("") // mysql init
}
