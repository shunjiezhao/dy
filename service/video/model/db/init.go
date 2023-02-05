package db

import (
	"first/pkg/constants"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"sync"
)

var VideoDb *gorm.DB
var once sync.Once

// InitVideo init VideoDb
func InitVideo() {
	once.Do(func() {
		var err error
		VideoDb, err = gorm.Open(mysql.Open(constants.MySQLDefaultDSN),
			&gorm.Config{
				PrepareStmt:            true,
				SkipDefaultTransaction: true,
			},
		)
		if err != nil {
			panic(err)
		}

		err = VideoDb.Set("gorm:table_options", "ENGINE=InnoDB DEFAULT CHARSET=utf8mb4").AutoMigrate(&FavouriteVideo{})

		//if err = DB.Use(gormopentracing.New()); err != nil {
		//	panic(err)
		//}
	})
}
