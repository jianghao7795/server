package utils

import (
	global "server/model"

	"gorm.io/gorm"
)

func MergeQuery(db *gorm.DB, search []global.SearchValue) (err error) {
	for _, val := range search {
		switch val.SearchType {
		case "like":
			db = db.Where(val.Name+" like ?", val.Value)
		case "equal":
			db = db.Where(val.Name+" = ?", val.Value)
		case "increase":
			db = db.Order(val.Name)
		case "decrease":
			db = db.Order(val.Name + " desc") // desc
		default:
		}
	}
	return
}
