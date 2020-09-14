package pkg

import (
	"log"
	"word-scrambler/models"

	"github.com/jinzhu/gorm"
)

// DB -
var DB *gorm.DB

// ConnectMysql -
func ConnectMysql() {
	var err error
	DB, err = gorm.Open("mysql", "root:thomialghani@tcp(mysql_img:3306)/word_scrambler?charset=utf8mb4&parseTime=True&loc=Local")
	if err != nil {
		log.Println(err)
		log.Fatal("Failed to connect mysql")
	}
	log.Println("make connection success")
	DB.LogMode(true)
}

// AutoMigrate -
func AutoMigrate() {
	DB.AutoMigrate(
		&models.Question{},
		&models.UserAnswers{},
		&models.User{},
		&models.UserSummary{},
	)
	CreateIndex()
}

func CreateIndex() {
	DB.Model(&models.UserSummary{}).AddUniqueIndex("idx_user_summary_user_id", "user_id")
}
