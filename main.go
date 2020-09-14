package main

import (
	"log"
	"word-scrambler/pkg"
	"word-scrambler/server"

	_ "github.com/jinzhu/gorm/dialects/mysql"
)

func main() {
	pkg.ConnectMysql()
	pkg.AutoMigrate()
	if err := pkg.InitRedis(); err != nil {
		log.Fatal(err)
	}
	log.Println("Connection to redis success")
	server.Server()
}
