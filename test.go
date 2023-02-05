package main

import (
	"first/pkg/constants"
	db2 "first/service/user/model/db"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func main() {
	db, err := gorm.Open(mysql.Open(constants.MySQLDefaultDSN),
		&gorm.Config{
			PrepareStmt:            true,
			SkipDefaultTransaction: true,
		},
	)
	if err != nil {
		panic(err.Error())
	}

	db.Migrator().DropTable(&db2.Comment{})
	err = db.AutoMigrate(&db2.Comment{})
	db.Create(&db2.Comment{
		Id:      3,
		Uuid:    2,
		VideoId: 2,
		Content: "CS多少分",
	})
	if err != nil {
		panic(err)
	}
	com := &db2.Comment{}
	db.Debug().Preload("User").Find(&com)
	fmt.Println(com)

	//db.Create()

}

/*

1. 有用户id的话
select * from  video
(
	select video_id,author_uuid,title,play_url,cover_url,favourite_count,comment_count,
created_at from video_info order by limit 50
) t


*/
